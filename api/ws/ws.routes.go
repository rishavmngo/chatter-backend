package ws

import (
	"github.com/gorilla/mux"
	chatEntity "github.com/rishavmngo/chat-backend-v2/chat-entity"
)

func Routes(router *mux.Router, m *chatEntity.Manager) {

	router.HandleFunc("/connection/{token}", m.EstablishWs)

}
