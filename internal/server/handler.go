package server

import (
	"buoy-hub/internal/client"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// HandleConnection upgrades the HTTP request to a WebSocket
// and starts listening for messages from the buoy.
func (s *Server) HandleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	clientID := fmt.Sprintf("buoy-%d", time.Now().UnixNano())
	c := client.New(clientID, ws)
	s.clients[clientID] = c

	log.Printf("[+] Buoy connected: %s (total: %d)", clientID, len(s.clients))
	ws.WriteMessage(websocket.TextMessage, []byte("Welcome! Your ID is: "+clientID))

	defer s.disconnect(clientID, ws)

	s.readLoop(c)
}

// readLoop continuously reads messages from a connected buoy.
func (s *Server) readLoop(c *client.Client) {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Read error from %s: %v", c.ID, err)
			}
			break
		}

		log.Printf("[MSG] %s @ %s: %s", c.ID, time.Now().Format("15:04:05"), string(message))
		s.handleMessage(c, message)
	}
}

// handleMessage processes an incoming message from a buoy.
// This is where you'll add GPS parsing, image saving, DB writes, etc.
func (s *Server) handleMessage(c *client.Client, message []byte) {
	reply := fmt.Sprintf("[%s] Echo: %s", time.Now().Format("15:04:05"), string(message))
	if err := c.Conn.WriteMessage(websocket.TextMessage, []byte(reply)); err != nil {
		log.Printf("Write error to %s: %v", c.ID, err)
	}
}

// disconnect cleans up a buoy connection.
func (s *Server) disconnect(clientID string, ws *websocket.Conn) {
	delete(s.clients, clientID)
	ws.Close()
	log.Printf("[-] Buoy disconnected: %s (total: %d)", clientID, len(s.clients))
}

// HealthHandler returns a simple JSON health check.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}
