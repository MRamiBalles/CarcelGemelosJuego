package engine

import (
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/item"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// MetabolismSystem manages hunger, thirst, and class-specific biological rules.
type MetabolismSystem struct {
	eventLog  *events.EventLog
	logger    *logger.Logger
	prisoners map[string]*prisoner.Prisoner
}

// Removed ResourceIntakePayload (moved to inventory_system.go)

// NewMetabolismSystem creates a new metabolism manager.
func NewMetabolismSystem(eventLog *events.EventLog, log *logger.Logger) *MetabolismSystem {
	return &MetabolismSystem{
		eventLog:  eventLog,
		logger:    log,
		prisoners: make(map[string]*prisoner.Prisoner),
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

	// Only update once per hour to avoid spam/float issues
	if payload.TickNumber%60 != 0 { // Assuming 60 ticks = 1 hour? Or just use GameHour change?
		// The payload gives GameHour. We should update on hour change.
		// But TimeTick is every minute.
		// Let's rely on simple decay per tick for smoothness or per hour.
		// Design says "Hardcore Timeline".
		// Let's do small decrement every tick (1 min real = 2 min game).
		// 100 Hunger / (21 days) is too slow.
		// 100 Hunger / (3 days) = ~33 per day.
		// ~1.5 per hour.
		// If TickRate is 1 min, and each tick is 2 game hours... wait.
		// Ticker.go says: "gameHour += 2". So 1 real minute = 2 game hours.
		// So 12 ticks = 1 day.
		// That's ultra fast. 21 days = 21 * 12 mins = 4 hours.
		// Decay needs to be aggressive.
	}

	for _, p := range ms.prisoners {
		// Mystic: Breatharian - No hunger decay, but Stamina decay
		if p.HasTrait(prisoner.TraitBreatharian) {
			p.Stamina -= 1 // Decay stamina slightly
			if p.Stamina < 0 {
				p.Stamina = 0
				// Maybe add "Weakness" state?
			}
			continue
		}

		// Normal: Decay Hunger/Thirst
		p.Hunger -= 2
		p.Thirst -= 3

		// Starvation check
		if p.Hunger <= 0 {
			p.Hunger = 0
			p.HP -= 5 // Starvation damage
		}
		if p.Thirst <= 0 {
			p.Thirst = 0
			p.HP -= 10 // Dehydration is faster
		}

		// Death check (handled by HealthSystem? Or here?)
		if p.HP <= 0 {
			// Trigger DEATH event?
			ms.logger.Warn("PLAYER DYING: " + p.ID)
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
