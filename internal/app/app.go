package app

import (
	"project_template/internal/someboundedcontext"
	"project_template/pkg/logger"
	"project_template/pkg/webserver"

	"go.uber.org/fx"
)

func generalModules() []fx.Option {
	return []fx.Option{
		logger.Module,
		someboundedcontext.Module,
	}
}

func Serve(configFile string) {
	fx.New(
		fx.Options(generalModules()...),
		fx.Provide(NewServeConfig(configFile)),
		webserver.Module,
	).Run()
}
