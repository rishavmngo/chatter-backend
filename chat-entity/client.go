package chatEntity

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	UID     uint
	Manager *Manager
	Conn    *websocket.Conn
	egress  chan Event
}

type ClientList map[uint]*Client

func NewClient(uid uint, manager *Manager, conn *websocket.Conn) *Client {

	return &Client{
		Manager: manager,
		UID:     uid,
		Conn:    conn,
		egress:  make(chan Event),
	}
}

func (client *Client) Reader() {
	defer func() {

		log.Println("removing client ", client.UID)

		client.Manager.RemoveClient(client.UID)

	}()

	for {

		_, payload, err := client.Conn.ReadMessage()

		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		var event Event

		err = json.Unmarshal(payload, &event)

		if err != nil {
			log.Printf("Event structure invalid")
			break
		}
		log.Println("event: ", event)
		client.Manager.routeEvent(event, client)
	}
}

func (c *Client) Writer() {
	defer func() {
		log.Println("removing client ", c.UID)
		c.Manager.RemoveClient(c.UID)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {

				if err := c.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed: ", err)
				}
				return

			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
		}
	}

}
