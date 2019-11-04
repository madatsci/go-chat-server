package repositories

import (
	"github.com/go-pg/pg"
	"github.com/madatsci/go-chat-server/internal/models"
)

type (
	ChatMessage interface {
		Create(message *models.ChatMessage) error
		GetLastPublicMessages(limit int) ([]models.ChatMessage, error)
		GetLastPrivateMessages(user models.User, limit int) ([]models.ChatMessage, error)
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
func (c *chatMessageRepository) Create(message *models.ChatMessage) error {
	if _, err := c.db.Model(message).Relation("User").Insert(); err != nil {
		return err
	}

	if err := c.db.Model(message).
		Column("chat_message.*").
		Relation("User").
		Relation("Receiver").
		Where("message.id=?", message.ID).
		First(); err != nil {
		return err
	}

	return nil
}

// GetAll retrieves all chat messages from database
func (c *chatMessageRepository) GetLastPublicMessages(limit int) ([]models.ChatMessage, error) {
	var chatMessages []models.ChatMessage

	if err := c.db.Model(&chatMessages).
		Column("chat_message.*").
		Where("receiver_id IS NULL").
		Order("id desc").
		Relation("Receiver").
		Relation("User").
		Limit(limit).
		Select(); err != nil {
		if err ==pg.ErrNoRows {
			return []models.ChatMessage{}, nil
		}

		return nil, err
	}

	return chatMessages, nil
}

func (c *chatMessageRepository) GetLastPrivateMessages(user models.User, limit int) ([]models.ChatMessage, error) {
	var chatMessages []models.ChatMessage

	if err := c.db.Model(&chatMessages).
		Column("chat_message.*").
		Where("user_id=? and (receiver_id > 0 or receiver_id = ?)", user.ID, user.ID).
		Order("id desc").
		Relation("Receiver").
		Relation("User").
		Limit(limit).
		Select(); err != nil {
		if err ==pg.ErrNoRows {
			return []models.ChatMessage{}, nil
		}

		return nil, err
	}

	return chatMessages, nil
}
