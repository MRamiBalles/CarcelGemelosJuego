// Package main is the entry point for the Cárcel de los Gemelos game server.
// It only handles dependency injection and server initialization.
// NO business logic belongs here.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/engine"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

func main() {
	log.Println("[JAIL-SERVER] Initializing 'Cárcel de los Gemelos' Authoritative Server...")

	appLogger := logger.NewLogger()
	appLogger.Info("Bootstrapping EventLog...")
	eventLog := events.NewEventLog()

	appLogger.Info("Bootstrapping Engine Subsystems...")
	gameEngine := engine.NewEngine(eventLog, appLogger)

	// Start engine processing in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gameEngine.Start(ctx)

	// TODO: Start WebSocket listener (T013)

	log.Println("[JAIL-SERVER] Server running. Press Ctrl+C to exit.")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[JAIL-SERVER] Shutting down...")
}
