package webserver

import (
	"net/http"

	"go.uber.org/fx"
)

type Router struct {
	mux     *http.ServeMux
	handler http.Handler
}

type RouterParams struct {
	fx.In
	Routes      []Route      `group:"routes"`
	Middlewares []Middleware `group:"middlewares"`
}

func NewRouter(params RouterParams) *Router {
	mux := http.NewServeMux()
	for _, r := range params.Routes {
		mux.Handle(r.Pattern(), r)
	}

	// Apply middlewares in reverse order so the first middleware in the slice
	// is the outermost wrapper (executed first on request, last on response)
	var handler http.Handler = mux
	for i := len(params.Middlewares) - 1; i >= 0; i-- {
		handler = params.Middlewares[i](handler)
	}

	return &Router{mux: mux, handler: handler}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}
