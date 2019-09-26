package internal

import (
	"github.com/madatsci/go-chat-server/internal/providers"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"github.com/madatsci/go-chat-server/internal/services"
	"github.com/madatsci/go-chat-server/api"
	"go.uber.org/fx"
)

func Run() {
	app := fx.New(
		fx.Provide(
			providers.NewConfig,
			providers.NewDB,
			providers.NewLogger,
			repositories.NewUserRepository,
			services.NewAccount,
		),

		fx.Invoke(
			api.New,
		),
	)

	app.Run()
}
