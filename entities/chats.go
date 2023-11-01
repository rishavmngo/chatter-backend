package entities

import "time"

type Chats struct {
	ID           uint      `json:"chat_id"`
	CreatedAt    time.Time `json:"created_at"`
	LastActiveAt time.Time `json:"last_active_at"`
	PartnerName  string    `json:"partner_name"`
	PartnerUID   uint      `json:"partner_uid"`
}
