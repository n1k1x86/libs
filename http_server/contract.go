package http_server

import (
	"context"
	"net/http"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
	WithMux(mux HTTPMux) Server
	GetAddr() string
}

type HTTPMux interface {
	HandleFunc(pattern string, fn http.HandlerFunc)
	Handle(patter string, fn http.Handler)
	GetMux() *http.ServeMux
}
