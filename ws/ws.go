package ws

import (
	"github.com/gorilla/websocket"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/internal/repositories"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type (
	WebSocket struct {
		logger *zap.SugaredLogger
		config *viper.Viper
		userRepo repositories.User
		history []Message
		hub map[string]User
		hmu sync.RWMutex
	}

	User struct {
		Model *models.User
		Conn *websocket.Conn
	}

	Options struct {
		fx.In

		Logger *zap.SugaredLogger
		Config *viper.Viper
		Lc fx.Lifecycle
		UserRepo repositories.User
	}
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func New(opts Options) {
	socket := &WebSocket{
		logger:   opts.Logger,
		config:   opts.Config,
		userRepo: opts.UserRepo,
		hub:      make(map[string]User),
	}
}
