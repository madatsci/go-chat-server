package api

import (
	"context"
	"github.com/spf13/viper"

	"github.com/labstack/echo"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"github.com/madatsci/go-chat-server/internal/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	// Api structure represents API interaction
	Api struct {
		logger         *zap.SugaredLogger
		config         *viper.Viper
		userRepo       repositories.User
		accountService services.Account
		echo           *echo.Echo
	}

	// ApiOptions represents options needed to instantiate Api
	ApiOptions struct {
		fx.In

		Logger         *zap.SugaredLogger
		Config         *viper.Viper
		UserRepo       repositories.User
		AccountService services.Account
		Lc             fx.Lifecycle
	}
)

// New creates a new instance of Api and injects it into fx.Lifecycle
func New(opts ApiOptions) {
	a := &Api{
		logger:         opts.Logger,
		config:         opts.Config,
		userRepo:       opts.UserRepo,
		accountService: opts.AccountService,
		echo:           echo.New(),
	}

	a.echo.HidePort = true
	a.echo.HideBanner = true

	// Endpoint
	a.echo.POST("/register", a.Register)
	a.echo.POST("/auth", a.Auth)

	a.echo.GET("/profile", a.Profile, a.AuthMiddleware)

	// Start & Stop server
	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := opts.Config.GetString("api.addr")
			opts.Logger.Infof("Starting server at: %s", addr)
			go func() {
				if err := a.echo.Start(addr); err != nil {
					opts.Logger.Errorf("error starting server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return a.echo.Shutdown(ctx)
		},
	})
}
