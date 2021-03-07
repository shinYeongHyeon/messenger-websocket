package messengerWebsocket

import (
	"github.com/gorilla/websocket"
	"github.com/shinYeongHyeon/messenger-websocket/chat"
	"github.com/shinYeongHyeon/messenger-websocket/db"
	"log"
	"sync"
	"time"
)

const (
	writeTimeout   = 10 * time.Second
	readTimeout    = 60 * time.Second
	pingPeriod     = 10 * time.Second
	maxMessageSize = 512
)

type conn struct {
	websocketConn     *websocket.Conn
	wg         sync.WaitGroup
	sub        db.ChatroomSubscription
	chatroomID int
	senderID   int
}

func newConn(websocketConn *websocket.Conn, chatroomID, senderID int) *conn {
	return &conn{
		websocketConn:     websocketConn,
		chatroomID: chatroomID,
		senderID:   senderID,
	}
}

func (c *conn) run() error {
	sub, err := db.NewChatroomSubscription(c.chatroomID)
	if err != nil {
		return err
	}
	c.sub = sub

	c.wg.Add(2)
	go c.readPump()
	go c.writePump()

	c.wg.Wait()
	return nil
}

func (c *conn) readPump() {
	defer c.wg.Done()
	defer c.sub.Close()

	c.websocketConn.SetReadLimit(maxMessageSize)
	c.websocketConn.SetReadDeadline(time.Now().Add(readTimeout))
	c.websocketConn.SetPongHandler(func(string) error {
		c.websocketConn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		var msg chat.Message
		if err := c.websocketConn.ReadJSON(&msg); err != nil {
			log.Println("err reading:", err)
			return
		}

		db.SendMessage(c.senderID, c.chatroomID, msg.Text)
	}
}

func (c *conn) writePump() {
	defer c.wg.Done()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case s, more := <-c.sub.C:
			if !more {
				return
			}
			c.websocketConn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.websocketConn.WriteJSON(s); err != nil {
				log.Println("err writing:", err)
				return
			}
		case <-ticker.C:
			c.websocketConn.WriteControl(
				websocket.PingMessage, nil, time.Now().Add(writeTimeout))
		}
	}
}
