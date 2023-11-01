package chatEntity

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var (
	sendMessageEvent      = "sendMessage"
	getAllChatsEvent      = "getAllChats"
	getAllMsgByChatID     = "getAllMsgByChatId"
	sendMessagePChatEvent = "sendMessagePChat"
)

type EventHandler func(event Event, c *Client) error
