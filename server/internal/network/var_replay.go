// Package network - var_replay.go
// T015: VAR Replay endpoint - JSON export of betrayal history.
//
// This is the "VAR of Betrayal" viewer. It allows the audience and
// moderators to replay the immutable history of events.
package network

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// VARReplayHandler provides the VAR replay API.
type VARReplayHandler struct {
	eventLog *events.EventLog
	logger   *logger.Logger
}

// NewVARReplayHandler creates a new VAR replay handler.
func NewVARReplayHandler(el *events.EventLog, log *logger.Logger) *VARReplayHandler {
	return &VARReplayHandler{
		eventLog: el,
		logger:   log,
	}
}

// ReplayEvent is a sanitized event for public viewing.
type ReplayEvent struct {
	ID        string                 `json:"id"`
	Timestamp string                 `json:"timestamp"`
	GameDay   int                    `json:"game_day"`
	Type      string                 `json:"type"`
	ActorName string                 `json:"actor_name"` // Anonymized unless revealed
	TargetName string                `json:"target_name,omitempty"`
	Summary   string                 `json:"summary"`
	Impact    string                 `json:"impact"`
	IsRevealed bool                  `json:"is_revealed"`
	Details   map[string]interface{} `json:"details,omitempty"` // Only if revealed
}

// ReplayResponse is the API response for VAR replay.
type ReplayResponse struct {
	GameID       string        `json:"game_id"`
	TotalEvents  int           `json:"total_events"`
	FilteredBy   string        `json:"filtered_by,omitempty"`
	GeneratedAt  string        `json:"generated_at"`
	Events       []ReplayEvent `json:"events"`
}

// HandleReplay returns the VAR replay for a game.
// GET /api/var/replay?game_id=XXX&day=N&type=BETRAYAL&revealed_only=true
func (vh *VARReplayHandler) HandleReplay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		vh.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		vh.jsonError(w, "Missing game_id", http.StatusBadRequest)
		return
	}

	// Optional filters
	dayStr := r.URL.Query().Get("day")
	eventType := r.URL.Query().Get("type")
	revealedOnly := r.URL.Query().Get("revealed_only") == "true"

	// Get all events (from in-memory log for now)
	allEvents := vh.eventLog.Replay()

	// Filter and convert to replay format
	var replayEvents []ReplayEvent
	filterDesc := ""

	for _, e := range allEvents {
		// Day filter
		if dayStr != "" {
			day, _ := strconv.Atoi(dayStr)
			if e.GameDay != day {
				continue
			}
			filterDesc = "Day " + dayStr
		}

		// Type filter
		if eventType != "" && string(e.Type) != eventType {
			continue
		}

		// Revealed filter
		if revealedOnly && !e.IsRevealed {
			continue
		}

		replayEvents = append(replayEvents, vh.convertToReplayEvent(e))
	}

	response := ReplayResponse{
		GameID:      gameID,
		TotalEvents: len(replayEvents),
		FilteredBy:  filterDesc,
		GeneratedAt: time.Now().Format(time.RFC3339),
		Events:      replayEvents,
	}

	vh.logger.Event("VAR_REPLAY", "AUDIENCE", "GameID:"+gameID+" Events:"+strconv.Itoa(len(replayEvents)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleEventDetail returns details of a specific event.
// GET /api/var/event/:eventID
func (vh *VARReplayHandler) HandleEventDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		vh.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	eventID := r.URL.Query().Get("event_id")
	if eventID == "" {
		vh.jsonError(w, "Missing event_id", http.StatusBadRequest)
		return
	}

	// Find the event
	allEvents := vh.eventLog.Replay()
	for _, e := range allEvents {
		if e.ID == eventID {
			detail := vh.convertToReplayEvent(e)
			
			// Include full details if revealed
			if e.IsRevealed {
				if payload, ok := e.Payload.(map[string]interface{}); ok {
					detail.Details = payload
				}
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(detail)
			return
		}
	}

	vh.jsonError(w, "Event not found", http.StatusNotFound)
}

// HandleStats returns aggregate statistics for the VAR.
// GET /api/var/stats?game_id=XXX
func (vh *VARReplayHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		vh.jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allEvents := vh.eventLog.Replay()

	// Calculate stats
	stats := map[string]int{
		"total_events":      len(allEvents),
		"noise_events":      0,
		"sanity_changes":    0,
		"betrayals":         0,
		"revealed_secrets":  0,
	}

	for _, e := range allEvents {
		switch e.Type {
		case events.EventTypeNoiseEvent:
			stats["noise_events"]++
		case events.EventTypeSanityChange:
			stats["sanity_changes"]++
		case events.EventTypeBetrayal:
			stats["betrayals"]++
		}
		if e.IsRevealed {
			stats["revealed_secrets"]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"generated_at": time.Now().Format(time.RFC3339),
		"stats":        stats,
	})
}

// RegisterRoutes sets up the VAR API routes.
func (vh *VARReplayHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/var/replay", vh.HandleReplay)
	mux.HandleFunc("/api/var/event", vh.HandleEventDetail)
	mux.HandleFunc("/api/var/stats", vh.HandleStats)
}

// convertToReplayEvent transforms an internal event to public format.
func (vh *VARReplayHandler) convertToReplayEvent(e events.GameEvent) ReplayEvent {
	summary := vh.summarizeEvent(e)
	impact := vh.determineImpact(e)

	actorName := e.ActorID
	if e.ActorID == "SYSTEM_TWINS" {
		actorName = "Los Gemelos"
	}

	return ReplayEvent{
		ID:         e.ID,
		Timestamp:  e.Timestamp.Format("15:04:05"),
		GameDay:    e.GameDay,
		Type:       string(e.Type),
		ActorName:  actorName,
		TargetName: e.TargetID,
		Summary:    summary,
		Impact:     impact,
		IsRevealed: e.IsRevealed,
	}
}

// summarizeEvent creates a human-readable summary.
func (vh *VARReplayHandler) summarizeEvent(e events.GameEvent) string {
	switch e.Type {
	case events.EventTypeNoiseEvent:
		return "Los Gemelos desataron una tortura de ruido."
	case events.EventTypeSanityChange:
		return "La cordura de un prisionero fue afectada."
	case events.EventTypeBetrayal:
		return "Hubo una traición en la prisión."
	case events.EventTypeTimeTick:
		return "El tiempo avanzó en la prisión."
	default:
		return "Algo ocurrió..."
	}
}

// determineImpact classifies the event impact.
func (vh *VARReplayHandler) determineImpact(e events.GameEvent) string {
	switch e.Type {
	case events.EventTypeNoiseEvent, events.EventTypeSanityChange, events.EventTypeBetrayal:
		return "NEGATIVE"
	case events.EventTypeResourceIntake:
		return "POSITIVE"
	default:
		return "NEUTRAL"
	}
}

// jsonError sends an error response.
func (vh *VARReplayHandler) jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
