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
	EventTypeResourceIntake       EventType = "RESOURCE_INTAKE"
	EventTypeNoiseEvent           EventType = "NOISE_EVENT"
	EventTypeSocialAction         EventType = "SOCIAL_ACTION"
	EventTypeVote                 EventType = "VOTE"
	EventTypeBetrayal             EventType = "BETRAYAL"
	EventTypePrivacyBreach        EventType = "PRIVACY_BREACH"
	EventTypeTimeTick             EventType = "TIME_TICK"
	EventTypeSanityChange         EventType = "SANITY_CHANGE"
	EventTypeLoyaltyChange        EventType = "LOYALTY_CHANGE"
	EventTypeToiletUse            EventType = "TOILET_USE"
	EventTypeDoorLock             EventType = "DOOR_LOCK"
	EventTypeDoorOpen             EventType = "DOOR_OPEN"
	EventTypeAudioTorture         EventType = "AUDIO_TORTURE"
	EventTypeAggressiveEmote      EventType = "AGGRESSIVE_EMOTE"
	EventTypeLockdownBang         EventType = "LOCKDOWN_BANG"
	EventTypeInsult               EventType = "INSULT"
	EventTypeSteal                EventType = "STEAL"
	EventTypeFinalDilemmaStart    EventType = "FINAL_DILEMMA_START"
	EventTypeFinalDilemmaDecision EventType = "FINAL_DILEMMA_DECISION"
	EventTypeOraclePainfulTruth   EventType = "ORACLE_PAINFUL_TRUTH"
)

// AudioTorturePayload holds the details for unavoidable sound events
type AudioTorturePayload struct {
	SoundName string `json:"soundName"`
	Duration  int    `json:"duration"` // in game minutes
}

// GameEvent represents an immutable record of an action in the game.
type GameEvent struct {
	ID         string      `json:"id"`
	Timestamp  time.Time   `json:"timestamp"`
	Type       EventType   `json:"type"`
	ActorID    string      `json:"actor_id"`  // Who performed the action
	TargetID   string      `json:"target_id"` // Who was affected (optional)
	Payload    interface{} `json:"payload"`   // Event-specific data
	GameDay    int         `json:"game_day"`
	IsRevealed bool        `json:"is_revealed"` // Exposed to audience?
}

// EventPersister defines how an event is durably stored.
type EventPersister interface {
	Append(event GameEvent) error
}

// EventLog is the in-memory append-only log of game events.
// In production, this would be backed by PostgreSQL/Redis or SQLite.
type EventLog struct {
	mu        sync.RWMutex
	events    []GameEvent
	persister EventPersister
}

// NewEventLog creates a new event log with an optional persister.
func NewEventLog(persister EventPersister) *EventLog {
	return &EventLog{
		events:    make([]GameEvent, 0),
		persister: persister,
	}
}

// Append adds a new event to the log. Events are immutable once appended.
func (el *EventLog) Append(event GameEvent) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.events = append(el.events, event)

	if el.persister != nil {
		// Write through to persistent storage
		// In a real high-throughput system this might be buffered/async
		go func(e GameEvent) {
			_ = el.persister.Append(e)
		}(event)
	}
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

// GenerateEventID creates a unique event identifier.
func GenerateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomSuffix()
}

// randomSuffix generates a short random string.
func randomSuffix() string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
