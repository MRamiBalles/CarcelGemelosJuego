// Package events provides the Event Sourcing system for the game.
// This is the "VAR of Betrayal" - an immutable log of all critical actions.
package events

import (
	"sync"
	"time"
)

// EventType defines the category of a game event.
type EventType string

const (
	EventTypeResourceIntake EventType = "RESOURCE_INTAKE"
	EventTypeNoiseEvent     EventType = "NOISE_EVENT"
	EventTypeSocialAction   EventType = "SOCIAL_ACTION"
	EventTypeVote           EventType = "VOTE"
	EventTypeBetrayal       EventType = "BETRAYAL"
	EventTypePrivacyBreach  EventType = "PRIVACY_BREACH"
)

// GameEvent represents an immutable record of an action in the game.
type GameEvent struct {
	ID         string      `json:"id"`
	Timestamp  time.Time   `json:"timestamp"`
	Type       EventType   `json:"type"`
	ActorID    string      `json:"actor_id"`    // Who performed the action
	TargetID   string      `json:"target_id"`   // Who was affected (optional)
	Payload    interface{} `json:"payload"`     // Event-specific data
	GameDay    int         `json:"game_day"`
	IsRevealed bool        `json:"is_revealed"` // Exposed to audience?
}

// EventLog is the in-memory append-only log of game events.
// In production, this would be backed by PostgreSQL/Redis.
type EventLog struct {
	mu     sync.RWMutex
	events []GameEvent
}

// NewEventLog creates a new empty event log.
func NewEventLog() *EventLog {
	return &EventLog{
		events: make([]GameEvent, 0),
	}
}

// Append adds a new event to the log. Events are immutable once appended.
func (el *EventLog) Append(event GameEvent) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.events = append(el.events, event)
}

// GetByActor returns all events performed by a specific actor.
func (el *EventLog) GetByActor(actorID string) []GameEvent {
	el.mu.RLock()
	defer el.mu.RUnlock()

	var result []GameEvent
	for _, e := range el.events {
		if e.ActorID == actorID {
			result = append(result, e)
		}
	}
	return result
}

// GetByDay returns all events that occurred on a specific game day.
func (el *EventLog) GetByDay(day int) []GameEvent {
	el.mu.RLock()
	defer el.mu.RUnlock()

	var result []GameEvent
	for _, e := range el.events {
		if e.GameDay == day {
			result = append(result, e)
		}
	}
	return result
}

// Replay returns the full history of events for state reconstruction.
// This is the "Reality Recap" system.
func (el *EventLog) Replay() []GameEvent {
	el.mu.RLock()
	defer el.mu.RUnlock()
	return el.events
}
