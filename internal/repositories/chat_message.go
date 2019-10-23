package repositories

import (
	"github.com/go-pg/pg"
	"github.com/madatsci/go-chat-server/internal/models"
)

type (
	ChatMessage interface {
		Create(from, to, text string) (*models.ChatMessage, error)
		GetAll() ([]models.ChatMessage, error)
	}

	chatMessageRepository struct {
		db *pg.DB
	}
)

func NewChatMessageRepository(db *pg.DB) ChatMessage {
	return &chatMessageRepository{
		db: db,
	}
}

// Create saves new chat message in database
func (c *chatMessageRepository) Create(from, to, text string) (*models.ChatMessage, error) {
	chatMessage := models.ChatMessage{
		From: from,
		To:   to,
		Text: text,
	}

	if _, err := c.db.Model(&chatMessage).Insert(); err != nil {
		return nil, err
	}

	return &chatMessage, nil
}

// GetAll retrieves all chat messages from database
func (c *chatMessageRepository) GetAll() ([]models.ChatMessage, error) {
	var chatMessages []models.ChatMessage

	if err := c.db.Model(&chatMessages).Select(); err != nil {
		return nil, err
	}

	return chatMessages, nil
}
