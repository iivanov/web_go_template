package webserver

import (
	"net/http"

	"go.uber.org/fx"
)

type Router struct {
	mux *http.ServeMux
}

type RouterParams struct {
	fx.In
	Routes []Route `group:"routes"`
}

func NewRouter(params RouterParams) *Router {
	mux := http.NewServeMux()
	for _, r := range params.Routes {
		mux.Handle(r.Pattern(), r)
	}
	return &Router{mux: mux}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
