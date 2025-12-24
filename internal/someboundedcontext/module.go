package someboundedcontext

import (
	controller "project_template/internal/someboundedcontext/controllers"
	"project_template/internal/someboundedcontext/services"
	"project_template/pkg/webserver"

	"go.uber.org/fx"
)

var Module = fx.Module("someboundedcontext",
	fx.Provide(
		services.NewUserService,
		webserver.AsRoute(controller.NewUserHandler),
		webserver.AsRoute(controller.NewUsersHandler),
		webserver.AsRoute(controller.NewCreateUserHandler),
	),
)
