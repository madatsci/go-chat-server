package internal

import (
	"github.com/madatsci/go-chat-server/api"
	"github.com/madatsci/go-chat-server/internal/providers"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"github.com/madatsci/go-chat-server/internal/services"
	"github.com/madatsci/go-chat-server/ws"
	"go.uber.org/fx"
)

// Run starts main application
func Run() {
	app := fx.New(
		fx.Provide(
			providers.NewConfig,
			providers.NewDB,
			providers.NewLogger,
			providers.NewBcryptHasher,
			repositories.NewUserRepository,
			repositories.NewChatMessageRepository,
			services.NewAccount,
		),

		fx.Invoke(
			api.New,
			ws.New,
		),
	)

	app.Run()
}
