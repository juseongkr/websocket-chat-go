package chat

import (
	"time"
)

type Message struct {
	Text     string    `json:"text"`
	Sender   string    `json:"sender"`
	SenderId int       `json:"senderId"`
	SentOn   time.Time `json:sentOn"`
}

type Room struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
