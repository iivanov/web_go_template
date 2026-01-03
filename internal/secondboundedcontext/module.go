package secondboundedcontext

import (
	"project_template/internal/secondboundedcontext/handlers"
	"project_template/pkg/messagebus"

	"go.uber.org/fx"
)

var Module = fx.Module("secondboundedcontext",
	fx.Provide(
		messagebus.AsHandler(handlers.NewUserCreatedHandler),
	),
)
