package server

import (
	"buoy-hub/internal/client"
	"net/http"

	"github.com/gorilla/websocket"
)

// Server manages all connected buoy clients
type Server struct {
	clients map[string]*client.Client
}

func New() *Server {
	return &Server{
		clients: make(map[string]*client.Client),
	}
}

// upgrader upgrades HTTP connections to WebSocket.
// CheckOrigin returns true to allow all origins.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
