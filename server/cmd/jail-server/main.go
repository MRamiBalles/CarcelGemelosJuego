// Package main is the entry point for the Cárcel de los Gemelos game server.
// It only handles dependency injection and server initialization.
// NO business logic belongs here.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("[JAIL-SERVER] Initializing 'Cárcel de los Gemelos' Authoritative Server...")

	// TODO: Wire up dependencies (T002+)
	// - IdentityService
	// - SocialEngine
	// - EnvironmentService (The Twins)
	// - EventLog (The VAR)
	// - Network (WebSockets)

	// TODO: Start game loop ticker (T007)
	// TODO: Start WebSocket listener (T013)

	log.Println("[JAIL-SERVER] Server running. Press Ctrl+C to exit.")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[JAIL-SERVER] Shutting down...")
}
