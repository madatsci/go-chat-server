package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/madatsci/go-chat-server/internal/services"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type (
	WebSocket struct {
		logger         *zap.SugaredLogger
		config         *viper.Viper
		accountService services.Account
		history        []Message
		hub            map[string]User
		hmu            sync.RWMutex
	}

	User struct {
		Model *models.User
		Conn  *websocket.Conn
	}

	Options struct {
		fx.In

		Logger         *zap.SugaredLogger
		Config         *viper.Viper
		Lc             fx.Lifecycle
		AccountService services.Account
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
		logger:         opts.Logger,
		config:         opts.Config,
		accountService: opts.AccountService,
		hub:            make(map[string]User),
	}

	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := opts.Config.GetString("ws.addr")
			opts.Logger.Infof("starting websocket server at: %s", addr)

			go func() {
				if err := http.ListenAndServe(addr, socket); err != nil {
					opts.Logger.Errorf("error starting websocket server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func (s *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Connecting
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("error creating upgrader: %v", err)
		return
	}
	defer c.Close()

	// Authentication
	token := r.URL.Query().Get("token")
	user, err := s.accountService.ValidateToken(token)
	if err != nil || user == nil {
		s.logger.Errorf("unable to find user by token '%s': %v", token, err)
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	s.hmu.Lock()
	s.hub[user.Email] = User{
		Model: user,
		Conn:  c,
	}
	s.hmu.Unlock()

	// Sending history to user
	s.logger.Info("sending history to user")
	for _, message := range s.history {
		if err := c.WriteJSON(message); err != nil {
			s.logger.Errorf("error sending history: %v", err)
		}
	}

	// Message handling
	s.logger.Info("start listening for incoming messages")
	for {
		var msg Message
		if err := c.ReadJSON(&msg); err != nil {
			s.logger.Errorf("error reading message: %v", err)
			return
		}

		s.logger.Infof("got incoming message: %v", msg)

		msg.From = user.Email
		msg.DateTime = time.Now()
		if msg.To != "" {
			s.hmu.RLock()
			userTo, ok := s.hub[msg.To]
			s.hmu.RUnlock()

			if !ok {
				s.logger.Errorf("unknown user to send message: %s", msg.To)
				continue
			}

			msg.To = userTo.Model.Email

			if err := c.WriteJSON(msg); err != nil {
				s.logger.Errorf("error sending message: %v", err)
			}

			if err := userTo.Conn.WriteJSON(msg); err != nil {
				s.logger.Errorf("error sending message: %v", err)
			}
		} else {
			s.hmu.RLock()
			for _, user := range s.hub {
				if err := user.Conn.WriteJSON(msg); err != nil {
					s.logger.Error("error sending message: %v", err)
					continue
				}
			}
			s.history = append(s.history, msg)
			s.hmu.RUnlock()
		}
	}
}
