package main

import (
	server "github.com/rishavmngo/chat-backend-v2/api"
	"github.com/rishavmngo/chat-backend-v2/postgres"
)

func main() {

	var s server.Server

	store := postgres.InitilizePostgresStore("postgres", "password", "chatter", "5432")

	s.Initilize(":3001", store)
	s.Run()
}
