package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/dissoupov/chuck_jokes/internal/config"
	"github.com/dissoupov/chuck_jokes/internal/version"
	"github.com/dissoupov/chuck_jokes/pkg/roles"
	"github.com/dissoupov/chuck_jokes/service/jokes"
	"github.com/dissoupov/chuck_jokes/service/status"
	"github.com/go-phorce/dolly/netutil"
	"github.com/go-phorce/dolly/rest"
	"github.com/go-phorce/dolly/rest/tlsconfig"
	"github.com/go-phorce/dolly/xhttp/authz"
	"github.com/go-phorce/dolly/xhttp/identity"
	"github.com/go-phorce/dolly/xlog"
	"github.com/go-phorce/dolly/xlog/logrotate"
	"github.com/juju/errors"
	"go.uber.org/dig"
	kp "gopkg.in/alecthomas/kingpin.v2"
)

var logger = xlog.NewPackageLogger("github.com/dissoupov/chuck_jokes/cmd/jokes", "main")

var serviceFactories = map[string]func(server rest.Server) interface{}{
	jokes.ServiceName:  jokes.Factory,
	status.ServiceName: status.Factory,
}

// return codes
const (
	rcError   = 1
	rcSuccess = 0
)

func main() {
	rc := rcSuccess

	app := newContainer(os.Args[1:])
	defer app.Close()

	err := app.start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		rc = rcError
	}
	os.Exit(rc)
}

// appFlags provides flags
type appFlags struct {
	cfgFile         *string
	cpu             *string
	isStderr        *bool
	bindHTTP        *string
	bindHTTPS       *string
	certRolesFile   *string
	apikeyRolesFile *string
	httpsCertFile   *string
	httpsKeyFile    *string
	httpsCAFile     *string
}

// app is the application container
type app struct {
	sigs      chan os.Signal
	container *dig.Container
	closers   []io.Closer
	closed    bool
	lock      sync.RWMutex

	args  []string
	flags *appFlags
	cfg   *config.Configuration
}

func newContainer(args []string) *app {
	return &app{
		sigs:      make(chan os.Signal, 2),
		args:      args,
		container: dig.New(),
		closers:   make([]io.Closer, 0, 8),
		flags:     new(appFlags),
	}
}

// OnClose adds a closer interface to the list to be called when application exits
func (a *app) OnClose(closer io.Closer) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if closer != nil {
		a.closers = append(a.closers, closer)
	}
}

// Close is called on application exit to free resources
func (a *app) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.closed {
		return errors.New("already closed")
	}

	a.closed = true
	// close in reverse order
	for i := len(a.closers) - 1; i >= 0; i-- {
		closer := a.closers[i]
		if closer != nil {
			closer.Close()
		}
	}
	logger.Warning("api=app.Close, status=service_stopped")

	return nil
}

func (a *app) loadConfig() error {
	app := kp.New("jokes", "Demo Web Server")
	app.HelpFlag.Short('h')
	verInfo := fmt.Sprintf("jokes %v", version.Current())
	app.Version(verInfo)

	flags := a.flags
	flags.cfgFile = app.Flag("cfg", "Location of the configuration file").Default(config.ConfigFileName).Short('c').String()
	flags.cpu = app.Flag("cpu", "Enable CPU profiling, specify a file to store CPU profiling info").String()
	flags.isStderr = app.Flag("std", "Output logs to stderr").Bool()
	flags.bindHTTP = app.Flag("bind-http", "Bind address for Public HTTP end-point").String()
	flags.bindHTTPS = app.Flag("bind-https", "Bind address for Public WebAPI end-point").String()
	flags.apikeyRolesFile = app.Flag("roles-apikey-file", "Location of the config file for API-Key role mapper").String()
	flags.certRolesFile = app.Flag("roles-cert-file", "Location of the config file for certificate role mapper").String()
	flags.httpsCertFile = app.Flag("https-cert-file", "Path to the server TLS cert file.").String()
	flags.httpsKeyFile = app.Flag("https-key-file", "Path to the server TLS key file.").String()
	flags.httpsCAFile = app.Flag("https-trusted-ca-file", "Path to the server TLS trusted CA file.").String()

	// Parse arguments
	kp.MustParse(app.Parse(a.args))

	cfg, _, err := config.LoadConfig(*flags.cfgFile)
	if err != nil {
		return errors.Annotatef(err, "failed to load configuration %q", *flags.cfgFile)
	}
	//logger.Infof("api=loadConfig, status=loaded, folder=%q", absCfgFile)
	a.cfg = cfg

	if *flags.httpsCertFile != "" {
		cfg.HTTPS.ServerTLS.CertFile = *flags.httpsCertFile
	}
	if *flags.httpsKeyFile != "" {
		cfg.HTTPS.ServerTLS.KeyFile = *flags.httpsKeyFile
	}
	if *flags.httpsCAFile != "" {
		cfg.HTTPS.ServerTLS.TrustedCAFile = *flags.httpsCAFile
	}
	if *flags.bindHTTP != "" {
		cfg.HTTP.BindAddr = *flags.bindHTTP
	}
	if *flags.bindHTTPS != "" {
		cfg.HTTPS.BindAddr = *flags.bindHTTPS
	}
	if *flags.certRolesFile != "" {
		cfg.Authz.CertMapper = *flags.certRolesFile
	}
	if *flags.apikeyRolesFile != "" {
		cfg.Authz.APIKeyMapper = *flags.apikeyRolesFile
	}
	a.container.Provide(func() (*config.Configuration, *appFlags) {
		return a.cfg, a.flags
	})

	return nil
}

