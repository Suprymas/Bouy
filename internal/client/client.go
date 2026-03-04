package client

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket buoy
type Client struct {
	ID   string
	Role string
	Conn *websocket.Conn
	mu   sync.Mutex
}

func New(id string, role string, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Role: role,
		Conn: conn,
	}
}

func (c *Client) Write(messageType int, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Conn.WriteMessage(messageType, payload)
}
