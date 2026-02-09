// Package perception provides the "eyes" of Los Gemelos.
// T016: EventLog Summarizer - Reads events and builds context for LLM decisions.
//
// This package transforms raw events into a digestible narrative that
// the Cognition module can process. It implements the "Perception" layer
// of the Perception-Cognition-Action loop.
package perception

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// PrisonState represents the current emotional/social state of the prison.
type PrisonState struct {
	GameID           string            `json:"game_id"`
	CurrentDay       int               `json:"current_day"`
	CurrentHour      int               `json:"current_hour"`
	TotalPrisoners   int               `json:"total_prisoners"`
	OnlinePrisoners  int               `json:"online_prisoners"`
	AverageSanity    float64           `json:"average_sanity"`
	TensionLevel     string            `json:"tension_level"` // LOW, MEDIUM, HIGH, CRITICAL
	RecentBetrayals  int               `json:"recent_betrayals"`
	AudienceActivity int               `json:"audience_activity"` // Sadism points spent recently
	PrisonerSummaries map[string]string `json:"prisoner_summaries"`
	NarrativeSummary string            `json:"narrative_summary"` // LLM-ready context
}

// Perceiver reads the EventLog and builds context for the Cognition module.
type Perceiver struct {
	eventLog *events.EventLog
	logger   *logger.Logger
}

// NewPerceiver creates a new perception module.
func NewPerceiver(el *events.EventLog, log *logger.Logger) *Perceiver {
	return &Perceiver{
		eventLog: el,
		logger:   log,
	}
}

// BuildPrisonState analyzes recent events and builds a comprehensive state.
func (p *Perceiver) BuildPrisonState(ctx context.Context, gameID string, currentDay int) (*PrisonState, error) {
	allEvents := p.eventLog.Replay()
	
	state := &PrisonState{
		GameID:            gameID,
		CurrentDay:        currentDay,
		PrisonerSummaries: make(map[string]string),
	}

	// Analyze events from the last 3 game days
	recentEvents := p.filterRecentEvents(allEvents, currentDay-3)
	
	// Count betrayals and calculate metrics
	sanitySum := 0.0
	sanityCount := 0
	
	for _, e := range recentEvents {
		switch e.Type {
		case events.EventTypeBetrayal:
			state.RecentBetrayals++
		case events.EventTypeSanityChange:
			if payload, ok := e.Payload.(map[string]interface{}); ok {
				if newSan, ok := payload["new_sanity"].(float64); ok {
					sanitySum += newSan
					sanityCount++
				}
			}
		}
		
		// Track audience activity
		if strings.HasPrefix(e.ActorID, "AUDIENCE_") {
			state.AudienceActivity++
		}
	}

	// Calculate average sanity
	if sanityCount > 0 {
		state.AverageSanity = sanitySum / float64(sanityCount)
	} else {
		state.AverageSanity = 75.0 // Default assumption
	}

	// Determine tension level
	state.TensionLevel = p.calculateTensionLevel(state)

	// Build narrative summary for LLM context
	state.NarrativeSummary = p.buildNarrativeSummary(state, recentEvents)

	p.logger.Event("PERCEPTION", "TWINS", "State built: Tension="+state.TensionLevel)

	return state, nil
}

// filterRecentEvents returns events from the specified day onwards.
func (p *Perceiver) filterRecentEvents(allEvents []events.GameEvent, sinceDay int) []events.GameEvent {
	var recent []events.GameEvent
	for _, e := range allEvents {
		if e.GameDay >= sinceDay {
			recent = append(recent, e)
		}
	}
	return recent
}

// calculateTensionLevel determines the overall tension in the prison.
func (p *Perceiver) calculateTensionLevel(state *PrisonState) string {
	score := 0

	// Low sanity increases tension
	if state.AverageSanity < 30 {
		score += 3
	} else if state.AverageSanity < 50 {
		score += 2
	} else if state.AverageSanity < 70 {
		score += 1
	}

	// Betrayals increase tension
	score += state.RecentBetrayals * 2

	// Audience activity increases tension (they're paying to see chaos)
	if state.AudienceActivity > 10 {
		score += 2
	} else if state.AudienceActivity > 5 {
		score += 1
	}

	// Late game is inherently more tense
	if state.CurrentDay >= 18 {
		score += 3
	} else if state.CurrentDay >= 14 {
		score += 2
	} else if state.CurrentDay >= 7 {
		score += 1
	}

	switch {
	case score >= 8:
		return "CRITICAL"
	case score >= 5:
		return "HIGH"
	case score >= 2:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

// buildNarrativeSummary creates an LLM-ready context string.
func (p *Perceiver) buildNarrativeSummary(state *PrisonState, events []events.GameEvent) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("=== INFORME DE SITUACIÓN: DÍA %d ===\n", state.CurrentDay))
	sb.WriteString(fmt.Sprintf("Nivel de Tensión: %s\n", state.TensionLevel))
	sb.WriteString(fmt.Sprintf("Cordura Promedio: %.1f%%\n", state.AverageSanity))
	sb.WriteString(fmt.Sprintf("Traiciones Recientes: %d\n", state.RecentBetrayals))
	sb.WriteString(fmt.Sprintf("Actividad de la Audiencia: %d intervenciones\n\n", state.AudienceActivity))

	// Add recent notable events
	sb.WriteString("=== EVENTOS NOTABLES ===\n")
	notable := 0
	for _, e := range events {
		if notable >= 5 {
			break // Limit to 5 notable events to control token usage
		}
		
		switch e.Type {
		case events.EventTypeBetrayal:
			sb.WriteString(fmt.Sprintf("- [DÍA %d] TRAICIÓN: %s traicionó.\n", e.GameDay, e.ActorID))
			notable++
		case events.EventTypeNoiseEvent:
			sb.WriteString(fmt.Sprintf("- [DÍA %d] TORTURA: Los Gemelos activaron ruido.\n", e.GameDay))
			notable++
		}
	}

	if notable == 0 {
		sb.WriteString("- Sin eventos notables recientes.\n")
	}

	sb.WriteString("\n=== DECISIÓN REQUERIDA ===\n")
	sb.WriteString("¿Qué acción deben tomar Los Gemelos para maximizar el drama sin destruir el juego?\n")

	return sb.String()
}

// GetPrisonerProfile builds a summary of a specific prisoner's recent behavior.
func (p *Perceiver) GetPrisonerProfile(prisonerID string, events []events.GameEvent) string {
	var actions []string
	sanityHistory := []int{}

	for _, e := range events {
		if e.ActorID == prisonerID || e.TargetID == prisonerID {
			switch e.Type {
			case events.EventTypeBetrayal:
				actions = append(actions, "TRAICIÓN")
			case events.EventTypeSocialAction:
				actions = append(actions, "INTERACCIÓN_SOCIAL")
			case events.EventTypeSanityChange:
				if payload, ok := e.Payload.(map[string]interface{}); ok {
					if newSan, ok := payload["new_sanity"].(float64); ok {
						sanityHistory = append(sanityHistory, int(newSan))
					}
				}
			}
		}
	}

	profile := fmt.Sprintf("Prisionero %s: %d acciones registradas. ", prisonerID, len(actions))
	if len(sanityHistory) > 0 {
		lastSanity := sanityHistory[len(sanityHistory)-1]
		profile += fmt.Sprintf("Cordura actual: %d%%. ", lastSanity)
	}

	return profile
}
