package status

import (
	"net/http"
	"strings"

	v1 "github.com/dissoupov/chuck_jokes/api/v1"
	"github.com/dissoupov/chuck_jokes/internal/config"
	"github.com/dissoupov/chuck_jokes/internal/version"
	"github.com/dissoupov/chuck_jokes/pkg/printer"
	"github.com/go-phorce/dolly/rest"
	"github.com/go-phorce/dolly/xhttp/header"
	"github.com/go-phorce/dolly/xhttp/identity"
	"github.com/go-phorce/dolly/xhttp/marshal"
	"github.com/go-phorce/dolly/xlog"
)

// ServiceName provides the Service Name for this package
const ServiceName = "status"

var logger = xlog.NewPackageLogger("github.com/dissoupov/chuck_jokes/service", "status")

// Service defines the Status service
type Service struct {
	conf   *config.Configuration
	server rest.Server
}

// Factory returns a factory of the service
func Factory(server rest.Server) interface{} {
	if server == nil {
		logger.Panic("status.Factory: invalid parameter")
	}

	return func(conf *config.Configuration) {
		svc := &Service{
			server: server,
			conf:   conf,
		}

		server.AddService(svc)
	}
}

// Name returns the service name
func (s *Service) Name() string {
	return ServiceName
}

// IsReady indicates that the service is ready to serve its end-points
func (s *Service) IsReady() bool {
	return true
}

// Close the subservices and it's resources
func (s *Service) Close() {
}

// Register adds the Status API endpoints to the overall URL router
func (s *Service) Register(r rest.Router) {
	r.GET(v1.PathForStatusVersion, versionHandler(s))
	r.GET(v1.PathForStatus, statusHandler(s))
	r.GET(v1.PathForStatusServer, statusHandler(s))
}

func versionHandler(s *Service) rest.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ rest.Params) {
		identity.ForRequest(r)

		w.Header().Set(header.ContentType, header.TextPlain)
		w.Write([]byte(s.server.Version()))
	}
}

func statusHandler(s *Service) rest.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ rest.Params) {
		identity.ForRequest(r)

		ver := version.Current()
		statusResponse := v1.ServerStatusResponse{
			Version: &v1.ServerVersion{
				Build:   ver.Build,
				Runtime: ver.Runtime,
			},
		}

		statusResponse.Status = &v1.ServerStatus{
			HostName:  s.server.HostName(),
			Port:      s.server.Port(),
			StartedAt: s.server.StartedAt(),
			Uptime:    s.server.Uptime(),
			Version:   s.server.Version(),
		}

		accept := r.Header.Get(header.Accept)
		if accept == "" || strings.EqualFold(accept, header.ApplicationJSON) {
			marshal.WriteJSON(w, r, statusResponse)
		} else {
			w.Header().Set(header.ContentType, header.TextPlain)
			printer.PrintServerStatus(w, statusResponse.Status)
		}
	}
}
