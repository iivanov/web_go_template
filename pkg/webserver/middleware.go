package webserver

import (
	"net/http"

	"go.uber.org/fx"
)

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// AsMiddleware annotates the given constructor to state that
// it provides a middleware to the "middlewares" group.
func AsMiddleware(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"middlewares"`),
	)
}
