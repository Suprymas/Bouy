package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	id   string
	conn *websocket.Conn
}

// Server manages WebSocket clients
type Server struct {
	clients map[string]*Client
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]*Client),
	}
}

// upgrader converts a regular HTTP connection into a WebSocket connection.
// CheckOrigin returning true disables origin checks — fixes the 403 error.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins
	},
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	clientID := fmt.Sprintf("client-%d", time.Now().UnixNano())
	client := &Client{id: clientID, conn: ws}
	s.clients[clientID] = client

	log.Printf("[+] New client connected: %s (total: %d)", clientID, len(s.clients))

	ws.WriteMessage(websocket.TextMessage, []byte("Welcome! Your ID is: "+clientID))

	defer func() {
		delete(s.clients, clientID)
		ws.Close()
		log.Printf("[-] Client disconnected: %s (total: %d)", clientID, len(s.clients))
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error receiving from %s: %v", clientID, err)
			}
			break
		}

		log.Printf("[MSG] %s @ %s: %s", clientID, time.Now().Format("15:04:05"), string(message))

		reply := fmt.Sprintf("[%s] Echo: %s", time.Now().Format("15:04:05"), string(message))
		if err := ws.WriteMessage(websocket.TextMessage, []byte(reply)); err != nil {
			log.Printf("Error sending reply to %s: %v", clientID, err)
			break
		}
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok"}`)
}

func main() {
	server := NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", server.handleConnection)
	mux.HandleFunc("/health", healthHandler)

	addr := ":8080"
	log.Printf("WebSocket server starting on ws://localhost%s/ws", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}


