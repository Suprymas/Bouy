package main

import (
	"buoy-hub/internal/server"
	"log"
	"net/http"
)

func main() {
	s := server.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.HandleConnection)
	mux.HandleFunc("/health", server.HealthHandler)

	addr := ":8080"
	log.Printf("Buoy hub starting on ws://localhost%s/ws", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
