package http_server

import (
	"context"
	"net/http"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
	WithMux(mux http.Handler) Server
	GetAddr() string
}

type HTTPMux interface {
	HandleFunc(pattern string, fn http.HandlerFunc)
	Handle(pattern string, fn http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Group(prefix string, middlewares ...Middleware) MiddlewareChain
}

type MiddlewareChain interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handleFunc http.HandlerFunc)
}
