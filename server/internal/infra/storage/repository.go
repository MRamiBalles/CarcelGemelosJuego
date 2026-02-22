// Package storage provides the persistence layer for the game server.
// This package implements the repository pattern to keep the domain pure.
package storage

import (
	"context"
	"time"
)

// GameEvent mirrors the domain event structure for persistence.
// The domain package should NOT import this; use interfaces instead.
type GameEvent struct {
	ID         string                 `json:"id" db:"id"`
	GameID     string                 `json:"game_id" db:"game_id"`
	Timestamp  time.Time              `json:"timestamp" db:"timestamp"`
	EventType  string                 `json:"event_type" db:"event_type"`
	ActorID    string                 `json:"actor_id" db:"actor_id"`
	TargetID   string                 `json:"target_id" db:"target_id"`
	Payload    map[string]interface{} `json:"payload" db:"payload"`
	GameDay    int                    `json:"game_day" db:"game_day"`
	IsRevealed bool                   `json:"is_revealed" db:"is_revealed"`
}

// EventRepository defines the interface for event persistence.
// The domain uses this interface; the implementation is in infra.
type EventRepository interface {
	// Append adds a new event to the immutable ledger.
	Append(ctx context.Context, event GameEvent) error

	// GetByGameID retrieves all events for a specific game (for replay).
	GetByGameID(ctx context.Context, gameID string) ([]GameEvent, error)

	// GetByActorID retrieves all events performed by an actor.
	GetByActorID(ctx context.Context, gameID, actorID string) ([]GameEvent, error)

	// GetByGameDay retrieves all events from a specific in-game day.
	GetByGameDay(ctx context.Context, gameID string, day int) ([]GameEvent, error)

	// GetByEventType retrieves all events of a specific type.
	GetByEventType(ctx context.Context, gameID string, eventType string) ([]GameEvent, error)

	// GetUnrevealed retrieves events not yet shown to the audience.
	GetUnrevealed(ctx context.Context, gameID string) ([]GameEvent, error)

	// MarkRevealed marks an event as revealed to the audience.
	MarkRevealed(ctx context.Context, eventID string) error
}

// PrisonerSnapshot represents the current state of a prisoner for quick reads.
type PrisonerSnapshot struct {
	PrisonerID  string    `json:"prisoner_id" db:"prisoner_id"`
	GameID      string    `json:"game_id" db:"game_id"`
	Name        string    `json:"name" db:"name"`
	Archetype   string    `json:"archetype" db:"archetype"`
	Hunger      int       `json:"hunger" db:"hunger"`
	Thirst      int       `json:"thirst" db:"thirst"`
	Sanity      int       `json:"sanity" db:"sanity"`
	Dignity     int       `json:"dignity" db:"dignity"`
	Loyalty     int       `json:"loyalty" db:"loyalty"`
	Empathy     int       `json:"empathy" db:"empathy"`
	IsIsolated  bool      `json:"is_isolated" db:"is_isolated"`
	IsSleeper   bool      `json:"is_sleeper" db:"is_sleeper"`
	IsWithdraw  bool      `json:"is_withdraw" db:"is_withdraw"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
}

// SnapshotRepository defines the interface for prisoner state snapshots.
type SnapshotRepository interface {
	// Upsert updates or inserts a prisoner snapshot.
	Upsert(ctx context.Context, snapshot PrisonerSnapshot) error

	// GetByPrisonerID retrieves a specific prisoner's snapshot.
	GetByPrisonerID(ctx context.Context, prisonerID string) (*PrisonerSnapshot, error)

	// GetByGameID retrieves all snapshots for a game.
	GetByGameID(ctx context.Context, gameID string) ([]PrisonerSnapshot, error)

	// RebuildFromEvents reconstructs all snapshots from the event log.
	RebuildFromEvents(ctx context.Context, gameID string, events []GameEvent) error
}
