package engine

import (
	"context"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// Engine is the central orchestrator that wires up the Event Sourcing log to the game mechanics.
type Engine struct {
	eventLog *events.EventLog
	logger   *logger.Logger
	ticker   *Ticker

	// Sub-systems
	sanitySystem     *SanitySystem
	socialSystem     *SocialSystem
	chaosSystem      *ChaosSystem
	metabolismSystem *MetabolismSystem
	lockdownSystem   *LockdownSystem
	noiseManager     *NoiseManager

	// State
	lastProcessedEvent int
	prisoners          map[string]*prisoner.Prisoner
}

// NewEngine initializes the core game systems and dependencies.
func NewEngine(eventLog *events.EventLog, log *logger.Logger) *Engine {
	e := &Engine{
		eventLog: eventLog,
		logger:   log,
		ticker:   NewTicker(eventLog, log),

		sanitySystem:     NewSanitySystem(eventLog, log),
		socialSystem:     NewSocialSystem(eventLog, log),
		chaosSystem:      NewChaosSystem(eventLog, log),
		metabolismSystem: NewMetabolismSystem(eventLog, log),
		lockdownSystem:   NewLockdownSystem(eventLog, log),
		noiseManager:     NewNoiseManager(eventLog, log),

		lastProcessedEvent: 0,
		prisoners:          make(map[string]*prisoner.Prisoner),
	}

	return e
}

// Start spawns the Ticker and the EventProcessor loop.
func (e *Engine) Start(ctx context.Context) {
	e.logger.Info("Starting core game engine...")

	// Start the main game clock
	go e.ticker.Start(ctx)

	// Start the event processing loop
	go e.processEvents(ctx)
}

// RegisterPrisoner adds a new player to all relevant subsystems.
func (e *Engine) RegisterPrisoner(p *prisoner.Prisoner) {
	e.prisoners[p.ID] = p
	e.sanitySystem.RegisterPrisoner(p)
	e.socialSystem.RegisterPrisoner(p)
	e.chaosSystem.RegisterPrisoner(p)
	e.metabolismSystem.RegisterPrisoner(p)
	e.logger.Info("Prisoner registered with engine sub-systems: " + p.ID)
}

// GetPrisoners returns a snapshot of the current state of all players.
// Used by the AI Perceiver to evaluate Dignity and Traits.
func (e *Engine) GetPrisoners() map[string]*prisoner.Prisoner {
	return e.prisoners
}

// processEvents listens to the EventLog and dispatches items to subsystems.
func (e *Engine) processEvents(ctx context.Context) {
	pollInterval := time.NewTicker(100 * time.Millisecond) // Poll the event log for new events
	defer pollInterval.Stop()

	for {
		select {
		case <-ctx.Done():
			e.logger.Info("EventProcessor stopped.")
			return
		case <-pollInterval.C:
			allEvents := e.eventLog.Replay()
			newEventsCount := len(allEvents) - e.lastProcessedEvent

			if newEventsCount > 0 {
				newEvents := allEvents[e.lastProcessedEvent:]
				for _, event := range newEvents {
					e.dispatch(event)
				}
				e.lastProcessedEvent = len(allEvents)
			}
		}
	}
}

// dispatch routes a standard GameEvent to the appropriate subsystems based on its type.
func (e *Engine) dispatch(event events.GameEvent) {
	switch event.Type {
	case events.EventTypeTimeTick:
		e.lockdownSystem.OnTimeTick(event)
		e.metabolismSystem.OnTimeTick(event)
		e.sanitySystem.OnTimeTick(event)
		e.socialSystem.OnTimeTick(event)
		e.chaosSystem.OnTimeTick(event)

		// Unmarshal payload if we need it for NoiseManager specifically
		if payload, ok := event.Payload.(TimeTickPayload); ok {
			e.noiseManager.OnTimeTick(payload)
		}

	case events.EventTypeNoiseEvent:
		e.sanitySystem.OnNoiseEvent(event)

	case events.EventTypeAudioTorture:
		e.sanitySystem.OnAudioTortureEvent(event)

	case events.EventTypeToiletUse:
		e.sanitySystem.OnToiletUseEvent(event)

	case events.EventTypeResourceIntake:
		e.metabolismSystem.OnResourceIntake(event)

	case events.EventTypeInsult:
		e.sanitySystem.OnInsultEvent(event)

	case events.EventTypeAggressiveEmote:
		e.socialSystem.OnAggressiveEmote(event)

	case events.EventTypeSteal:
		e.chaosSystem.OnStealEvent(event)

	case events.EventTypeLockdownBang:
		e.chaosSystem.OnLockdownBang(event)

	case events.EventTypeFinalDilemmaDecision:
		e.socialSystem.OnFinalDilemmaDecision(event)
	}
}
