package jokes

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	v1 "github.com/dissoupov/chuck_jokes/api/v1"
	"github.com/dissoupov/chuck_jokes/internal/config"
	"github.com/go-phorce/dolly/metrics"
	"github.com/go-phorce/dolly/rest"
	"github.com/go-phorce/dolly/xhttp/header"
	"github.com/go-phorce/dolly/xhttp/httperror"
	"github.com/go-phorce/dolly/xhttp/identity"
	"github.com/go-phorce/dolly/xhttp/marshal"
	"github.com/go-phorce/dolly/xlog"
	"github.com/juju/errors"
)

// ServiceName provides the Service Name for this package
const ServiceName = "jokes"

var logger = xlog.NewPackageLogger("github.com/dissoupov/chuck_jokes/service", "jokes")

// Service defines the Status service
type Service struct {
	conf   *config.Configuration
	server rest.Server

	qurl string
}

// Factory returns a factory of the service
func Factory(server rest.Server) interface{} {
	if server == nil {
		logger.Panic("jokes.Factory: invalid parameter")
	}

	return func(conf *config.Configuration) {

		q := url.Values{}
		q.Add("firstName", "John")
		q.Add("lastName", "Doe")
		q.Add("limitTo", url.QueryEscape(strings.Join(conf.Jokes.Categories, ",")))

		svc := &Service{
			server: server,
			conf:   conf,
			qurl:   conf.Jokes.JokesSvcURL + "?" + q.Encode(),
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
	r.GET(v1.DefaultPathForJokes, s.jokesHandler())
	r.GET(v1.PathForJokes, s.jokesHandler())
}

func (s *Service) jokesHandler() rest.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ rest.Params) {
		identity.ForRequest(r)

		var uerr, jerr error
		var user *v1.UserName
		var joke *v1.JokeResponse

		var wg sync.WaitGroup
		wg.Add(2)

		// Fetch the name
		go func() {
			defer wg.Done()
			user, uerr = s.fetchName()
		}()

		go func() {
			defer wg.Done()
			joke, jerr = s.fetchJoke()
		}()

		wg.Wait()

		if uerr != nil || jerr != nil || joke.ResponseType != "success" {
			marshal.WriteJSON(w, r, httperror.WithUnexpected("failed to provide response"))
			return
		}

		jtext := strings.ReplaceAll(joke.Value.Joke, "John", user.First)
		jtext = strings.ReplaceAll(jtext, "Doe", user.Last)

		w.Header().Set(header.ContentType, header.TextPlain)
		w.Write([]byte(jtext))
	}
}

func (s *Service) fetchName() (name *v1.UserName, err error) {
	measure := createMetrics("name")
	defer measure(500, err)

	logger.Debugf("api=fetchName")

	var res *http.Response
	res, err = http.Get(s.conf.Jokes.NameSvcURL)
	if err != nil {
		logger.Errorf("failed to fetch name: %v", err)
		return nil, errors.Trace(err)
	}

	name = new(v1.UserName)
	if err = marshal.Decode(res.Body, name); err != nil {
		logger.Errorf("failed to parse name response: %v", err)
		return nil, errors.Trace(err)
	}
	return name, nil
}

func (s *Service) fetchJoke() (joke *v1.JokeResponse, err error) {
	measure := createMetrics("joke")
	defer measure(500, err)

	logger.Debugf("api=fetchJoke")

	var res *http.Response
	res, err = http.Get(s.qurl)
	if err != nil {
		logger.Errorf("failed to fetch joke: %v", err)
		return nil, errors.Trace(err)
	}

	joke = new(v1.JokeResponse)
	if err = marshal.Decode(res.Body, joke); err != nil {
		logger.Errorf("failed to parse name response: %v", err)
		return nil, errors.Trace(err)
	}
	return joke, nil
}

var (
	keyForOutboundPerf   = []string{"outbound", "perf"}
	keyForOutboundStatus = []string{"outbound", "status"}
)

func createMetrics(api string) func(rc int, err error) {
	started := time.Now().UTC()
	tagAPI := metrics.Tag{Name: "api", Value: api}
	return func(rc int, err error) {
		if err == nil {
			metrics.MeasureSince(keyForOutboundPerf, started, tagAPI)
			metrics.IncrCounter(keyForOutboundStatus, 1, tagAPI, metrics.Tag{Name: "status", Value: strconv.Itoa(rc)})
		} else {
			metrics.IncrCounter(keyForOutboundStatus, 1, tagAPI, metrics.Tag{Name: "status", Value: "error"})
		}
	}
}
