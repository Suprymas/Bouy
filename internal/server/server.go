package server

import (
	"buoy-hub/internal/client"
	"buoy-hub/internal/db"
	"buoy-hub/internal/storage"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Server manages all connected buoy clients and database access
type Server struct {
	clients map[string]*client.Client
	db      *db.DB
	storage *storage.Storage
	latest  map[string]buoyState
	links   map[string]string
	logs    []LogEntry
	mu      sync.RWMutex
}

func New(database *db.DB, store *storage.Storage) *Server {
	return &Server{
		clients: make(map[string]*client.Client),
		db:      database,
		storage: store,
		latest:  make(map[string]buoyState),
		links:   make(map[string]string),
	}
}

// upgrader upgrades HTTP connections to WebSocket.
// CheckOrigin returns true to allow all origins.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
