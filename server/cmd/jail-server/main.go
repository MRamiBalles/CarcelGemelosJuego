// Package main is the entry point for the Cárcel de los Gemelos game server.
// It only handles dependency injection and server initialization.
// NO business logic belongs here.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/engine"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/network"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/gorilla/websocket"
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

	appLogger.Info("Bootstrapping WebSocket Hub...")
	hub := network.NewHub(appLogger)
	go hub.Run(ctx)
	hub.StartEventPoller(ctx, eventLog)

	// Setup API Routes
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, appLogger)
	})

	http.HandleFunc("/api/trigger-oracle", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type requestParams struct {
			Target  string `json:"target"`
			Message string `json:"message"`
		}
		var req requestParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		eventLog.StoreEvent(events.GameEvent{
			Type:     events.EventTypeOraclePainfulTruth,
			ActorID:  "Audience",
			TargetID: req.Target,
			Payload:  map[string]string{"message": "The Oracle whispers a painful truth: " + req.Message},
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Oracle event triggered"})
	})

	http.HandleFunc("/api/trigger-torture", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type requestParams struct {
			SoundID string `json:"soundId"`
		}
		var req requestParams
		_ = json.NewDecoder(r.Body).Decode(&req)

		// Create audio torture payload
		payload := events.AudioTorturePayload{
			SoundName: req.SoundID,
			Duration:  30, // 30 in-game minutes of torture
		}

		eventLog.StoreEvent(events.GameEvent{
			Type:     events.EventTypeAudioTorture,
			ActorID:  "Audience",
			TargetID: "ALL",
			Payload:  payload,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Audio Torture dispatched"})
	})

	go func() {
		log.Println("[JAIL-SERVER] HTTP API & WS Server listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("[JAIL-SERVER] Server running. Press Ctrl+C to exit.")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[JAIL-SERVER] Shutting down...")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow cross-origin requests for Next.js dev server
	},
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *network.Hub, w http.ResponseWriter, r *http.Request, log *logger.Logger) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Failed to upgrade websocket connection")
		return
	}

	client := network.NewClient(hub, conn)
	client.Register()

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
