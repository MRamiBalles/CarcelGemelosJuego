// Package storage - postgres.go
// T011: PostgreSQL implementation of EventRepository.
package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// PostgresEventRepository implements EventRepository using PostgreSQL.
type PostgresEventRepository struct {
	db *sql.DB
}

// NewPostgresEventRepository creates a new PostgreSQL event repository.
func NewPostgresEventRepository(db *sql.DB) *PostgresEventRepository {
	return &PostgresEventRepository{db: db}
}

// Append inserts a new event into the immutable ledger.
func (r *PostgresEventRepository) Append(ctx context.Context, event GameEvent) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO event_log (id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.ExecContext(ctx, query,
		event.ID,
		event.GameID,
		event.Timestamp,
		event.EventType,
		event.ActorID,
		event.TargetID,
		payloadJSON,
		event.GameDay,
		event.IsRevealed,
	)

	if err != nil {
		return fmt.Errorf("failed to append event: %w", err)
	}

	return nil
}

// GetByGameID retrieves all events for a game (the full VAR replay).
func (r *PostgresEventRepository) GetByGameID(ctx context.Context, gameID string) ([]GameEvent, error) {
	query := `
		SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed
		FROM event_log
		WHERE game_id = $1
		ORDER BY timestamp ASC
	`

	return r.queryEvents(ctx, query, gameID)
}

// GetByActorID retrieves all events performed by an actor.
func (r *PostgresEventRepository) GetByActorID(ctx context.Context, gameID, actorID string) ([]GameEvent, error) {
	query := `
		SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed
		FROM event_log
		WHERE game_id = $1 AND actor_id = $2
		ORDER BY timestamp ASC
	`

	return r.queryEvents(ctx, query, gameID, actorID)
}

// GetByGameDay retrieves all events from a specific in-game day.
func (r *PostgresEventRepository) GetByGameDay(ctx context.Context, gameID string, day int) ([]GameEvent, error) {
	query := `
		SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed
		FROM event_log
		WHERE game_id = $1 AND game_day = $2
		ORDER BY timestamp ASC
	`

	return r.queryEvents(ctx, query, gameID, day)
}

// GetByEventType retrieves all events of a specific type.
func (r *PostgresEventRepository) GetByEventType(ctx context.Context, gameID string, eventType string) ([]GameEvent, error) {
	query := `
		SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed
		FROM event_log
		WHERE game_id = $1 AND event_type = $2
		ORDER BY timestamp ASC
	`

	return r.queryEvents(ctx, query, gameID, eventType)
}

// GetUnrevealed retrieves events not yet shown to the audience.
func (r *PostgresEventRepository) GetUnrevealed(ctx context.Context, gameID string) ([]GameEvent, error) {
	query := `
		SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed
		FROM event_log
		WHERE game_id = $1 AND is_revealed = FALSE
		ORDER BY timestamp ASC
	`

	return r.queryEvents(ctx, query, gameID)
}

// MarkRevealed marks an event as revealed to the audience.
// NOTE: This is a controlled exception to immutability - only for reveal status.
func (r *PostgresEventRepository) MarkRevealed(ctx context.Context, eventID string) error {
	query := `UPDATE event_log SET is_revealed = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, eventID)
	return err
}

// queryEvents is a helper to execute queries and scan results.
func (r *PostgresEventRepository) queryEvents(ctx context.Context, query string, args ...interface{}) ([]GameEvent, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []GameEvent
	for rows.Next() {
		var e GameEvent
		var payloadJSON []byte
		var targetID sql.NullString

		err := rows.Scan(
			&e.ID,
			&e.GameID,
			&e.Timestamp,
			&e.EventType,
			&e.ActorID,
			&targetID,
			&payloadJSON,
			&e.GameDay,
			&e.IsRevealed,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if targetID.Valid {
			e.TargetID = targetID.String
		}

		if err := json.Unmarshal(payloadJSON, &e.Payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		events = append(events, e)
	}

	return events, nil
}

// Ensure PostgresEventRepository implements EventRepository
var _ EventRepository = (*PostgresEventRepository)(nil)
