package http_server

import "net/http"

type Middleware func(next http.Handler) http.Handler

type mux struct {
	mux *http.ServeMux
}

func NewMux() HTTPMux {
	return &mux{
		mux: http.NewServeMux(),
	}
}

func (h *mux) HandleFunc(pattern string, fn http.HandlerFunc) {
	h.mux.HandleFunc(pattern, fn)
}

func (h *mux) Handle(pattern string, fn http.Handler) {
	h.mux.Handle(pattern, fn)
}

func (h *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *mux) Group(prefix string, middlewares ...Middleware) MiddlewareChain {
	return &group{
		mux:         h,
		prefix:      prefix,
		middlewares: middlewares,
	}
}
