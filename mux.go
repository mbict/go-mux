package mux

import (
	"net/http"
	"sort"
)

type route struct {
	pattern string
	handler http.Handler
}

type Mux struct {
	routes []*route
}

func New() *Mux {
	return &Mux{}
}

func (mux *Mux) Handler(req *http.Request) (h http.Handler, pattern string) {
	h, pattern = mux.matchRoute(req.URL.Path)
	if h == nil {
		h = http.NotFoundHandler()
	}
	return
}

func (mux *Mux) Handle(pattern string, h http.Handler) {
	if h == nil {
		panic("mux: nil handler")
	}
	if h, p := mux.matchRoute(pattern); h != nil && p == pattern {
		panic("mux: duplicate path for `" + pattern + "`")
	}

	mux.routes = append(mux.routes, &route{
		pattern: pattern,
		handler: h,
	})

	// we always sort the routes from longest to shortest pattern
	sort.Slice(mux.routes, func(i, j int) bool {
		return len(mux.routes[i].pattern) >= len(mux.routes[j].pattern)
	})
}

func (mux *Mux) HandleFunc(pattern string, hf http.HandlerFunc) {
	var h http.Handler
	if hf != nil {
		h = http.HandlerFunc(hf)
	}

	mux.Handle(pattern, h)
}

func (mux *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h, _ := mux.Handler(req)
	h.ServeHTTP(rw, req)
}

func (mux *Mux) matchRoute(path string) (h http.Handler, pattern string) {
	for _, r := range mux.routes {
		if matchPath(r.pattern, path) {
			return r.handler, r.pattern
		}
	}
	return nil, ""
}

func matchPath(pattern, path string) bool {
	n := len(pattern)
	return len(path) >= n && path[0:n] == pattern
}
