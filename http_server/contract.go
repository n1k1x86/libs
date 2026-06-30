package http_server

import (
	"context"
	"net/http"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
	WithMux(mux *http.ServeMux) Server
	GetAddr() string
}
