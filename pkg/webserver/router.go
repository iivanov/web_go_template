package webserver

import (
	"log/slog"
	"net/http"

	"go.uber.org/fx"
)

type Router struct {
	mux     *http.ServeMux
	handler http.Handler
}

type RouterParams struct {
	fx.In
	Logger      *slog.Logger
	Routes      []Route      `group:"routes"`
	AppRoutes   []AppRoute   `group:"approutes"`
	Middlewares []Middleware `group:"middlewares"`
}

func NewRouter(params RouterParams) *Router {
	mux := http.NewServeMux()
	for _, r := range params.Routes {
		mux.Handle(r.Pattern(), r)
	}
	for _, r := range params.AppRoutes {
		adapter := &appRouteAdapter{route: r, logger: params.Logger}
		mux.Handle(adapter.Pattern(), adapter)
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
