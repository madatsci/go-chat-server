package api

import (
	"context"

	"github.com/labstack/echo"
	"github.com/madatsci/go-chat-server/internal/providers"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"github.com/madatsci/go-chat-server/internal/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	Api struct {
		logger         *zap.SugaredLogger
		config         *providers.Config
		userRepo       repositories.User
		accountService services.Account
	}

	ApiOptions struct {
		fx.In

		Logger         *zap.SugaredLogger
		Config         *providers.Config
		UserRepo       repositories.User
		AccountService services.Account
		Lc             fx.Lifecycle
	}
)

func New(opts ApiOptions) {
	a := &Api{
		logger:         opts.Logger,
		config:         opts.Config,
		userRepo:       opts.UserRepo,
		accountService: opts.AccountService,
	}

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	// Endpoint
	e.POST("/register", a.Register)
	e.POST("/auth", a.Auth)

	e.GET("/profile", a.Profile, a.AuthMiddleware)

	// Start & Stop server
	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			opts.Logger.Infof("Starting server at: %s", opts.Config.ListenAddr)
			go e.Start(opts.Config.ListenAddr)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}
