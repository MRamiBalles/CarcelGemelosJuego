// Package engine - noise.go
// T008: Noise Event Generation System
//
// The Twins can trigger audio tortures at random intervals or on demand.
// This system DOES NOT modify Prisoner state - it emits NoiseEvents.
package engine

import (
	"math/rand"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// NoiseType defines the category of audio torture.
type NoiseType string

const (
	NoiseSiren      NoiseType = "SIREN"
	NoiseCryingBaby NoiseType = "CRYING_BABY"
	NoiseScratching NoiseType = "SCRATCHING"
	NoiseWhiteNoise NoiseType = "WHITE_NOISE"
	NoiseAlarm      NoiseType = "ALARM"
)

// NoiseEventPayload is the data for a noise torture event.
type NoiseEventPayload struct {
	NoiseType   NoiseType `json:"noise_type"`
	Intensity   int       `json:"intensity"`    // 1-3 (mild to severe)
	DurationSec int       `json:"duration_sec"` // How long it lasts
	TargetZone  string    `json:"target_zone"`  // "ALL", "BLOCK_A", specific cell ID
	Reason      string    `json:"reason"`       // For audit: "RANDOM", "PUNISHMENT", "AUDIENCE_VOTE"
}

// NoiseManager handles the generation and scheduling of noise events.
type NoiseManager struct {
	eventLog     *events.EventLog
	logger       *logger.Logger
	noiseChance  float64 // Probability per tick (0.0 - 1.0)
	currentDay   int
}

// NewNoiseManager creates a new noise manager.
func NewNoiseManager(eventLog *events.EventLog, log *logger.Logger) *NoiseManager {
	return &NoiseManager{
		eventLog:    eventLog,
		logger:      log,
		noiseChance: 0.15, // 15% chance per tick
	}
}

// OnTimeTick is called by the event subscriber when a TimeTickEvent occurs.
// This is the event-driven hook that avoids direct coupling to the Ticker.
func (nm *NoiseManager) OnTimeTick(tickPayload TimeTickPayload) {
	nm.currentDay = tickPayload.GameDay

	// Higher chance of noise at night (The Twins are cruel)
	chance := nm.noiseChance
	if tickPayload.IsNightTime {
		chance = 0.30 // 30% at night
	}

	// Day 15+ ramps up the pressure
	if nm.currentDay >= 15 {
		chance += 0.10
	}

	if rand.Float64() < chance {
		nm.triggerRandomNoise("RANDOM")
	}
}

// TriggerPunishment allows The Twins (or Audience) to manually trigger noise.
func (nm *NoiseManager) TriggerPunishment(targetZone string, reason string) {
	nm.generateNoiseEvent(NoiseSiren, 3, 60, targetZone, reason)
}

// triggerRandomNoise generates a random noise event.
func (nm *NoiseManager) triggerRandomNoise(reason string) {
	noiseTypes := []NoiseType{NoiseSiren, NoiseCryingBaby, NoiseScratching, NoiseWhiteNoise, NoiseAlarm}
	selectedNoise := noiseTypes[rand.Intn(len(noiseTypes))]
	intensity := rand.Intn(3) + 1      // 1-3
	duration := 30 + rand.Intn(90)     // 30-120 seconds

	nm.generateNoiseEvent(selectedNoise, intensity, duration, "ALL", reason)
}

// generateNoiseEvent creates and logs an immutable noise event.
func (nm *NoiseManager) generateNoiseEvent(noiseType NoiseType, intensity int, duration int, target string, reason string) {
	payload := NoiseEventPayload{
		NoiseType:   noiseType,
		Intensity:   intensity,
		DurationSec: duration,
		TargetZone:  target,
		Reason:      reason,
	}

	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeNoiseEvent,
		ActorID:   "SYSTEM_TWINS",
		Payload:   payload,
		GameDay:   nm.currentDay,
	}

	nm.eventLog.Append(event)
	nm.logger.Event("NOISE_TORTURE", "TWINS", 
		string(noiseType)+" | Intensity:"+string(rune('0'+intensity))+" | Target:"+target+" | Reason:"+reason)
}
