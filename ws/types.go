package ws

import (
	"github.com/madatsci/go-chat-server/internal/models"
	"time"
)

const eventMessage = "message"
const eventUserJoin = "join"

type (
	// MessageEvent represents incoming message event from websocket
	MessageEvent struct {
		From     string    `json:"from"`
		To       string    `json:"to"`
		Text     string    `json:"text"`
		DateTime time.Time `json:"date_time"`
	}

	// UserJoinEvent represents message that is used to notify users about new user entered
	UserJoinEvent struct {
		User string `json:"user"`
	}

	// Event represents websocket event for client
	Event struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}
)

// NewMessageEvent creates new Event instance with type "message"
func NewMessageEvent(message models.ChatMessage) Event {
	data := MessageEvent{
		From:     message.User.Email,
		Text:     message.Text,
		DateTime: message.CreatedAt,
	}

	if message.Receiver != nil {
		data.To = message.Receiver.Email
	}

	return Event{
		Type: eventMessage,
		Data: data,
	}
}

// NewUserJoinEvent creates new Event instance with type "join"
func NewUserJoinEvent(user models.User) Event {
	data := UserJoinEvent{
		User: user.Email,
	}

	return Event{
		Type: eventUserJoin,
		Data: data,
	}
}
