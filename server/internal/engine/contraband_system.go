package engine

import (
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// LootEventPayload describes contraband discovered dynamically.
type LootEventPayload struct {
	TargetID  string `json:"target_id"`
	ItemName  string `json:"item_name"`
	SanityBuf int    `json:"sanity_buf"`
}

// SnitchEventPayload handles players betraying others.
type SnitchEventPayload struct {
	ActorID       string `json:"actor_id"`
	TargetID      string `json:"target_id"`
	Success       bool   `json:"success"`
	PotStolen     int    `json:"pot_stolen"`
	SanityPenalty int    `json:"sanity_penalty"`
}

// ContrabandSystem manages hidden loot, buffering player stats illicitly, and the ActionSnitch command.
type ContrabandSystem struct {
	prisoners map[string]*prisoner.Prisoner
	eventLog  *events.EventLog
	logger    *logger.Logger
	// Track who has contraband: TargetID -> true
	hasContraband map[string]bool
}

func NewContrabandSystem(eventLog *events.EventLog, log *logger.Logger) *ContrabandSystem {
	cs := &ContrabandSystem{
		prisoners:     make(map[string]*prisoner.Prisoner),
		eventLog:      eventLog,
		logger:        log,
		hasContraband: make(map[string]bool),
	}

	return cs
}

func (cs *ContrabandSystem) RegisterPrisoner(p *prisoner.Prisoner) {
	cs.prisoners[p.ID] = p
}

// GenerateLoot drops a hidden item randomly (e.g., called by Twins API or Cron).
func (cs *ContrabandSystem) GenerateLoot(targetID string, itemName string, sanityBuf int) {
	target, exists := cs.prisoners[targetID]
	if !exists {
		return
	}

	cs.hasContraband[targetID] = true
	target.Sanity += sanityBuf
	if target.Sanity > 100 {
		target.Sanity = 100
	}

	// Logging an obscured event to not expose it directly to full clients
	cs.logger.Event("LOOT_ACQUIRED", target.ID, "Found hidden item: "+itemName)

	// Emit event for event-sourcing (Reality Recap)
	cs.eventLog.Append(events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeLootAcquired,
		ActorID:   "SYSTEM_GEMELOS",
		TargetID:  targetID,
		Payload: LootEventPayload{
			TargetID:  targetID,
			ItemName:  itemName,
			SanityBuf: sanityBuf,
		},
	})
}

// OnSocialAction intercepts ActionSnitch
func (cs *ContrabandSystem) OnSocialAction(event events.GameEvent) {
	// Need to check if the action was a Snitch action
	payload, ok := event.Payload.(SocialActionPayload)
	if !ok {
		// Map decoding
		m, mapped := event.Payload.(map[string]interface{})
		if mapped {
			payload.ActorID, _ = m["actor_id"].(string)
			payload.TargetID, _ = m["target_id"].(string)
			payload.ActionType, _ = m["action_type"].(string)
		} else {
			return
		}
	}

	if payload.ActionType != "ActionSnitch" {
		return
	}

	actor, actorExists := cs.prisoners[payload.ActorID]
	target, targetExists := cs.prisoners[payload.TargetID]

	if !actorExists || !targetExists {
		return
	}

	actor.Loyalty -= 20
	if actor.Loyalty < -100 {
		actor.Loyalty = -100
	}

	snitchPayload := SnitchEventPayload{
		ActorID:  actor.ID,
		TargetID: target.ID,
	}

	// Evaluate Snitch
	if cs.hasContraband[target.ID] {
		// Success: Target is caught. Snitch gets 50% of target Pot (simulated here as random static value for now)
		snitchPayload.Success = true
		snitchPayload.PotStolen = 500

		target.Sanity -= 40
		if target.Sanity < 0 {
			target.Sanity = 0
		}
		delete(cs.hasContraband, target.ID)

		cs.logger.Event("SNITCH_SUCCESS", actor.ID, "Caught "+target.ID+" with contraband")

		// Instantly send target to Isolation Room
		cs.eventLog.Append(events.GameEvent{
			ID:        events.GenerateEventID(),
			Timestamp: time.Now(),
			Type:      events.EventTypeIsolationChanged,
			ActorID:   "SYSTEM_GEMELOS",
			TargetID:  target.ID,
			GameDay:   event.GameDay,
			Payload: IsolationChangePayload{
				TargetID:   target.ID,
				IsIsolated: true,
			},
		})

	} else {
		// Fail: Actor lied to the Twins
		snitchPayload.Success = false
		snitchPayload.SanityPenalty = 30

		actor.Sanity -= 30
		if actor.Sanity < 0 {
			actor.Sanity = 0
		}

		cs.logger.Event("SNITCH_FAIL", actor.ID, "Lied about "+target.ID+". Punished by Twins.")
	}

	// Record betrayal
	cs.eventLog.Append(events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeBetrayal,
		ActorID:   actor.ID,
		TargetID:  target.ID,
		GameDay:   event.GameDay,
		Payload:   snitchPayload,
	})
}
