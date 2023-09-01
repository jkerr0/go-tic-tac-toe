package handlers

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Channel struct {
	connections []*websocket.Conn
	writerMutex sync.Mutex
	messages    []string
}

func (c *Channel) Join(connection *websocket.Conn) {
	c.connections = append(c.connections, connection)
}

func (c *Channel) Leave(connection *websocket.Conn) error {
	rmIndex := -1
	for i, conn := range c.connections {
		if conn == connection {
			rmIndex = i
			break
		}
	}
	if rmIndex == -1 {
		return errors.New("connection not found")
	}
	c.connections = append(c.connections[:rmIndex], c.connections[rmIndex+1:]...)
	return nil
}

func (c *Channel) Broadcast(message string) {
	c.messages = append(c.messages, message)
	c.writerMutex.Unlock()
}

func OpenChannel() *Channel {
	c := &Channel{}
	go c.writer()
	return c
}

func (c *Channel) writer() {
	for {
		c.writerMutex.Lock()
		for _, message := range c.messages {
			for _, ws := range c.connections {
				if ws == nil {
					continue
				}
				err := ws.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					c.Leave(ws)
				}
			}
		}
		c.messages = make([]string, 0)
	}
}
