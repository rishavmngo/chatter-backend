package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rishavmngo/chat-backend-v2/api/auth"
	"github.com/rishavmngo/chat-backend-v2/api/ws"
	chatEntity "github.com/rishavmngo/chat-backend-v2/chat-entity"
	intrf "github.com/rishavmngo/chat-backend-v2/interfac"
)

type Server struct {
	DB      intrf.Store
	Port    string
	Router  *mux.Router
	Manager *chatEntity.Manager
}

func (s *Server) Initilize(port string, db intrf.Store) {
	s.DB = db
	s.Port = port
	s.Router = mux.NewRouter()
	s.InitilizeRoutes()
}

func (s *Server) InitilizeRoutes() {

	log.Println("routes registered")
	s.Manager = chatEntity.NewManager(s.DB)
	s.Subroute("/auth", auth.Routes)
	s.SubrouteWs("/chat", ws.Routes)

}

type InitRouterType func(*mux.Router, intrf.Store)
type InitRouterTypeWs func(*mux.Router, *chatEntity.Manager)

func (server *Server) Subroute(path string, initRouter InitRouterType) {
	subrouter := server.Router.PathPrefix(path).Subrouter()
	initRouter(subrouter, server.DB)
}

func (server *Server) SubrouteWs(path string, initRouter InitRouterTypeWs) {
	subrouter := server.Router.PathPrefix(path).Subrouter()
	initRouter(subrouter, server.Manager)
}

func (s *Server) Run() {

	http.ListenAndServe(s.Port, s.Router)

}
