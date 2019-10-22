package ws

import "time"

type Message struct {
	From     string    `json:"from"`
	To       string    `json:"to"`
	Text     string    `json:"text"`
	DateTime time.Time `json:"date_time"`
}
