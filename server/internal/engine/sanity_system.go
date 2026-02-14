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

// OnAudioTortureEvent handles inescapable audio torture.
// It reuses the noise logic but the client treats it differently (bypass volume).
func (ss *SanitySystem) OnAudioTortureEvent(event events.GameEvent) {
	ss.OnNoiseEvent(event)
}

// ToiletUsePayload carries data about a toilet usage action.
type ToiletUsePayload struct {
	ActorID string `json:"actor_id"`
	CellID  string `json:"cell_id"`
}

// OnToiletUseEvent handles the "Toilet of Shame" mechanic.
// If a prisoner uses the toilet while their cellmate is NOT facing the wall,
// both take massive sanity damage.
func (ss *SanitySystem) OnToiletUseEvent(event events.GameEvent) {
	payload, ok := event.Payload.(ToiletUsePayload)
	if !ok {
		ss.logger.Error("Failed to parse ToiletUsePayload")
		return
	}

	user, exists := ss.prisoners[payload.ActorID]
	if !exists {
		return
	}

	// 1. Apply Dignity Loss to User
	user.Dignity -= 15
	if user.Dignity < 0 {
		user.Dignity = 0
	}
	ss.logger.Event("DIGNITY_LOSS", user.ID, "Used toilet openly")

	// 2. Check for witnesses (Cellmates)
	for _, neighbor := range ss.prisoners {
		// Same cell, different person
		if neighbor.CellID == user.CellID && neighbor.ID != user.ID {
			// Check facing
			isFacingWall := neighbor.HasState(prisoner.StateFacingWall)
			
			if !isFacingWall {
				// PRIVACY BREACH!
				
				// Damage Witness
				drainWitness := 20
				// Mystic mitigation
				if neighbor.Archetype == prisoner.ArchetypeMystic && neighbor.Sanity > 20 {
					drainWitness /= 2
				}
				
				neighbor.Sanity -= drainWitness
				if neighbor.Sanity < 0 {
					neighbor.Sanity = 0
				}
				
				ss.emitSanityChange(neighbor.ID, -drainWitness, "PRIVACY_WITNESS", event.ID)

				// Damage User (Shame)
				drainUser := 10
				user.Sanity -= drainUser
				if user.Sanity < 0 {
					user.Sanity = 0
				}
				ss.emitSanityChange(user.ID, -drainUser, "PRIVACY_VIOLATION", event.ID)
			}
		}
	}
}

// emitSanityChange is a helper to log sanity changes
func (ss *SanitySystem) emitSanityChange(targetID string, delta int, cause string, causeID string) {
	p, exists := ss.prisoners[targetID]
	if !exists {
		return
	}

	payload := SanityChangePayload{
		PrisonerID:   targetID,
		PreviousSan:  p.Sanity - delta, // Approx
		NewSanity:    p.Sanity,
		Delta:        delta,
		Cause:        cause,
		CauseEventID: causeID,
	}

	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeSanityChange,
		ActorID:   "SYSTEM_SANITY",
		TargetID:  targetID,
		Payload:   payload,
		GameDay:   p.DayInGame,
	}

	ss.eventLog.Append(event)
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
			ss.emitSanityChange(p.ID, sanityMod, "WITHDRAWAL", "SYSTEM_TICK")
		}
	}
}
