package messengerWebsocket

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Conn struct {
	websocketConn *websocket.Conn
	send          chan []byte
	waitGroup     sync.WaitGroup
}

func NewConn(c *websocket.Conn) Conn {
	return Conn {
		websocketConn: c,
		send:          make(chan []byte),
	}
}

func (c *Conn) Run() {
	c.waitGroup.Add(2)
	go c.readPump()
	go c.writePump()
	c.waitGroup.Wait()
}

func (c *Conn) readPump() {
	defer c.waitGroup.Done()

	for {
		readType, msg, err := c.websocketConn.ReadMessage()

		if err != nil {
			log.Println("error in reading : ", err)
			return
		}

		if readType != websocket.TextMessage {
			log.Println("not a text message")
			continue
		}

		log.Println("echoing : ", string(msg))

		c.send <- msg
	}
}

func (c *Conn) writePump() {
	defer c.waitGroup.Done()

	for {
		select {
		case msg, more := <-c.send:
			if !more {
				return
			}

			if err := c.websocketConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("error writing", err)
				return
			}
		}
	}
}
