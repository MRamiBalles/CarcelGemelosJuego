// Package engine contains the game loop and simulation logic.
// This is the heartbeat of "Cárcel de los Gemelos".
//
// ARCHITECTURAL RULE: The Engine does NOT mutate Prisoner state directly.
// It emits TimeTickEvents to the EventLog. Subsystems subscribe and react.
package engine

import (
	"context"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// TickRate defines how often the game world updates (in real time).
const TickRate = 1 * time.Minute // 1 real minute = ~2 in-game hours

// TimeTickPayload is the data attached to each TimeTickEvent.
type TimeTickPayload struct {
	GameDay      int   `json:"game_day"`
	GameHour     int   `json:"game_hour"`      // 0-23 in-game
	TickNumber   int64 `json:"tick_number"`
	IsNightTime  bool  `json:"is_night_time"`  // 22:00-06:00
	IsMealWindow bool  `json:"is_meal_window"` // 08:00, 14:00, 20:00
}

// Ticker manages the game loop heartbeat.
// It does NOT know about Prisoners or Sanity - only time progression.
type Ticker struct {
	eventLog   *events.EventLog
	logger     *logger.Logger
	tickNumber int64
	gameDay    int
	gameHour   int
	stopChan   chan struct{}
}

// NewTicker creates a new game ticker.
func NewTicker(eventLog *events.EventLog, log *logger.Logger) *Ticker {
	return &Ticker{
		eventLog:   eventLog,
		logger:     log,
		tickNumber: 0,
		gameDay:    1,
		gameHour:   6, // Start at 6:00 AM
		stopChan:   make(chan struct{}),
	}
}

// Start begins the game loop. Call in a goroutine.
func (t *Ticker) Start(ctx context.Context) {
	t.logger.Info("Engine Ticker started. The Twins are watching...")

	ticker := time.NewTicker(TickRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.logger.Info("Engine Ticker stopped by context.")
			return
		case <-t.stopChan:
			t.logger.Info("Engine Ticker stopped manually.")
			return
		case <-ticker.C:
			t.tick()
		}
	}
}

// Stop gracefully stops the ticker.
func (t *Ticker) Stop() {
	close(t.stopChan)
}

// tick processes a single game tick.
func (t *Ticker) tick() {
	t.tickNumber++
	t.advanceTime()

	payload := TimeTickPayload{
		GameDay:      t.gameDay,
		GameHour:     t.gameHour,
		TickNumber:   t.tickNumber,
		IsNightTime:  t.gameHour >= 22 || t.gameHour < 6,
		IsMealWindow: t.gameHour == 8 || t.gameHour == 14 || t.gameHour == 20,
	}

	// Emit the TimeTickEvent - subsystems will react
	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeTimeTick,
		ActorID:   "SYSTEM_TWINS",
		Payload:   payload,
		GameDay:   t.gameDay,
	}

	t.eventLog.Append(event)
	t.logger.Event("TIME_TICK", "TWINS", 
		"Day "+string(rune('0'+t.gameDay))+" Hour "+string(rune('0'+t.gameHour/10))+string(rune('0'+t.gameHour%10)))
}

// advanceTime moves the in-game clock forward.
func (t *Ticker) advanceTime() {
	t.gameHour += 2 // Each tick = 2 in-game hours

	if t.gameHour >= 24 {
		t.gameHour = 0
		t.gameDay++

		if t.gameDay > 21 {
			t.logger.Warn("DAY 21 REACHED. Endgame triggered.")
			// TODO: Trigger Duo Dilemma resolution
		}
	}
}

// GetCurrentTime returns the current in-game time.
func (t *Ticker) GetCurrentTime() (day int, hour int) {
	return t.gameDay, t.gameHour
}
