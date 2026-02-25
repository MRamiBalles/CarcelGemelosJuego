package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// SQLiteEventRepository implements EventRepository for SQLite.
type SQLiteEventRepository struct {
	db *sql.DB
}

func NewSQLiteEventRepository(db *sql.DB) *SQLiteEventRepository {
	return &SQLiteEventRepository{db: db}
}

func (r *SQLiteEventRepository) Append(ctx context.Context, event GameEvent) error {
	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO events (id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		event.ID, event.GameID, event.Timestamp, event.EventType, event.ActorID,
		event.TargetID, string(payloadBytes), event.GameDay, event.IsRevealed,
	)
	if err != nil {
		return fmt.Errorf("failed to append event: %w", err)
	}
	return nil
}

func (r *SQLiteEventRepository) getMany(ctx context.Context, query string, args ...interface{}) ([]GameEvent, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []GameEvent
	for rows.Next() {
		var e GameEvent
		var payloadStr string
		err := rows.Scan(
			&e.ID, &e.GameID, &e.Timestamp, &e.EventType, &e.ActorID,
			&e.TargetID, &payloadStr, &e.GameDay, &e.IsRevealed,
		)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(payloadStr), &e.Payload); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *SQLiteEventRepository) GetByGameID(ctx context.Context, gameID string) ([]GameEvent, error) {
	query := `SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed FROM events WHERE game_id = ? ORDER BY timestamp ASC`
	return r.getMany(ctx, query, gameID)
}

func (r *SQLiteEventRepository) GetByActorID(ctx context.Context, gameID, actorID string) ([]GameEvent, error) {
	query := `SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed FROM events WHERE game_id = ? AND actor_id = ? ORDER BY timestamp ASC`
	return r.getMany(ctx, query, gameID, actorID)
}

func (r *SQLiteEventRepository) GetByGameDay(ctx context.Context, gameID string, day int) ([]GameEvent, error) {
	query := `SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed FROM events WHERE game_id = ? AND game_day = ? ORDER BY timestamp ASC`
	return r.getMany(ctx, query, gameID, day)
}

func (r *SQLiteEventRepository) GetByEventType(ctx context.Context, gameID string, eventType string) ([]GameEvent, error) {
	query := `SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed FROM events WHERE game_id = ? AND event_type = ? ORDER BY timestamp ASC`
	return r.getMany(ctx, query, gameID, eventType)
}

func (r *SQLiteEventRepository) GetUnrevealed(ctx context.Context, gameID string) ([]GameEvent, error) {
	query := `SELECT id, game_id, timestamp, event_type, actor_id, target_id, payload, game_day, is_revealed FROM events WHERE game_id = ? AND is_revealed = 0 ORDER BY timestamp ASC`
	return r.getMany(ctx, query, gameID)
}

func (r *SQLiteEventRepository) MarkRevealed(ctx context.Context, eventID string) error {
	query := `UPDATE events SET is_revealed = 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, eventID)
	return err
}

// ---------------------------------------------------------
// SQLiteSnapshotRepository
// ---------------------------------------------------------

type SQLiteSnapshotRepository struct {
	db *sql.DB
}

func NewSQLiteSnapshotRepository(db *sql.DB) *SQLiteSnapshotRepository {
	return &SQLiteSnapshotRepository{db: db}
}

func (r *SQLiteSnapshotRepository) Upsert(ctx context.Context, snapshot PrisonerSnapshot) error {
	query := `
		INSERT INTO prisoners (prisoner_id, game_id, name, archetype, sanity, dignity, pot_contribution, is_isolated, is_sleeper, is_withdraw, cell_id, last_updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(prisoner_id) DO UPDATE SET
			name=excluded.name,
			archetype=excluded.archetype,
			sanity=excluded.sanity,
			dignity=excluded.dignity,
			pot_contribution=excluded.pot_contribution,
			is_isolated=excluded.is_isolated,
			is_sleeper=excluded.is_sleeper,
			is_withdraw=excluded.is_withdraw,
			cell_id=excluded.cell_id,
			last_updated=excluded.last_updated
	`
	_, err := r.db.ExecContext(ctx, query,
		snapshot.PrisonerID, snapshot.GameID, snapshot.Name, snapshot.Archetype,
		snapshot.Sanity, snapshot.Dignity, 0.0, snapshot.IsIsolated, snapshot.IsSleeper, snapshot.IsWithdraw, snapshot.CellID, time.Now(),
	)
	return err
}

func (r *SQLiteSnapshotRepository) GetByPrisonerID(ctx context.Context, prisonerID string) (*PrisonerSnapshot, error) {
	query := `SELECT prisoner_id, game_id, name, archetype, sanity, dignity, is_isolated, is_sleeper, is_withdraw, cell_id FROM prisoners WHERE prisoner_id = ?`
	var p PrisonerSnapshot
	err := r.db.QueryRowContext(ctx, query, prisonerID).Scan(
		&p.PrisonerID, &p.GameID, &p.Name, &p.Archetype, &p.Sanity, &p.Dignity, &p.IsIsolated, &p.IsSleeper, &p.IsWithdraw, &p.CellID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *SQLiteSnapshotRepository) GetByGameID(ctx context.Context, gameID string) ([]PrisonerSnapshot, error) {
	query := `SELECT prisoner_id, game_id, name, archetype, sanity, dignity, is_isolated, is_sleeper, is_withdraw, cell_id FROM prisoners WHERE game_id = ?`
	rows, err := r.db.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snaps []PrisonerSnapshot
	for rows.Next() {
		var p PrisonerSnapshot
		if err := rows.Scan(&p.PrisonerID, &p.GameID, &p.Name, &p.Archetype, &p.Sanity, &p.Dignity, &p.IsIsolated, &p.IsSleeper, &p.IsWithdraw, &p.CellID); err != nil {
			return nil, err
		}
		snaps = append(snaps, p)
	}
	return snaps, nil
}

func (r *SQLiteSnapshotRepository) RebuildFromEvents(ctx context.Context, gameID string, events []GameEvent) error {
	// Future optimization
	return nil
}