func (a *app) initLogs() error {
	cfg := a.cfg
	if cfg.Logger.Directory != "" {
		var sink io.Writer
		if *a.flags.isStderr {
			sink = os.Stderr
			xlog.SetFormatter(xlog.NewColorFormatter(sink, true))
		}
		logRotate, err := logrotate.Initialize(cfg.Logger.Directory, cfg.ServiceName, cfg.Logger.MaxAgeDays, cfg.Logger.MaxSizeMb, true, sink)
		if err != nil {
			return errors.Annotate(err, "failed to initialize log rotate")
		}

		a.OnClose(logRotate)
	} else {
		formatter := xlog.NewColorFormatter(os.Stderr, true)
		xlog.SetFormatter(formatter)
	}

	// Set log levels for each repo
	if cfg.LogLevels != nil {
		for _, ll := range cfg.LogLevels {
			l, _ := xlog.ParseLevel(ll.Level)
			if ll.Repo == "*" {
				xlog.SetGlobalLogLevel(l)
			} else {
				xlog.SetPackageLogLevel(ll.Repo, ll.Package, l)
			}
			logger.Debugf("api=start, logger=%q, level=%v", ll.Repo, l)
		}
	}

	logger.Infof("api=initLogs, status=service_starting, version='%v', runtime='%v', args=%v, config=%q",
		version.Current(), runtime.Version(), os.Args, *a.flags.cfgFile)

	return nil
}

func (a *app) start() error {
	err := a.loadConfig()
	if err != nil {
		return errors.Trace(err)
	}

	err = a.initLogs()
	if err != nil {
		return errors.Trace(err)
	}

	// if this service is started on boot, ensure that network is available
	ipaddr, err := netutil.WaitForNetwork(30 * time.Second)
	if err != nil {
		return errors.Annotate(err, "unable to resolve local IP")
	}
	logger.Infof("api=start, ipaddr=%q", ipaddr)

	err = a.container.Provide(func(cfg *config.Configuration) (rest.Authz, error) {
		var azp rest.Authz
		if len(cfg.Authz.Allow) > 0 ||
			len(cfg.Authz.AllowAny) > 0 ||
			len(cfg.Authz.AllowAnyRole) > 0 {
			azp, err = authz.New(&authz.Config{
				Allow:        cfg.Authz.Allow,
				AllowAny:     cfg.Authz.AllowAny,
				AllowAnyRole: cfg.Authz.AllowAnyRole,
				LogAllowed:   cfg.Authz.GetLogAllowed(),
				LogDenied:    cfg.Authz.GetLogDenied(),
			})
			if err != nil {
				return nil, errors.Trace(err)
			}

			p, err := roles.New(
				cfg.Authz.APIKeyMapper,
				cfg.Authz.CertMapper,
			)
			if err != nil {
				return nil, errors.Trace(err)
			}
			identity.SetGlobalIdentityMapper(p.IdentityMapper)
		}
		return azp, nil
	})
	if err != nil {
		return errors.Trace(err)
	}

	stopServers := func(servers []rest.Server) {
		for _, running := range servers {
			running.StopHTTP()
		}
	}
	servers := []rest.Server{}

	for _, svcCfg := range []*config.HTTPServer{&a.cfg.HTTP, &a.cfg.HTTPS} {
		httpServer, err := createHTTPServer(ipaddr, svcCfg, a.container)
		if err != nil {
			logger.Errorf("api=start, reason=createHTTPServer, service=%s, err=[%v]", svcCfg.ServiceName, errors.ErrorStack(err))
			stopServers(servers)

			return errors.Trace(err)
		}
		servers = append(servers, httpServer)
	}

	// register for signals, and wait to be shutdown
	signal.Notify(a.sigs, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGABRT)
	// Block until a signal is received.
	sig := <-a.sigs
	logger.Warningf("api=start, status=shuting_signal, sig=%v", sig)

	stopServers(servers)

	// let to stop
	time.Sleep(time.Second * 3)

	return nil
}

