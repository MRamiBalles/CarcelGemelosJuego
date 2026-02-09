// Package engine - sanity_system.go
// T009: Sanity System - Reacts to NoiseEvents and applies Sanity changes.
//
// This is a SUBSCRIBER to events from the EventLog.
// It processes NoiseEvents and emits SanityChangeEvents.
package engine

import (
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/rules"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// SanityChangePayload records the sanity modification for audit.
type SanityChangePayload struct {
	PrisonerID   string `json:"prisoner_id"`
	PreviousSan  int    `json:"previous_sanity"`
	NewSanity    int    `json:"new_sanity"`
	Delta        int    `json:"delta"`
	Cause        string `json:"cause"`        // "NOISE", "PRIVACY", "SOCIAL"
	CauseEventID string `json:"cause_event_id"` // Links to the triggering event
}

// SanitySystem processes events and applies sanity modifications.
type SanitySystem struct {
	eventLog  *events.EventLog
	logger    *logger.Logger
	prisoners map[string]*prisoner.Prisoner // In-memory state (would be DB in prod)
}

// NewSanitySystem creates a new sanity processing system.
func NewSanitySystem(eventLog *events.EventLog, log *logger.Logger) *SanitySystem {
	return &SanitySystem{
		eventLog:  eventLog,
		logger:    log,
		prisoners: make(map[string]*prisoner.Prisoner),
	}
}

// RegisterPrisoner adds a prisoner to the system's tracked state.
func (ss *SanitySystem) RegisterPrisoner(p *prisoner.Prisoner) {
	ss.prisoners[p.ID] = p
}

// OnNoiseEvent is the subscriber hook for NoiseEvents.
// It calculates sanity drain and emits SanityChangeEvents.
func (ss *SanitySystem) OnNoiseEvent(noiseEvent events.GameEvent) {
	payload, ok := noiseEvent.Payload.(NoiseEventPayload)
	if !ok {
		ss.logger.Error("Failed to parse NoiseEventPayload")
		return
	}

	for _, p := range ss.prisoners {
		// Skip if not in target zone (simplified: ALL affects everyone)
		if payload.TargetZone != "ALL" && payload.TargetZone != p.ID {
			continue
		}

		// Calculate drain using pure domain rules
		params := rules.SanityDrainParams{
			IsMystic:     p.Archetype == prisoner.ArchetypeMystic,
			NoiseLevel:   payload.Intensity,
			IsPrivacyHit: false,
		}
		drain := rules.CalculateSanityDrain(p, params)

		// Apply the drain
		previousSanity := p.Sanity
		p.Sanity -= drain
		if p.Sanity < 0 {
			p.Sanity = 0
			ss.logger.Warn("BREAKDOWN: " + p.Name + " has reached 0 Sanity!")
		}

		// Emit SanityChangeEvent for audit trail
		changePayload := SanityChangePayload{
			PrisonerID:   p.ID,
			PreviousSan:  previousSanity,
			NewSanity:    p.Sanity,
			Delta:        -drain,
			Cause:        "NOISE",
			CauseEventID: noiseEvent.ID,
		}

		changeEvent := events.GameEvent{
			ID:        events.GenerateEventID(),
			Timestamp: time.Now(),
			Type:      events.EventTypeSanityChange,
			ActorID:   "SYSTEM_SANITY",
			TargetID:  p.ID,
			Payload:   changePayload,
			GameDay:   noiseEvent.GameDay,
		}

		ss.eventLog.Append(changeEvent)
		ss.logger.Event("SANITY_DRAIN", p.ID, 
			"Drain:"+string(rune('0'+drain/10))+string(rune('0'+drain%10))+" | New:"+string(rune('0'+p.Sanity/10))+string(rune('0'+p.Sanity%10)))
	}
}

// ProcessWithdrawal handles Simon's 5-day "Cold Turkey" mechanic.
func (ss *SanitySystem) ProcessWithdrawal(gameDay int) {
	for _, p := range ss.prisoners {
		if p.Archetype != prisoner.ArchetypeRedeemed {
			continue
		}

		p.DayInGame = gameDay
		sanityMod, _ := rules.ProcessWithdrawal(p)

		if sanityMod != 0 {
			previousSanity := p.Sanity
			p.Sanity += sanityMod
			if p.Sanity > 100 {
				p.Sanity = 100
			}
			if p.Sanity < 0 {
				p.Sanity = 0
			}

			// Emit event for withdrawal effect
			changePayload := SanityChangePayload{
				PrisonerID:  p.ID,
				PreviousSan: previousSanity,
				NewSanity:   p.Sanity,
				Delta:       sanityMod,
				Cause:       "WITHDRAWAL",
			}

			changeEvent := events.GameEvent{
				ID:        events.GenerateEventID(),
				Timestamp: time.Now(),
				Type:      events.EventTypeSanityChange,
				ActorID:   "SYSTEM_SANITY",
				TargetID:  p.ID,
				Payload:   changePayload,
				GameDay:   gameDay,
			}

			ss.eventLog.Append(changeEvent)
		}
	}
}
