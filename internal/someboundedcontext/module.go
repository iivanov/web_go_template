package someboundedcontext

import (
	controller "project_template/internal/someboundedcontext/controllers"
	"project_template/internal/someboundedcontext/repositories"
	"project_template/internal/someboundedcontext/services"
	"project_template/pkg/webserver"

	"go.uber.org/fx"
)

var Module = fx.Module("someboundedcontext",
	fx.Provide(
		services.NewUserService,
		repositories.NewUserRepository,
		webserver.AsAppRoute(controller.NewUserHandler),
		webserver.AsAppRoute(controller.NewUsersHandler),
		webserver.AsAppRoute(controller.NewCreateUserHandler),
	),
)
