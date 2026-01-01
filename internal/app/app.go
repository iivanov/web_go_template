package app

import (
	"project_template/internal/someboundedcontext"
	"project_template/pkg/database"
	"project_template/pkg/logger"
	"project_template/pkg/telemetry"
	"project_template/pkg/webserver"

	"go.uber.org/fx"
)

func generalModules() []fx.Option {
	return []fx.Option{
		logger.Module,
		telemetry.Module,
		someboundedcontext.Module,
		database.Module,
	}
}

func middlewares() []fx.Option {
	return []fx.Option{
		fx.Provide(webserver.AsMiddleware(telemetry.NewHTTPMiddleware)),
		fx.Provide(webserver.AsMiddleware(telemetry.NewHTTPMetricsMiddleware)),
	}
}

func Serve(configFile string) {
	fx.New(
		fx.Options(generalModules()...),
		fx.Options(middlewares()...),
		fx.Provide(NewServeConfig(configFile)),
		webserver.Module,
	).Run()
}
