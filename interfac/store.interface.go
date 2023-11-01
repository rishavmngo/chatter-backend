package intrf

import "github.com/rishavmngo/chat-backend-v2/entities"

type Store interface {
	GetUser(*entities.User) error
	GetAllChatsByUID(UID uint) ([]*entities.Chats, error)
	GetUserByUID(*entities.User) error
	GetAllMsgByChatID(ChatID uint) ([]*entities.Message, error)
	GetOtherMember(chat_id, partner_id uint) (uint, error)
	AddMsgToPChat(message *entities.Message) error
}
