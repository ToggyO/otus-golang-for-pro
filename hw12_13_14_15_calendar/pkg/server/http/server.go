package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, cfg configuration.Configuration) Server {
	return Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: handler,
			// TODO: вынести в конфиг
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start(_ context.Context) error {
	s.check()
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.check()
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) check() {
	if s.httpServer == nil {
		panic("There no http server instance is provided!")
	}
}
