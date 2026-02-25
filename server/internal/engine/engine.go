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
	isolationSystem  *IsolationSystem
	pollingSystem    *PollingSystem
	patioSystem      *PatioSystem
	contrabandSystem *ContrabandSystem
	inventorySystem  *InventorySystem

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
		isolationSystem:  NewIsolationSystem(eventLog, log),
		pollingSystem:    NewPollingSystem(eventLog, log),
		patioSystem:      NewPatioSystem(eventLog, log),
		contrabandSystem: NewContrabandSystem(eventLog, log),
		inventorySystem:  NewInventorySystem(eventLog, log),

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

// OverrideTime allows external bootstrapping commands to set the internal clock directly.
func (e *Engine) OverrideTime(day, hour int, tickNumber int64) {
	e.ticker.SetTime(day, hour, tickNumber)
}

// RegisterPrisoner adds a new player to all relevant subsystems.
func (e *Engine) RegisterPrisoner(p *prisoner.Prisoner) {
	e.prisoners[p.ID] = p
	e.sanitySystem.RegisterPrisoner(p)
	e.socialSystem.RegisterPrisoner(p)
	e.chaosSystem.RegisterPrisoner(p)
	e.metabolismSystem.RegisterPrisoner(p)
	e.isolationSystem.RegisterPrisoner(p)
	e.pollingSystem.RegisterPrisoner(p)
	e.patioSystem.RegisterPrisoner(p)
	e.contrabandSystem.RegisterPrisoner(p)
	e.inventorySystem.RegisterPrisoner(p)
	e.logger.Info("Prisoner registered with engine sub-systems: " + p.ID)
}

// GetPrisoners returns a snapshot of the current state of all players.
// Used by the AI Perceiver to evaluate Dignity and Traits.
func (e *Engine) GetPrisoners() map[string]*prisoner.Prisoner {
	return e.prisoners
}

// GetCurrentTime returns the current in-game day and hour from the ticker.
func (e *Engine) GetCurrentTime() (int, int) {
	return e.ticker.GetCurrentTime()
}

// GetNoiseManager exposes the built-in NoiseManager for the Twins AI Executor.
func (e *Engine) GetNoiseManager() *NoiseManager {
	return e.noiseManager
}

// GetPollingSystem exposes the polling system for API endpoints.
func (e *Engine) GetPollingSystem() *PollingSystem {
	return e.pollingSystem
}

// GetEventLog exposes the event log for the client to inject player actions.
func (e *Engine) GetEventLog() *events.EventLog {
	return e.eventLog
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
				e.checkMedicalEvacuations() // F5: Simon Factor
			}
		}
	}
}

// checkMedicalEvacuations enforces the Simon Factor (F5)
func (e *Engine) checkMedicalEvacuations() {
	for _, p := range e.prisoners {
		// If already dead/evacuated, skip
		if p.HP <= 0 && p.HasState(prisoner.StateDead) {
			continue
		}

		if p.HP <= 0 || p.Sanity <= 0 {
			e.logger.Error("MEDICAL EVACUATION TRIGGERED FOR " + p.Name)
			p.AddState(prisoner.StateDead, 0) // Permanent

			// Optional: Emitting a formal event could be done here if needed
			// Let's just log it and mark the state for now.
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
		e.isolationSystem.OnTimeTick(event)
		e.patioSystem.OnTimeTick(event)

		// Unmarshal payload if we need it for NoiseManager specifically
		if payload, ok := event.Payload.(TimeTickPayload); ok {
			e.noiseManager.OnTimeTick(payload)
		}

	case events.EventTypeIsolationChanged:
		e.isolationSystem.OnIsolationChanged(event)

	case events.EventTypeDoorLock, events.EventTypeDoorOpen:
		e.metabolismSystem.OnDoorLockEvent(event)

	case events.EventTypeNoiseEvent:
		e.sanitySystem.OnNoiseEvent(event)
		e.metabolismSystem.OnSleepInterruptEvent(event)

	case events.EventTypeAudioTorture:
		e.sanitySystem.OnAudioTortureEvent(event)
		e.metabolismSystem.OnSleepInterruptEvent(event)

	case events.EventTypePollCreated:
		e.pollingSystem.OnPollCreated(event)

	case events.EventTypePollResolved:
		e.pollingSystem.OnPollResolved(event)

	case events.EventTypeToiletUse:
		e.sanitySystem.OnToiletUseEvent(event)

	case events.EventTypeItemConsumed:
		e.metabolismSystem.OnItemConsumed(event)

	case events.EventTypeElixirGiven:
		e.metabolismSystem.OnElixirGiven(event)

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

	case events.EventTypeAudienceExpulsion:
		// Expel the prisoner immediately
		if p, ok := e.prisoners[event.TargetID]; ok {
			e.logger.Error("AUDIENCE EXPULSION EXECUTED FOR " + p.Name)
			p.AddState(prisoner.StateDead, 0)
		}
	}
}
