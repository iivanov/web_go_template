package someboundedcontext

import (
	controller "project_template/internal/someboundedcontext/controllers"
	"project_template/internal/someboundedcontext/entities"
	"project_template/internal/someboundedcontext/repositories"
	"project_template/internal/someboundedcontext/services"
	"project_template/pkg/webserver"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Module("someboundedcontext",
	fx.Provide(
		services.NewUserService,
		repositories.NewUserRepository,
		webserver.AsRoute(controller.NewUserHandler),
		webserver.AsRoute(controller.NewUsersHandler),
		webserver.AsRoute(controller.NewCreateUserHandler),
	),
	fx.Invoke(func(db *gorm.DB) error {
		return db.AutoMigrate(&entities.User{})
	}),
)
