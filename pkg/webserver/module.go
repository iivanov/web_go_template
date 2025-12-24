package webserver

import (
	"net/http"

	"go.uber.org/fx"
)

var Module = fx.Module("webserver",
	fx.Provide(
		NewRouter,
		NewHTTPServer,
	),
	fx.Invoke(func(*http.Server) {}),
)
