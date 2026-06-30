package http_server

import (
	"context"
	"errors"
	"net/http"
)

type httpServer struct {
	s *http.Server
}

func NewHTTPServer(cfg HTTPServerConfig) Server {
	s := &http.Server{
		Addr:                         cfg.Addr,
		DisableGeneralOptionsHandler: cfg.DisableGeneralOptionsHandler,
		ReadTimeout:                  cfg.ReadTimeout,
		WriteTimeout:                 cfg.WriteTimeout,
		ReadHeaderTimeout:            cfg.ReadHeaderTimeout,
		IdleTimeout:                  cfg.IdleTimeout,
	}

	return &httpServer{
		s: s,
	}
}

func (s *httpServer) WithMux(mux *http.ServeMux) Server {

	s.s.Handler = mux

	return s
}

func (s *httpServer) GetAddr() string {
	return s.s.Addr
}

func (s *httpServer) Start() error {
	if err := s.s.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	err := s.s.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
