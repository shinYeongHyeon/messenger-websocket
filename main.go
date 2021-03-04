package main

import (
	"github.com/gorilla/websocket"
	"github.com/shinYeongHyeon/messenger-websocket/messengerWebsocket"
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

func main() {
	error := http.ListenAndServe(":8081", http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			conn, err := webSocketUpgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println("err : ", err)
				return
			}

			c := messengerWebsocket.NewConn(conn)
			c.Run()
		}))

	if error != nil {
		log.Println("Listen Err")
	}
}