var tlsStrToClientAuthMap = map[string]tls.ClientAuthType{
	"NoClientCert":               tls.NoClientCert,
	"RequestClientCert":          tls.RequestClientCert,
	"RequireAnyClientCert":       tls.RequireAnyClientCert,
	"VerifyClientCertIfGiven":    tls.VerifyClientCertIfGiven,
	"RequireAndVerifyClientCert": tls.RequireAndVerifyClientCert,
}

func createHTTPServer(
	ipaddr string,
	cfgHTTPServer *config.HTTPServer,
	container *dig.Container,
) (rest.Server, error) {
	var err error
	var server rest.Server
	var tlsCfg *tls.Config
	var tlsloader *tlsconfig.KeypairReloader

	err = container.Invoke(func(
		cfg *config.Configuration,
		azp rest.Authz,
	) error {
		if cfgHTTPServer.ServerTLS.KeyFile != "" && cfgHTTPServer.ServerTLS.CertFile != "" {
			clientauthType := tls.VerifyClientCertIfGiven
			if ct, ok := tlsStrToClientAuthMap[cfgHTTPServer.ServerTLS.GetClientCertAuth()]; ok {
				clientauthType = ct
			}
			tlsCfg, err = tlsconfig.NewServerTLSFromFiles(
				cfgHTTPServer.ServerTLS.CertFile,
				cfgHTTPServer.ServerTLS.KeyFile,
				cfgHTTPServer.ServerTLS.TrustedCAFile,
				clientauthType)
			if err != nil {
				return errors.Trace(err)
			}

			tlsloader, err = tlsconfig.NewKeypairReloader(
				cfgHTTPServer.ServerTLS.CertFile,
				cfgHTTPServer.ServerTLS.KeyFile,
				5*time.Second)
			if err != nil {
				return errors.Annotatef(err, "api=createHTTPServer, reason=NewKeypairReloader, cert=%q, key=%q",
					cfgHTTPServer.ServerTLS.CertFile, cfgHTTPServer.ServerTLS.KeyFile)
			}
			tlsCfg.GetCertificate = tlsloader.GetKeypairFunc()
		}

		server, err = rest.New(version.Current().String(), ipaddr, cfgHTTPServer, tlsCfg, nil, azp, nil, nil)
		if err != nil {
			return errors.Annotatef(err, "api=createHTTPServer, reason=unable_initialize_service, name=%q", cfgHTTPServer.ServiceName)
		}
		return nil
	})
	if err != nil {
		if tlsloader != nil {
			tlsloader.Close()
		}
		return nil, errors.Trace(err)
	}

	server.OnEvent(rest.ServerStoppedEvent, func(evt rest.ServerEvent) {
		if tlsloader != nil {
			tlsloader.Close()
		}
	})

	for _, name := range cfgHTTPServer.Services {
		sf := serviceFactories[name]
		if sf == nil {
			return nil, errors.Errorf("service factory is not registered: %q", name)
		}
		err = container.Invoke(sf(server))
		if err != nil {
			return nil, errors.Annotatef(err, "api=createHTTPServer, reason='unable to start HTTP service', service=%q, factory=%s",
				cfgHTTPServer.ServiceName, name)
		}
	}

	if err := server.StartHTTP(); err != nil {
		return nil, errors.Annotatef(err, "api=createHTTPServer, reason='unable to start HTTP service', service=%q",
			cfgHTTPServer.ServiceName)
	}

	return server, nil
}
