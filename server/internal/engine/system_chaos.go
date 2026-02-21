package engine

import (
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// ChaosSystem manages active abilities for Chaos Agents and Deceivers.
type ChaosSystem struct {
	eventLog      *events.EventLog
	logger        *logger.Logger
	prisoners     map[string]*prisoner.Prisoner
	delayedEvents map[string]int // EventID -> GameHour it should be revealed
}

// StealPayload holds data for a theft action.
type StealPayload struct {
	ActorID  string `json:"actor_id"`
	TargetID string `json:"target_id"` // Victim or Container
	ItemType string `json:"item_type"` // e.g., "RICE", "WATER"
	Amount   int    `json:"amount"`
}

// LockdownBangPayload runs Aída's Poltergeist ability.
type LockdownBangPayload struct {
	ActorID      string `json:"actor_id"`
	TargetCellID string `json:"target_cell_id"`
}

// NewChaosSystem creates a new system for disruptive mechanics.
func NewChaosSystem(eventLog *events.EventLog, log *logger.Logger) *ChaosSystem {
	return &ChaosSystem{
		eventLog:      eventLog,
		logger:        log,
		prisoners:     make(map[string]*prisoner.Prisoner),
		delayedEvents: make(map[string]int),
	}
}

// RegisterPrisoner adds a prisoner to be tracked.
func (cs *ChaosSystem) RegisterPrisoner(p *prisoner.Prisoner) {
	cs.prisoners[p.ID] = p
}

// OnTimeTick checks for delayed events (Héctor's Smooth Criminal mechanic).
func (cs *ChaosSystem) OnTimeTick(event events.GameEvent) {
	payload, ok := event.Payload.(TimeTickPayload)
	if !ok {
		return
	}

	currentTime := (payload.GameDay * 24) + payload.GameHour

	for eventID, revealTime := range cs.delayedEvents {
		if currentTime >= revealTime {
			// Time to reveal the steal event
			cs.revealEvent(eventID)
			delete(cs.delayedEvents, eventID)
			cs.logger.Event("SECRET_REVEALED", "TWINS", "Delayed event "+eventID+" is now public.")
		}
	}
}

// OnStealEvent processes theft and checks for TraitSmoothCriminal to delay the VAR log.
func (cs *ChaosSystem) OnStealEvent(event events.GameEvent) {
	payload, ok := event.Payload.(StealPayload)
	if !ok {
		return
	}

	actor, exists := cs.prisoners[payload.ActorID]
	if !exists {
		return
	}

	isRevealed := true

	// Deceiver (Héctor): Smooth Criminal trait delays the VAR reveal by 12 game hours (T037)
	if actor.HasTrait(prisoner.TraitSmoothCriminal) {
		isRevealed = false
		currentHour := (event.GameDay * 24) + 6 // Rough estimation, should ideally pass current exact time
		cs.delayedEvents[event.ID] = currentHour + 12
		cs.logger.Info("SmoothCriminal triggered for " + actor.ID + ": Event delayed 12h")
	}

	// Persist the Steal Event
	stealEvent := events.GameEvent{
		ID:         event.ID,
		Timestamp:  time.Now(),
		Type:       events.EventTypeSteal,
		ActorID:    actor.ID,
		TargetID:   payload.TargetID,
		Payload:    payload,
		GameDay:    event.GameDay,
		IsRevealed: isRevealed,
	}

	cs.eventLog.Append(stealEvent)
}

// OnLockdownBang handles Aída's Poltergeist ability during Lockdown.
func (cs *ChaosSystem) OnLockdownBang(event events.GameEvent) {
	payload, ok := event.Payload.(LockdownBangPayload)
	if !ok {
		return
	}

	actor, exists := cs.prisoners[payload.ActorID]
	if !exists {
		return
	}

	// Ensure the ability is only casted by someone with the trait
	if !actor.HasTrait(prisoner.TraitInsomniac) { // Approximating Poltergeist to Insomniac as they are linked
		cs.logger.Warn("Unauthorized LockdownBang by " + actor.ID)
		return
	}

	// This action emits a NoiseEvent targeted at the cell
	noisePayload := NoiseEventPayload{
		NoiseType:   "BANGING_BARS",
		Intensity:   3, // High local intensity
		DurationSec: 60,
		TargetZone:  payload.TargetCellID,
		Reason:      "POLTERGEIST",
	}

	noiseEvent := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeNoiseEvent,
		ActorID:   actor.ID,
		TargetID:  payload.TargetCellID,
		Payload:   noisePayload,
		GameDay:   event.GameDay,
	}

	// Appending this will inherently trigger SanitySystem.OnNoiseEvent
	cs.eventLog.Append(noiseEvent)
	cs.logger.Event("POLTERGEIST_TRIGGERED", actor.ID, "Target:"+payload.TargetCellID)
}

// revealEvent modifies the event log (in a real DB this would UPDATE the row).
func (cs *ChaosSystem) revealEvent(eventID string) {
	// In an append-only system, we can't truly UPDATE.
	// We append a REVEAL event instead.
	// The frontend/VarsReplay will interpret this correctly to display the old hidden event.
	revealEvent := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      "SECRET_REVEALED",
		ActorID:   "SYSTEM_TWINS",
		TargetID:  eventID, // Pointing to the hidden event
		Payload: map[string]interface{}{
			"revealed_event_id": eventID,
		},
		IsRevealed: true,
	}
	cs.eventLog.Append(revealEvent)
}
