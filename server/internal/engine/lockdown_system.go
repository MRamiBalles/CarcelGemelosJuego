package engine

import (
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// LockdownSystem manages the automated door schedule.
// It listens to TimeTickEvents and triggers DoorLock/DoorOpen actions.
type LockdownSystem struct {
	eventLog *events.EventLog
	logger   *logger.Logger
}

// NewLockdownSystem creates a new lockdown manager.
func NewLockdownSystem(eventLog *events.EventLog, log *logger.Logger) *LockdownSystem {
	return &LockdownSystem{
		eventLog: eventLog,
		logger:   log,
	}
}

// OnTimeTick checks the game hour and triggers lockdown events.
func (ls *LockdownSystem) OnTimeTick(event events.GameEvent) {
	payload, ok := event.Payload.(TimeTickPayload)
	if !ok {
		return
	}

	// 22:00 -> LOCKDOWN
	if payload.GameHour == 22 {
		ls.logger.Info("LOCKDOWN: Closing all cells.")
		ls.emitLockEvent(events.EventTypeDoorLock, payload.GameDay)
	}

	// 08:00 -> UNLOCK
	if payload.GameHour == 8 {
		ls.logger.Info("UNLOCK: Opening all cells.")
		ls.emitLockEvent(events.EventTypeDoorOpen, payload.GameDay)
	}
}

func (ls *LockdownSystem) emitLockEvent(eventType events.EventType, day int) {
	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      eventType,
		ActorID:   "SYSTEM_LOCKDOWN",
		TargetID:  "ALL_CELLS",
		Payload:   map[string]string{"action": string(eventType)},
		GameDay:   day,
	}
	ls.eventLog.Append(event)
}
