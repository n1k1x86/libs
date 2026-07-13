package http_server

import (
	"net/http"
	"strings"
)

type group struct {
	mux         HTTPMux
	prefix      string
	middlewares []Middleware
}

func (g *group) Handle(pattern string, handler http.Handler) {
	prefix := g.addPrefix(pattern)

	g.mux.Handle(prefix, Chain(handler, g.middlewares...))
}

func (g *group) HandleFunc(pattern string, handleFunc http.HandlerFunc) {
	g.Handle(pattern, handleFunc)
}

func (g *group) addPrefix(pattern string) string {
	method, path, found := strings.Cut(pattern, " ")
	if found {
		return method + " " + g.prefix + path
	}

	return g.prefix + pattern
}

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
