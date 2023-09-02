package handlers

import (
	"errors"

	"github.com/gorilla/websocket"
)

type ChannelPool struct {
	pool map[int]*Channel
}

func NewChannelPool() *ChannelPool {
	return &ChannelPool{pool: map[int]*Channel{}}
}

func (cp *ChannelPool) Join(connection *websocket.Conn, id int) *Channel {
	if cp.pool[id] == nil {
		cp.pool[id] = OpenChannel()
	}
	cp.pool[id].Join(connection)
	return cp.pool[id]
}

func (cp *ChannelPool) Leave(connection *websocket.Conn, id int) {
	cp.pool[id].Leave(connection)
	if cp.pool[id].IsEmpty() {
		cp.pool[id].Close()
		cp.pool[id] = nil
	}
}

type Channel struct {
	connections []*websocket.Conn
	notify      chan bool
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
	c.notify <- true
}

func OpenChannel() *Channel {
	c := &Channel{
		notify: make(chan bool)}
	go c.writer()
	return c
}

func (c *Channel) Close() {
	c.notify <- false
	close(c.notify)
}

func (c *Channel) IsEmpty() bool {
	return len(c.connections) == 0
}

func (c *Channel) writer() {
	for {
		shouldRun := <-c.notify
		if !shouldRun {
			break
		}
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
