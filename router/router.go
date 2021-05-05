// Package router lets you register handlers for path regexen.
package router

import (
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// RegexpRouter maps path regexen to HTTP handlers
type RegexpRouter struct {
	routes []*route
}

// Handler registers a handler for the given path regexp.
func (h *RegexpRouter) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

// HandleFunc registers a handler func for the given path regexp.
func (h *RegexpRouter) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

// ServeHTTP delegates HTTP requests to the handler that matches their path,
// serves 404 if no match is found.
func (h *RegexpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
