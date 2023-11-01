package entities

import "time"

type Message struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	ChatID     uint      `json:"chat_id"`
	CreatedAt  time.Time `json:"message_created_at"`
	AuthorID   uint      `json:"author_id"`
	AuthorName string    `json:"author_name"`
}
