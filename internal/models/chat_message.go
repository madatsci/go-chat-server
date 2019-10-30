package models

import "time"

type ChatMessage struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	User       *User     `json:"user"`
	ReceiverID int64     `json:"receiver_id"`
	Receiver   *User     `json:"receiver"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
