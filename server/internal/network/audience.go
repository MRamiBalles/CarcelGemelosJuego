// Package network - audience.go
// T014: AudienceBridge - REST API for mobile "Pay-to-Torture" integration.
//
// The Audience is the Third Twin. They can spend "Sadism Points" to
// trigger noise events, reveal secrets, or manipulate resources.
package network

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/engine"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// AudienceBridge handles mobile app interactions.
type AudienceBridge struct {
	noiseManager *engine.NoiseManager
	eventLog     *events.EventLog
	wsHub        *Hub
	logger       *logger.Logger
}

// NewAudienceBridge creates a new audience interaction handler.
func NewAudienceBridge(nm *engine.NoiseManager, el *events.EventLog, hub *Hub, log *logger.Logger) *AudienceBridge {
	return &AudienceBridge{
		noiseManager: nm,
		eventLog:     el,
		wsHub:        hub,
		logger:       log,
	}
}

// TorturRequest is the payload for triggering a noise event.
type TortureRequest struct {
	GameID      string `json:"game_id"`
	TargetZone  string `json:"target_zone"`  // "ALL", "BLOCK_A", or prisoner ID
	NoiseType   string `json:"noise_type"`   // "SIREN", "CRYING_BABY", etc.
	Intensity   int    `json:"intensity"`    // 1-3
	SadismCost  int    `json:"sadism_cost"`  // Points spent by audience
	AudienceID  string `json:"audience_id"`  // Who triggered it
}

// RevealRequest is the payload for revealing a secret event.
type RevealRequest struct {
	GameID     string `json:"game_id"`
	EventID    string `json:"event_id"`
	SadismCost int    `json:"sadism_cost"`
	AudienceID string `json:"audience_id"`
}

// HandleTorture is the endpoint for audience-triggered noise events.
// POST /api/audience/torture
func (ab *AudienceBridge) HandleTorture(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ab.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TortureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ab.jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate
	if req.GameID == "" || req.SadismCost < 10 {
		ab.jsonError(w, "Invalid game_id or insufficient sadism_cost", http.StatusBadRequest)
		return
	}

	// TODO: Verify audience has enough Sadism Points (payment integration)

	// Trigger the torture via NoiseManager
	reason := "AUDIENCE_VOTE:" + req.AudienceID
	ab.noiseManager.TriggerPunishment(req.TargetZone, reason)

	// Log the audience action
	ab.logger.Event("AUDIENCE_TORTURE", req.AudienceID, 
		"Target:"+req.TargetZone+" Cost:"+string(rune('0'+req.SadismCost)))

	// Notify connected players
	ab.wsHub.BroadcastToGame(req.GameID, Message{
		Type:      MsgTypeEvent,
		Timestamp: time.Now().Unix(),
		Payload: map[string]interface{}{
			"event_type": "AUDIENCE_INTERVENTION",
			"message":    "La Audiencia ha hablado... Los Gemelos escuchan.",
		},
	})

	ab.jsonSuccess(w, map[string]interface{}{
		"success":      true,
		"message":      "Torture triggered",
		"sadism_spent": req.SadismCost,
	})
}

// HandleReveal is the endpoint for revealing hidden events.
// POST /api/audience/reveal
func (ab *AudienceBridge) HandleReveal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ab.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RevealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ab.jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GameID == "" || req.EventID == "" {
		ab.jsonError(w, "Missing game_id or event_id", http.StatusBadRequest)
		return
	}

	// TODO: Mark event as revealed in database
	// For now, log the action
	ab.logger.Event("AUDIENCE_REVEAL", req.AudienceID, "EventID:"+req.EventID)

	// Notify connected players that something was revealed
	ab.wsHub.BroadcastToGame(req.GameID, Message{
		Type:      MsgTypeEvent,
		Timestamp: time.Now().Unix(),
		Payload: map[string]interface{}{
			"event_type": "SECRET_REVEALED",
			"event_id":   req.EventID,
			"message":    "Un secreto ha sido expuesto...",
		},
	})

	ab.jsonSuccess(w, map[string]interface{}{
		"success":   true,
		"revealed":  true,
		"event_id":  req.EventID,
	})
}

// HandleGameStatus returns the current state for audience viewing.
// GET /api/audience/game/:gameID/status
func (ab *AudienceBridge) HandleGameStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ab.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		ab.jsonError(w, "Missing game_id", http.StatusBadRequest)
		return
	}

	// Return connected prisoners and basic stats
	connected := ab.wsHub.GetConnectedPrisoners(gameID)
	
	ab.jsonSuccess(w, map[string]interface{}{
		"game_id":            gameID,
		"connected_prisoners": connected,
		"online_count":       len(connected),
		"timestamp":          time.Now().Unix(),
	})
}

// RegisterRoutes sets up the audience API routes.
func (ab *AudienceBridge) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/audience/torture", ab.HandleTorture)
	mux.HandleFunc("/api/audience/reveal", ab.HandleReveal)
	mux.HandleFunc("/api/audience/status", ab.HandleGameStatus)
}

// jsonError sends an error response.
func (ab *AudienceBridge) jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// jsonSuccess sends a success response.
func (ab *AudienceBridge) jsonSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
