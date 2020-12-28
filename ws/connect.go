package ws

import (
	"github.com/gorilla/websocket"
	"github.com/juseongkr/websocket-chat-go/chat"
	"github.com/juseongkr/websocket-chat-go/db"
	"log"
	"sync"
	"time"
)

const (
	writeTimeout   = time.Second * 10
	readTimeout    = time.Second * 60
	pingPeriod     = time.Second * 10
	maxMessageSize = 512
)

type connect struct {
	wsConn     *websocket.Conn
	sub        db.ChatroomSubscription
	senderId   int
	chatroomId int
	waitGruop  sync.WaitGroup
}

func newConn(wsConn *websocket.Conn, chatroomId, senderId int) *connect {
	return &connect{
		wsConn:     wsConn,
		senderId:   senderId,
		chatroomId: chatroomId,
	}
}

func (c *connect) readPump() {
	defer c.waitGruop.Done()
	defer c.sub.Close()

	c.wsConn.SetReadLimit(maxMessageSize)
	c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
	c.wsConn.SetPongHandler(func(string) error {
		c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		var msg chat.Message
		if err := c.wsConn.ReadJSON(&msg); err != nil {
			log.Println("reading error:", err)
			return
		}

		db.SendMessage(c.senderId, c.chatroomId, msg.Text)
	}
}

func (c *connect) writePump() {
	defer c.waitGruop.Done()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case s, more := <-c.sub.C:
			if !more {
				return
			}

			c.wsConn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.wsConn.WriteJSON(s); err != nil {
				log.Println("writting error:", err)
				return
			}

		case <-ticker.C:
			c.wsConn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeTimeout))
		}
	}
}

func (c *connect) run() error {
	sub, err := db.NewChatroomSubscription(c.chatroomId)
	if err != nil {
		return err
	}
	c.sub = sub
	c.waitGruop.Add(2)

	go c.readPump()
	go c.writePump()

	c.waitGruop.Wait()

	return nil
}
