package models

import "time"

type ChatMessage struct {
	ID       int64     `json:"id"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Text     string    `json:"text"`
	DateTime time.Time `json:"date_time"`
}
