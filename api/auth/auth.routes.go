package auth

import (
	"github.com/gorilla/mux"
	"github.com/rishavmngo/chat-backend-v2/entities"
	intrf "github.com/rishavmngo/chat-backend-v2/interfac"
	middlewares "github.com/rishavmngo/chat-backend-v2/middleware"
)

func Routes(router *mux.Router, db intrf.Store) {
	var user entities.UserController

	user = &controller{db}

	router.HandleFunc("/login", middlewares.SetMiddlewareJSON(user.Login)).Methods("POST")
	router.HandleFunc("/register", user.Register).Methods("POST")
	router.HandleFunc("/userByUID", middlewares.SetMiddlewareJSON(user.GetUserByUID)).Methods("GET")
}
