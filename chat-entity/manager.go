package chatEntity

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rishavmngo/chat-backend-v2/entities"
	intrf "github.com/rishavmngo/chat-backend-v2/interfac"
	jwts "github.com/rishavmngo/chat-backend-v2/jwt"
	"github.com/rishavmngo/chat-backend-v2/responses"
	"github.com/rishavmngo/chat-backend-v2/utils"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

type Manager struct {
	clients  ClientList
	DB       intrf.Store
	handlers map[string]EventHandler
	sync.RWMutex
}

type ManagerController interface {
	EstablishWs()
}

func NewManager(db intrf.Store) *Manager {

	m := &Manager{
		clients:  make(ClientList),
		DB:       db,
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()

	return m

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type httpHandler func(w http.ResponseWriter, r *http.Request)

func (m *Manager) EstablishWs(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	resp := mux.Vars(r)["token"]

	log.Println("here is the token: ", resp)

	uid, err := jwts.ExtractTokenIDTest(resp)

	if err != nil {
		log.Println(err)
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected with uid: ", uid)

	client := NewClient(uid, m, conn)
	m.addClient(client)
	log.Println(len(m.clients))

	go client.Reader()
	go client.Writer()

}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client.UID] = client
}

func (m *Manager) RemoveClient(uid uint) {
	m.Lock()
	defer m.Unlock()

	delete(m.clients, uid)
}

func (m *Manager) setupEventHandlers() {
	m.handlers[sendMessageEvent] = func(event Event, c *Client) error {

		return nil
	}

	m.handlers[getAllChatsEvent] = func(event Event, c *Client) error {

		var chats []*entities.Chats
		chats, err := m.DB.GetAllChatsByUID(c.UID)

		if err != nil {
			log.Println("error while fetching all chats")
			return err
		}

		data, err := json.Marshal(chats)

		if err != nil {
			log.Println("Error marshaling data")
			return err
		}

		var outgoingEvent Event
		outgoingEvent.Type = getAllChatsEvent
		outgoingEvent.Payload = data

		c.egress <- outgoingEvent
		return nil
	}
	m.handlers[getAllMsgByChatID] = func(event Event, c *Client) error {
		var messages []*entities.Message
		messages = make([]*entities.Message, 0)
		data, err := utils.UnmarshalUnstructredData(event.Payload)

		if err != nil {
			return err
		}
		id, ok := data["chat_id"].(float64)

		if !ok {
			log.Println("error at mapping chat_id to float64")
			return errors.New("error at mapping chat_id to float64")
		}

		messages, err = m.DB.GetAllMsgByChatID(uint(id))

		if err != nil {
			log.Println("Error at getting msg by ChatID")
			return err
		}

		if len(messages) == 0 {
			messages = make([]*entities.Message, 0)
		}
		resp, err := json.Marshal(messages)

		if err != nil {
			log.Println("Error marshaling data")
			return err
		}
		var outgoingEvent Event
		outgoingEvent.Type = getAllMsgByChatID
		outgoingEvent.Payload = resp

		c.egress <- outgoingEvent
		return nil
	}
	m.handlers[sendMessagePChatEvent] = func(event Event, c *Client) error {
		var message entities.Message

		err := json.Unmarshal(event.Payload, &message)

		if err != nil {
			log.Println(err)
			return err
		}
		otherId, err := m.DB.GetOtherMember(message.ChatID, message.AuthorID)

		if err != nil {
			log.Println(err)
			return err
		}

		err = m.DB.AddMsgToPChat(&message)

		if err != nil {
			log.Println(err)
			return err
		}

		messagesArr, err := m.DB.GetAllMsgByChatID(message.ChatID)

		if err != nil {
			log.Println(err)
			return err
		}

		resp, err := json.Marshal(messagesArr)

		if err != nil {
			log.Println("Error marshaling data")
			return err
		}
		var outgoingEvent Event
		outgoingEvent.Type = getAllMsgByChatID
		outgoingEvent.Payload = resp

		client, ok := m.clients[otherId]

		if ok {

			client.egress <- outgoingEvent
		}

		c.egress <- outgoingEvent

		return nil
	}
}

func (m *Manager) routeEvent(event Event, client *Client) error {

	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}
