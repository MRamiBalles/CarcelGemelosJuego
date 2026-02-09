-- T031: Database indexes for VAR Replay performance
-- These indexes optimize the most common query patterns

-- Index for event replay by game (most common query)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_log_game_timestamp 
ON event_log (game_id, timestamp DESC);

-- Index for actor-specific queries (prisoner history)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_log_actor 
ON event_log (actor_id, timestamp DESC) 
WHERE actor_id IS NOT NULL;

-- Index for event type filtering (VAR filtering)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_log_type 
ON event_log (event_type, timestamp DESC);

-- Composite index for the most common VAR query
-- "Get all events for game X, day Y, ordered by time"
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_log_game_day 
ON event_log (game_id, (payload->>'game_day')::int, timestamp DESC);

-- JSONB index for payload searching (GIN)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_log_payload 
ON event_log USING GIN (payload jsonb_path_ops);

-- Partial index for unrevealed secrets (hidden events)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_event_log_hidden 
ON event_log (game_id, timestamp DESC) 
WHERE (payload->>'is_revealed')::boolean = false;

-- ============================================
-- Analyze tables after index creation
-- ============================================
ANALYZE event_log;
ANALYZE games;
ANALYZE prisoners;

-- ============================================
-- Query optimization hints for the application
-- ============================================
COMMENT ON INDEX idx_event_log_game_timestamp IS 
    'Primary index for VAR replay - use for GetByGameID queries';

COMMENT ON INDEX idx_event_log_actor IS 
    'Use for prisoner history reconstruction';

COMMENT ON INDEX idx_event_log_hidden IS 
    'Use for "Los Gemelos can see all" queries - reveals secrets';
