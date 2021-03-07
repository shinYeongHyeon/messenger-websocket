package main

import (
	"github.com/gorilla/websocket"
	"github.com/shinYeongHyeon/messenger-websocket/api"
	"github.com/shinYeongHyeon/messenger-websocket/db"
	"log"
	"net/http"
)

var webSocketUpgrader = websocket.Upgrader {
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: 	 func(r *http.Request) bool {
		return true
	},
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	if err := db.Connect("host=localhost port=5432 user=postgres dbname=go-chat sslmode=disable"); err != nil {
		log.Fatalln(err)
	}

	if err := http.ListenAndServe(":8080", api.Handler()); err != nil {
		log.Fatalln(err)
	}
}
