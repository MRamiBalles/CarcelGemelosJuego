-- T011: Hardcore Mechanics Schema Update

-- 1. Add Pot Contribution tracking (Single Winner Rule)
ALTER TABLE prisoners ADD COLUMN IF NOT EXISTS pot_contribution DECIMAL(12,2) NOT NULL DEFAULT 0.00;

-- 2. Add DNA/Traits/States columns to Snapshots for fast lookup
ALTER TABLE prisoner_snapshots ADD COLUMN IF NOT EXISTS pot_contribution DECIMAL(12,2) NOT NULL DEFAULT 0.00;
ALTER TABLE prisoner_snapshots ADD COLUMN IF NOT EXISTS states JSONB DEFAULT '{}'; -- Active states
ALTER TABLE prisoner_snapshots ADD COLUMN IF NOT EXISTS traits JSONB DEFAULT '[]'; -- Immutable traits

-- 3. The Dilemma Table (Day 21 Decisions)
CREATE TABLE IF NOT EXISTS dilemma_decisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id),
    prisoner_id UUID NOT NULL REFERENCES prisoners(id),
    partner_id UUID NOT NULL REFERENCES prisoners(id),
    decision VARCHAR(20) NOT NULL, -- "BETRAY", "COLLABORATE"
    decision_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    outcome VARCHAR(50), -- "WON_ALL", "LOST_ALL", "SHARED"
    
    CONSTRAINT fk_dilemma_game FOREIGN KEY (game_id) REFERENCES games(id)
);

-- Index for fast retrieval of decisions
CREATE INDEX idx_dilemma_game_id ON dilemma_decisions(game_id);
