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
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/engine"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/infra/ai"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/infra/storage"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/network"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins"
	"github.com/gorilla/websocket"
)

// SQLitePersisterAdapter translates domain events to storage events.
type SQLitePersisterAdapter struct {
	repo *storage.SQLiteEventRepository
}

func (a *SQLitePersisterAdapter) Append(event events.GameEvent) error {
	payloadBytes, _ := json.Marshal(event.Payload)
	var payloadMap map[string]interface{}
	json.Unmarshal(payloadBytes, &payloadMap)

	storageEvent := storage.GameEvent{
		ID:         event.ID,
		GameID:     "GAME_1", // Default singleton game ID
		Timestamp:  event.Timestamp,
		EventType:  string(event.Type),
		ActorID:    event.ActorID,
		TargetID:   event.TargetID,
		Payload:    payloadMap,
		GameDay:    event.GameDay,
		IsRevealed: event.IsRevealed,
	}
	return a.repo.Append(context.Background(), storageEvent)
}

func bootstrapPrisoners(ctx context.Context, repo *storage.SQLiteSnapshotRepository, eng *engine.Engine, appLogger *logger.Logger) {
	appLogger.Info("Checking DB for existing prisoners...")
	snaps, err := repo.GetByGameID(ctx, "GAME_1")
	if err != nil {
		appLogger.Error("Failed to query DB for prisoners: " + err.Error())
		return
	}

	if len(snaps) == 0 {
		appLogger.Info("Database empty. Seeding Initial 6 Prisoners...")
		starters := []*prisoner.Prisoner{
			prisoner.NewPrisoner("P001", "Frank", prisoner.ArchetypeVeteran),
			prisoner.NewPrisoner("P002", "Dakota", prisoner.ArchetypeExplosive),
			prisoner.NewPrisoner("P003", "Aída", prisoner.ArchetypeChaos),
			prisoner.NewPrisoner("P004", "Tartaria", prisoner.ArchetypeMystic),
			prisoner.NewPrisoner("P005", "Héctor", prisoner.ArchetypeDeceiver),
			prisoner.NewPrisoner("P006", "Ylenia", prisoner.ArchetypeToxic),
		}
		for _, p := range starters {
			snap := storage.PrisonerSnapshot{
				PrisonerID: p.ID,
				GameID:     "GAME_1",
				Name:       p.Name,
				Archetype:  string(p.Archetype),
				IsIsolated: p.IsIsolated,
				Sanity:     p.Sanity,
				Dignity:    p.Dignity,
			}
			repo.Upsert(ctx, snap)
			eng.RegisterPrisoner(p)
		}
	} else {
		appLogger.Info("Reconstructing Prisoners from SQLite State...")
		for _, snap := range snaps {
			p := prisoner.NewPrisoner(snap.PrisonerID, snap.Name, prisoner.Archetype(snap.Archetype))
			p.Sanity = snap.Sanity
			p.Dignity = snap.Dignity
			p.IsIsolated = snap.IsIsolated
			eng.RegisterPrisoner(p)
		}
	}
}

func main() {
	log.Println("[JAIL-SERVER] Initializing 'Cárcel de los Gemelos' Authoritative Server...")

	appLogger := logger.NewLogger()

	appLogger.Info("Initializing SQLite database 'jail.db'...")
	db, err := storage.InitSQLite("jail.db")
	if err != nil {
		appLogger.Error("Failed to initialize SQLite: " + err.Error())
		os.Exit(1)
	}
	eventRepo := storage.NewSQLiteEventRepository(db)
	eventPersister := &SQLitePersisterAdapter{repo: eventRepo}

	appLogger.Info("Bootstrapping EventLog...")
	eventLog := events.NewEventLog(eventPersister)

	appLogger.Info("Bootstrapping Engine Subsystems...")
	gameEngine := engine.NewEngine(eventLog, appLogger)

	// Start engine processing in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	snapRepo := storage.NewSQLiteSnapshotRepository(db)
	bootstrapPrisoners(ctx, snapRepo, gameEngine, appLogger)

	// Attempt to recover the last known game clock state
	var tickPayloadStr string
	err = db.QueryRowContext(ctx, "SELECT payload FROM events WHERE game_id = ? AND event_type = ? ORDER BY timestamp DESC LIMIT 1", "GAME_1", events.EventTypeTimeTick).Scan(&tickPayloadStr)
	if err == nil {
		var tickPayload engine.TimeTickPayload
		if err := json.Unmarshal([]byte(tickPayloadStr), &tickPayload); err == nil {
			gameEngine.OverrideTime(tickPayload.GameDay, tickPayload.GameHour, tickPayload.TickNumber)
			appLogger.Info("Restored Game Clock from Database.")
		}
	}

	gameEngine.Start(ctx)

	// Automated State Backup Routine
	go func() {
		backupTicker := time.NewTicker(5 * time.Second)
		defer backupTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-backupTicker.C:
				for _, p := range gameEngine.GetPrisoners() {
					snap := storage.PrisonerSnapshot{
						PrisonerID: p.ID,
						GameID:     "GAME_1",
						Name:       p.Name,
						Archetype:  string(p.Archetype),
						IsIsolated: p.IsIsolated,
						Sanity:     p.Sanity,
						Dignity:    p.Dignity,
					}
					_ = snapRepo.Upsert(ctx, snap)
				}
			}
		}
	}()

	appLogger.Info("Bootstrapping WebSocket Hub...")
	hub := network.NewHub(appLogger)
	go hub.Run(ctx)
	hub.StartEventPoller(ctx, eventLog)

	appLogger.Info("Bootstrapping AI Cognition (Los Gemelos)...")
	budgetGate := ai.NewBudgetGate(10.0, 50.0) // $10/day, $50/month safety net
	llmProvider := ai.NewOpenAIProvider(budgetGate)

	aiMind := twins.NewTwinsMind(eventLog, gameEngine, llmProvider, gameEngine.GetNoiseManager(), hub, appLogger)
	aiMind.SetGameID("GAME_1")
	go aiMind.Start(ctx)

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

		eventLog.Append(events.GameEvent{
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

		eventLog.Append(events.GameEvent{
			Type:     events.EventTypeAudioTorture,
			ActorID:  "Audience",
			TargetID: "ALL",
			Payload:  payload,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Audio Torture dispatched"})
	})

	http.HandleFunc("/api/twins/force-decision", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		decision, err := aiMind.ForceDecision(ctx, 1) // default to day 1 for testing
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "ok",
			"decision": decision,
		})
	})

	http.HandleFunc("/api/poll/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload engine.PollPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		if payload.PollID == "" || len(payload.Options) < 2 {
			http.Error(w, "Invalid poll definition", http.StatusBadRequest)
			return
		}

		eventLog.Append(events.GameEvent{
			ID:        events.GenerateEventID(),
			Timestamp: time.Now(),
			Type:      events.EventTypePollCreated,
			ActorID:   "SYSTEM_ADMIN",
			TargetID:  "ALL",
			Payload:   payload,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Poll created successfully"})
	})

	http.HandleFunc("/api/poll/vote", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		type voteReq struct {
			PollID string `json:"poll_id"`
			Option string `json:"option"`
		}
		var req voteReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		gameEngine.GetPollingSystem().CastVote(req.PollID, req.Option)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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
