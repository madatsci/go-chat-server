package internal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/madatsci/go-chat-server/api"
	"github.com/madatsci/go-chat-server/internal/providers"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"github.com/madatsci/go-chat-server/internal/services"
	"github.com/madatsci/go-chat-server/ws"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
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

// Migrate applies database migrations
func Migrate(dir string) {
	app := fx.New(
		fx.Provide(
			providers.NewConfig,
			providers.NewLogger,
		),

		fx.Invoke(func(v *viper.Viper, logger *zap.SugaredLogger) {
			if err := migrate(dir, v); err != nil {
				logger.Errorf("error applying migrations: %v", err)
			}

			logger.Infof("migrations successfully applied")
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func migrate(dir string, v *viper.Viper) error {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		v.GetString("db.user"),
		v.GetString("db.password"),
		v.GetString("db.addr"),
		v.GetString("db.db"),
	)

	db, err := sql.Open("postgres", conn)

	if err != nil {
		return err
	}

	if err := goose.Run("up", db, dir); err != nil {
		return err
	}

	return nil
}
