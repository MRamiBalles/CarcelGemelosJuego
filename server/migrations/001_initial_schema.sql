-- T010: PostgreSQL Schema for "Cárcel de los Gemelos"
-- This is an APPEND-ONLY LEDGER. NO UPDATE/DELETE allowed.
-- The history of the prison is sacred.

-- ==================================================
-- GAMES (Partidas)
-- ==================================================
CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, ACTIVE, COMPLETED
    current_day INT NOT NULL DEFAULT 1,
    current_hour INT NOT NULL DEFAULT 6,
    winner_id UUID -- NULL until endgame
);

-- ==================================================
-- PRISONERS (Prisioneros)
-- ==================================================
CREATE TABLE IF NOT EXISTS prisoners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id),
    user_id UUID NOT NULL, -- External auth reference
    name VARCHAR(100) NOT NULL,
    archetype VARCHAR(20) NOT NULL, -- VETERAN, MYSTIC, SHOWMAN, REDEEMED
    cellmate_id UUID, -- Partner reference
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_eliminated BOOLEAN NOT NULL DEFAULT FALSE,
    eliminated_day INT
);

-- ==================================================
-- EVENT LOG (El VAR de la Traición)
-- IMMUTABLE: INSERT ONLY. History is truth.
-- ==================================================
CREATE TABLE IF NOT EXISTS event_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    event_type VARCHAR(50) NOT NULL,
    actor_id VARCHAR(100) NOT NULL, -- Prisoner ID or "SYSTEM_TWINS"
    target_id VARCHAR(100), -- Affected prisoner (optional)
    payload JSONB NOT NULL, -- Event-specific data (immutable)
    game_day INT NOT NULL,
    is_revealed BOOLEAN NOT NULL DEFAULT FALSE, -- Exposed to audience?
    
    -- Index for fast replay by game
    CONSTRAINT fk_game FOREIGN KEY (game_id) REFERENCES games(id)
);

-- Indexes for common queries
CREATE INDEX idx_event_log_game_id ON event_log(game_id);
CREATE INDEX idx_event_log_actor_id ON event_log(actor_id);
CREATE INDEX idx_event_log_event_type ON event_log(event_type);
CREATE INDEX idx_event_log_game_day ON event_log(game_day);
CREATE INDEX idx_event_log_timestamp ON event_log(timestamp);

-- ==================================================
-- SECURITY: Prevent UPDATE/DELETE on event_log
-- ==================================================
CREATE OR REPLACE FUNCTION prevent_event_modification()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Event log is immutable. Updates and deletes are prohibited.';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_event_log_immutable
BEFORE UPDATE OR DELETE ON event_log
FOR EACH ROW
EXECUTE FUNCTION prevent_event_modification();

-- ==================================================
-- PRISONER SNAPSHOTS (Estado actual para lecturas rápidas)
-- This is a MATERIALIZED cache, rebuilt from events.
-- ==================================================
CREATE TABLE IF NOT EXISTS prisoner_snapshots (
    prisoner_id UUID PRIMARY KEY REFERENCES prisoners(id),
    game_id UUID NOT NULL REFERENCES games(id),
    hunger INT NOT NULL DEFAULT 100,
    thirst INT NOT NULL DEFAULT 100,
    sanity INT NOT NULL DEFAULT 100,
    dignity INT NOT NULL DEFAULT 100,
    loyalty INT NOT NULL DEFAULT 50,
    empathy INT NOT NULL DEFAULT 50,
    is_sleeper BOOLEAN NOT NULL DEFAULT FALSE,
    is_withdraw BOOLEAN NOT NULL DEFAULT FALSE,
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index for fast lookups
CREATE INDEX idx_prisoner_snapshots_game_id ON prisoner_snapshots(game_id);

COMMENT ON TABLE event_log IS 'The VAR of Betrayal. Immutable append-only ledger of all game events.';
COMMENT ON TABLE prisoner_snapshots IS 'Fast-read cache of current prisoner state. Rebuilt from event_log.';
