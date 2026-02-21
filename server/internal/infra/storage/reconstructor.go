// Package storage - reconstructor.go
// T012: Reality Recap - Rebuilds prisoner state from the event log.
// This is the core of Event Sourcing: state = f(events).
package storage

import (
	"context"
	"fmt"
)

// Reconstructor rebuilds prisoner state from the event log.
// This is used for:
// 1. Login "Reality Recap" - show what happened while offline
// 2. Snapshot rebuilding after cache invalidation
// 3. Auditing and debugging
type Reconstructor struct {
	eventRepo EventRepository
}

// NewReconstructor creates a new state reconstructor.
func NewReconstructor(eventRepo EventRepository) *Reconstructor {
	return &Reconstructor{eventRepo: eventRepo}
}

// RebuiltState holds the reconstructed state for a prisoner.
type RebuiltState struct {
	PrisonerID string
	Hunger     int
	Thirst     int
	Sanity     int
	Dignity    int
	Loyalty    int
	IsWithdraw bool
}

// RecapEvent is a simplified event for the "Reality Recap" screen.
type RecapEvent struct {
	Timestamp  string `json:"timestamp"`
	EventType  string `json:"event_type"`
	Summary    string `json:"summary"` // Human-readable description
	Impact     string `json:"impact"`  // "POSITIVE", "NEGATIVE", "NEUTRAL"
	IsRevealed bool   `json:"is_revealed"`
}

// RebuildPrisonerState reconstructs a prisoner's current state from events.
func (r *Reconstructor) RebuildPrisonerState(ctx context.Context, gameID, prisonerID string, initialState RebuiltState) (*RebuiltState, error) {
	events, err := r.eventRepo.GetByActorID(ctx, gameID, prisonerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events for prisoner: %w", err)
	}

	// Also get events where this prisoner was the TARGET
	targetEvents, err := r.eventRepo.GetByGameID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game events: %w", err)
	}

	state := initialState

	// Process events in chronological order
	for _, e := range events {
		r.applyEventToState(&state, e)
	}

	// Process events where prisoner was target
	for _, e := range targetEvents {
		if e.TargetID == prisonerID {
			r.applyEventToState(&state, e)
		}
	}

	return &state, nil
}

// GenerateRecap creates the "Reality Recap" for a prisoner since a given day.
func (r *Reconstructor) GenerateRecap(ctx context.Context, gameID, prisonerID string, sinceDay int) ([]RecapEvent, error) {
	allEvents, err := r.eventRepo.GetByGameID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	var recap []RecapEvent

	for _, e := range allEvents {
		// Filter: only events from the requested day onwards
		if e.GameDay < sinceDay {
			continue
		}

		// Filter: only events relevant to this prisoner
		if e.ActorID != prisonerID && e.TargetID != prisonerID && e.ActorID != "SYSTEM_TWINS" {
			continue
		}

		recap = append(recap, RecapEvent{
			Timestamp:  e.Timestamp.Format("15:04 Day 2"),
			EventType:  e.EventType,
			Summary:    r.summarizeEvent(e, prisonerID),
			Impact:     r.determineImpact(e),
			IsRevealed: e.IsRevealed,
		})
	}

	return recap, nil
}

// applyEventToState modifies state based on event type.
func (r *Reconstructor) applyEventToState(state *RebuiltState, event GameEvent) {
	switch event.EventType {
	case "SANITY_CHANGE":
		if delta, ok := event.Payload["delta"]; ok {
			if d, ok := delta.(float64); ok {
				state.Sanity += int(d)
			}
		}
	case "LOYALTY_CHANGE":
		if delta, ok := event.Payload["delta"]; ok {
			if d, ok := delta.(float64); ok {
				state.Loyalty += int(d)
			}
		}
	case "RESOURCE_INTAKE":
		if hunger, ok := event.Payload["hunger_delta"]; ok {
			if h, ok := hunger.(float64); ok {
				state.Hunger += int(h)
			}
		}
		if thirst, ok := event.Payload["thirst_delta"]; ok {
			if t, ok := thirst.(float64); ok {
				state.Thirst += int(t)
			}
		}
	}

	// Clamp values
	if state.Sanity < 0 {
		state.Sanity = 0
	}
	if state.Sanity > 100 {
		state.Sanity = 100
	}
}

// summarizeEvent creates a human-readable summary.
func (r *Reconstructor) summarizeEvent(event GameEvent, observerID string) string {
	switch event.EventType {
	case "NOISE_EVENT":
		return "Los Gemelos desataron una tortura de ruido."
	case "SANITY_CHANGE":
		if event.TargetID == observerID {
			return "Tu cordura fue afectada."
		}
		return "La cordura de otro prisionero fue afectada."
	case "BETRAYAL":
		return "Hubo una traición..."
	default:
		return "Algo ocurrió en la prisión."
	}
}

// determineImpact classifies the event impact.
func (r *Reconstructor) determineImpact(event GameEvent) string {
	switch event.EventType {
	case "NOISE_EVENT", "SANITY_CHANGE", "BETRAYAL":
		return "NEGATIVE"
	case "RESOURCE_INTAKE":
		return "POSITIVE"
	default:
		return "NEUTRAL"
	}
}
