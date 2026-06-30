package http_server

import "net/http"

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

func (h *mux) GetMux() *http.ServeMux {
	return h.mux
}
