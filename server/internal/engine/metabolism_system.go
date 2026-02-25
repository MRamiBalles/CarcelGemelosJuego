package engine

import (
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/item"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// MetabolismSystem manages hunger, thirst, and class-specific biological rules.
type MetabolismSystem struct {
	eventLog          *events.EventLog
	logger            *logger.Logger
	prisoners         map[string]*prisoner.Prisoner
	lastHourProcessed int
}

// Removed ResourceIntakePayload (moved to inventory_system.go)

// NewMetabolismSystem creates a new metabolism manager.
func NewMetabolismSystem(eventLog *events.EventLog, log *logger.Logger) *MetabolismSystem {
	return &MetabolismSystem{
		eventLog:          eventLog,
		logger:            log,
		prisoners:         make(map[string]*prisoner.Prisoner),
		lastHourProcessed: -1,
	}
}

// RegisterPrisoner adds a prisoner to be tracked.
func (ms *MetabolismSystem) RegisterPrisoner(p *prisoner.Prisoner) {
	ms.prisoners[p.ID] = p
}

// OnTimeTick updates vital stats based on time passage.
func (ms *MetabolismSystem) OnTimeTick(event events.GameEvent) {
	payload, ok := event.Payload.(TimeTickPayload)
	if !ok {
		return
	}

	// We process metabolism per game hour to keep the math predictable
	if payload.GameHour == ms.lastHourProcessed {
		return
	}
	ms.lastHourProcessed = payload.GameHour

	for _, p := range ms.prisoners {
		// 1. Stamina and Sleep Mechanics
		if p.HasState(prisoner.StateAsleep) {
			// Regenerate stamina if fully resting and not starving
			if p.Hunger > 0 {
				p.Stamina += 6 // 60 per night (10 hours)
				if p.Stamina > 100 {
					p.Stamina = 100
				}
				// Remove Exhausted state if we regained enough stamina
				if p.Stamina > 10 && p.HasState(prisoner.StateExhausted) {
					delete(p.States, prisoner.StateExhausted)
				}
			} else {
				ms.logger.Event("POOR_SLEEP", p.ID, "Cannot rest due to starvation")
			}
		} else {
			// Decay stamina while awake
			staminaDrain := 3
			if p.HasTrait(prisoner.TraitInsomniac) {
				staminaDrain = 1 // Aída needs less sleep
			}
			p.Stamina -= staminaDrain
			if p.Stamina <= 0 {
				p.Stamina = 0
				p.AddState(prisoner.StateExhausted, 9999) // Indefinite until sleeps
				ms.logger.Warn("EXHAUSTION: " + p.Name + " has collapsed from fatigue!")
			}
		}

		// 2. Hydration Drain (Fast)
		p.Thirst -= 5
		if p.Thirst <= 0 {
			p.Thirst = 0
			p.HP -= 10 // Dehydration damage
			ms.logger.Warn("DEHYDRATION: " + p.Name + " is taking damage!")
		}

		// 3. Starvation Logic (Moderate)
		if !p.HasTrait(prisoner.TraitBreatharian) {
			p.Hunger -= 2
			if p.Hunger <= 0 {
				p.Hunger = 0
				p.HP -= 5 // Starvation damage
				ms.logger.Warn("STARVATION: " + p.Name + " is taking damage!")
			}
		} else {
			// Mystic (Tartaria): 21 Day Water Fasting
			// No food intake required. Hunger decays but causes no HP damage.
			p.Hunger -= 2
			if p.Hunger < 0 {
				p.Hunger = 0
			}
		}

		// Death Check
		if p.HP <= 0 {
			ms.logger.Warn("CRITICAL: " + p.ID + " requires medical evacuation! HP reached 0.")
		}
	}
}

// OnDoorLockEvent handles sleep cycle initiation and termination.
func (ms *MetabolismSystem) OnDoorLockEvent(event events.GameEvent) {
	payload, ok := event.Payload.(events.DoorLockPayload)
	if !ok {
		return
	}

	for _, p := range ms.prisoners {
		if payload.CellID == "ALL" || payload.CellID == p.CellID {
			if payload.IsLocked {
				p.AddState(prisoner.StateAsleep, 9999)
				ms.logger.Info("SLEEP: " + p.Name + " goes to sleep.")
			} else {
				delete(p.States, prisoner.StateAsleep)
				ms.logger.Info("WAKE: " + p.Name + " woke up.")
			}
		}
	}
}

// OnSleepInterruptEvent handles waking up prisoners due to noise or torture.
func (ms *MetabolismSystem) OnSleepInterruptEvent(event events.GameEvent) {
	if event.Type == events.EventTypeAudioTorture {
		for _, p := range ms.prisoners {
			if p.HasState(prisoner.StateAsleep) {
				delete(p.States, prisoner.StateAsleep)
				ms.logger.Warn("SLEEP_INTERRUPTED: " + p.Name + " was awakened by audio torture!")
			}
		}
		return
	}

	if event.Type == events.EventTypeNoiseEvent {
		payload, ok := event.Payload.(NoiseEventPayload)
		if ok {
			for _, p := range ms.prisoners {
				if payload.TargetZone == "ALL" || payload.TargetZone == p.CellID {
					if p.HasState(prisoner.StateAsleep) {
						delete(p.States, prisoner.StateAsleep)
						ms.logger.Warn("SLEEP_INTERRUPTED: " + p.Name + " was awakened by noise!")
					}
				}
			}
		}
	}
}

// OnItemConsumed handles the effects of eating/drinking items.
func (ms *MetabolismSystem) OnItemConsumed(event events.GameEvent) {
	payload, ok := event.Payload.(ItemConsumedPayload)
	if !ok {
		ms.logger.Error("Failed to parse ItemConsumedPayload")
		return
	}

	p, exists := ms.prisoners[payload.PrisonerID]
	if !exists {
		return
	}

	def, ok := item.GetItem(payload.ItemType)
	if !ok {
		ms.logger.Error("Unknown item type consumed: " + string(payload.ItemType))
		return
	}

	// Mystic Logic: Cannot eat solids
	if p.HasTrait(prisoner.TraitBreatharian) {
		if def.IsFood {
			// VIOLATION!
			p.Sanity -= 50
			p.HP -= 20
			ms.logger.Warn("MYSTIC VIOLATION: " + p.Name + " ate solid food!")
			return
		}
	}

	// Normal intake
	totalNutrition := def.Nutrition * payload.Quantity
	totalHydration := def.Hydration * payload.Quantity
	totalSanityMod := def.SanityMod * payload.Quantity

	p.Hunger += totalNutrition
	p.Thirst += totalHydration
	p.Sanity += totalSanityMod

	// Cap at 100
	if p.Thirst > 100 {
		p.Thirst = 100
	}
	if p.Hunger > 100 {
		p.Hunger = 100
	}
	if p.Sanity > 100 {
		p.Sanity = 100
	}
	// Cap at 0 bottom
	if p.Sanity < 0 {
		p.Sanity = 0
	}
}
