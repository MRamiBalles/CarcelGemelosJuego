// Package cache provides Redis-based caching for quick state reads.
// T012: Redis cache for prisoner snapshots (not the source of truth).
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// RedisClient is an interface for Redis operations.
// This allows for easy mocking in tests.
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
}

// PrisonerCache provides fast access to prisoner state snapshots.
type PrisonerCache struct {
	client     RedisClient
	expiration time.Duration
}

// NewPrisonerCache creates a new prisoner cache instance.
func NewPrisonerCache(client RedisClient) *PrisonerCache {
	return &PrisonerCache{
		client:     client,
		expiration: 15 * time.Minute, // Cache expires after 15 minutes
	}
}

// PrisonerState represents the cached state of a prisoner.
type PrisonerState struct {
	PrisonerID string `json:"prisoner_id"`
	Name       string `json:"name"`
	Archetype  string `json:"archetype"`
	Hunger     int    `json:"hunger"`
	Thirst     int    `json:"thirst"`
	Sanity     int    `json:"sanity"`
	Dignity    int    `json:"dignity"`
	Loyalty    int    `json:"loyalty"`
	IsSleeper  bool   `json:"is_sleeper"`
	LastSync   int64  `json:"last_sync"` // Unix timestamp
}

// SetPrisonerState caches the current state of a prisoner.
func (c *PrisonerCache) SetPrisonerState(ctx context.Context, gameID string, state PrisonerState) error {
	key := c.prisonerKey(gameID, state.PrisonerID)
	
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal prisoner state: %w", err)
	}

	return c.client.Set(ctx, key, data, c.expiration)
}

// GetPrisonerState retrieves the cached state of a prisoner.
func (c *PrisonerCache) GetPrisonerState(ctx context.Context, gameID, prisonerID string) (*PrisonerState, error) {
	key := c.prisonerKey(gameID, prisonerID)
	
	data, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, err // Cache miss or error
	}

	var state PrisonerState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prisoner state: %w", err)
	}

	return &state, nil
}

// SetGameState caches the current state of all prisoners in a game.
// Uses Redis Hash for efficient storage.
func (c *PrisonerCache) SetGameState(ctx context.Context, gameID string, states map[string]PrisonerState) error {
	key := c.gameKey(gameID)
	
	values := make([]interface{}, 0, len(states)*2)
	for id, state := range states {
		data, err := json.Marshal(state)
		if err != nil {
			return fmt.Errorf("failed to marshal state for %s: %w", id, err)
		}
		values = append(values, id, string(data))
	}

	return c.client.HSet(ctx, key, values...)
}

// GetGameState retrieves the cached state of all prisoners in a game.
func (c *PrisonerCache) GetGameState(ctx context.Context, gameID string) (map[string]PrisonerState, error) {
	key := c.gameKey(gameID)
	
	data, err := c.client.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}

	states := make(map[string]PrisonerState)
	for id, jsonStr := range data {
		var state PrisonerState
		if err := json.Unmarshal([]byte(jsonStr), &state); err != nil {
			return nil, fmt.Errorf("failed to unmarshal state for %s: %w", id, err)
		}
		states[id] = state
	}

	return states, nil
}

// InvalidateGame removes all cached state for a game.
func (c *PrisonerCache) InvalidateGame(ctx context.Context, gameID string) error {
	key := c.gameKey(gameID)
	return c.client.Del(ctx, key)
}

// prisonerKey generates the Redis key for a specific prisoner.
func (c *PrisonerCache) prisonerKey(gameID, prisonerID string) string {
	return fmt.Sprintf("game:%s:prisoner:%s", gameID, prisonerID)
}

// gameKey generates the Redis key for all prisoners in a game.
func (c *PrisonerCache) gameKey(gameID string) string {
	return fmt.Sprintf("game:%s:prisoners", gameID)
}
