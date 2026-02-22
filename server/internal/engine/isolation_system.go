package engine

import (
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// IsolationChangePayload contains the target and state of isolation.
type IsolationChangePayload struct {
	TargetID   string `json:"target_id"`
	IsIsolated bool   `json:"is_isolated"`
}

// IsolationSystem handles the 24h punishment cell.
type IsolationSystem struct {
	prisoners map[string]*prisoner.Prisoner
	eventLog  *events.EventLog
	logger    *logger.Logger
}

func NewIsolationSystem(eventLog *events.EventLog, log *logger.Logger) *IsolationSystem {
	is := &IsolationSystem{
		prisoners: make(map[string]*prisoner.Prisoner),
		eventLog:  eventLog,
		logger:    log,
	}

	eventLog.Subscribe(events.EventTypeIsolationChanged, is.OnIsolationChanged)
	eventLog.Subscribe(events.EventTypeTimeTick, is.OnTimeTick)

	return is
}

func (is *IsolationSystem) RegisterPrisoner(p *prisoner.Prisoner) {
	is.prisoners[p.ID] = p
}

func (is *IsolationSystem) OnIsolationChanged(event events.GameEvent) {
	payload, ok := event.Payload.(IsolationChangePayload)
	if !ok {
		// Attempt to parse map[string]interface{} mostly from JSON DB recovery
		m, ok := event.Payload.(map[string]interface{})
		if ok {
			payload.TargetID, _ = m["target_id"].(string)
			payload.IsIsolated, _ = m["is_isolated"].(bool)
		} else {
			return
		}
	}

	p, exists := is.prisoners[payload.TargetID]
	if !exists {
		return
	}

	p.IsIsolated = payload.IsIsolated
	if p.IsIsolated {
		is.logger.Event("ISOLATION_ENTER", p.ID, "Prisoner sent to solitary")
	} else {
		is.logger.Event("ISOLATION_EXIT", p.ID, "Prisoner released from solitary")
	}
}

func (is *IsolationSystem) OnTimeTick(event events.GameEvent) {
	// Every tick (e.g., 2 in-game hours), process solitary effects
	tickPayload, ok := event.Payload.(TimeTickPayload)
	if !ok {
		// check map variant
		m, mapped := event.Payload.(map[string]interface{})
		if mapped {
			dayFloat, _ := m["game_day"].(float64)
			tickPayload.GameDay = int(dayFloat)
		} else {
			return
		}
	}

	for _, p := range is.prisoners {
		if !p.IsIsolated {
			continue
		}

		// Apply Isolation modifiers
		var sanityDelta int
		cause := "SOLITARY_CONFINEMENT"

		if p.Archetype == prisoner.ArchetypeVeteran {
			// Frank loves being alone
			sanityDelta = 5
			cause = "MISANTHROPE_ISOLATION"
		} else if p.Archetype == prisoner.ArchetypeToxic {
			// Ylenia/Labrador go crazy without drama
			sanityDelta = -5
			cause = "TOXIC_ISOLATION_WITHDRAWAL"
		} else {
			// standard prisoner goes slightly crazy
			sanityDelta = -2
		}

		oldSanity := p.Sanity
		p.Sanity += sanityDelta
		if p.Sanity > 100 {
			p.Sanity = 100
		}
		if p.Sanity < 0 {
			p.Sanity = 0
		}

		// Emit Sanity Change Event
		changeEvent := events.GameEvent{
			ID:        events.GenerateEventID(),
			Timestamp: time.Now(),
			Type:      events.EventTypeSanityChange,
			ActorID:   "SYSTEM_ISOLATION",
			TargetID:  p.ID,
			GameDay:   tickPayload.GameDay,
			Payload: SanityChangePayload{
				PrisonerID:   p.ID,
				PreviousSan:  oldSanity,
				NewSanity:    p.Sanity,
				Delta:        sanityDelta,
				Cause:        cause,
				CauseEventID: event.ID,
			},
		}

		is.eventLog.Append(changeEvent)
	}
}
