package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const (
	writeTimeout   = time.Second * 10
	readTimeout    = time.Second * 10
	pingPeriod     = time.Second * 5
	maxMessageSize = 512
)

type connect struct {
	wsConn    *websocket.Conn
	send      chan []byte
	waitGruop sync.WaitGroup
}

func newConn(wsConn *websocket.Conn) *connect {
	return &connect{
		wsConn: wsConn,
		send:   make(chan []byte),
	}
}

func (c *connect) readPump() {
	defer c.waitGruop.Done()

	c.wsConn.SetReadLimit(maxMessageSize)
	c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
	c.wsConn.SetPongHandler(func(string) error {
		log.Println("received pong")
		c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		typ, msg, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Println("reading error:", err)
			close(c.send)
			return
		}

		if typ != websocket.TextMessage {
			continue
		}

		c.send <- msg
	}
}

func (c *connect) writePump() {
	defer c.waitGruop.Done()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case s, more := <-c.send:
			if !more {
				return
			}

			log.Println("sending:", string(s))
			c.wsConn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.wsConn.WriteMessage(websocket.TextMessage, s); err != nil {
				log.Println("writting error:", err)
				return
			}

		case <-ticker.C:
			log.Println("sent ping")
			c.wsConn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeTimeout))
		}
	}
}

func (c *connect) run() {
	c.waitGruop.Add(2)

	go c.readPump()
	go c.writePump()

	c.waitGruop.Wait()
	c.wsConn.Close()
}
