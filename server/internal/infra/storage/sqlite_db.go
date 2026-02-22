package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// InitSQLite initializes the local SQLite database and creates the necessary schemas
// for persisting the game state, prisoners, and the immutable event log.
func InitSQLite(dbPath string) (*sql.DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	// Create tables
	if err := createSchemas(db); err != nil {
		return nil, fmt.Errorf("failed to create schemas: %w", err)
	}

	return db, nil
}

func createSchemas(db *sql.DB) error {
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS game_state (
			game_id TEXT PRIMARY KEY,
			current_day INTEGER NOT NULL DEFAULT 1,
			tension_level TEXT NOT NULL DEFAULT 'LOW',
			last_updated DATETIME NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS prisoners (
			prisoner_id TEXT PRIMARY KEY,
			game_id TEXT NOT NULL,
			name TEXT,
			archetype TEXT,
			sanity INTEGER NOT NULL,
			dignity INTEGER NOT NULL,
			pot_contribution REAL NOT NULL DEFAULT 0.0,
			is_isolated BOOLEAN NOT NULL DEFAULT 0,
			traits_json TEXT,
			is_sleeper BOOLEAN NOT NULL DEFAULT 0,
			is_withdraw BOOLEAN NOT NULL DEFAULT 0,
			last_updated DATETIME NOT NULL,
			FOREIGN KEY (game_id) REFERENCES game_state(game_id)
		);`,
		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			game_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			event_type TEXT NOT NULL,
			actor_id TEXT NOT NULL,
			target_id TEXT NOT NULL,
			payload TEXT NOT NULL,
			game_day INTEGER NOT NULL,
			is_revealed BOOLEAN NOT NULL DEFAULT 0,
			FOREIGN KEY (game_id) REFERENCES game_state(game_id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_events_game_id ON events(game_id);`,
		`CREATE INDEX IF NOT EXISTS idx_events_actor_id ON events(actor_id);`,
		`CREATE INDEX IF NOT EXISTS idx_events_game_day ON events(game_day);`,
	}

	for _, query := range schemas {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
