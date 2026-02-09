// Package action provides the "hands" of Los Gemelos.
// T018: System Event Emitter - Executes decisions as immutable events.
//
// This package transforms Cognition decisions into concrete game actions.
// All actions are recorded in the EventLog for full auditability.
package action

import (
	"context"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/engine"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/network"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/cognition"
)

// Executor translates Cognition decisions into game actions.
type Executor struct {
	noiseManager *engine.NoiseManager
	eventLog     *events.EventLog
	wsHub        *network.Hub
	logger       *logger.Logger
	currentDay   int
}

// NewExecutor creates a new action executor.
func NewExecutor(
	nm *engine.NoiseManager,
	el *events.EventLog,
	hub *network.Hub,
	log *logger.Logger,
) *Executor {
	return &Executor{
		noiseManager: nm,
		eventLog:     el,
		wsHub:        hub,
		logger:       log,
	}
}

// SetCurrentDay updates the current game day for event logging.
func (e *Executor) SetCurrentDay(day int) {
	e.currentDay = day
}

// Execute performs the decided action and logs it to the EventLog.
func (e *Executor) Execute(ctx context.Context, decision *cognition.Decision) error {
	if !decision.IsApproved {
		e.logger.Warn("ACTION: Attempted to execute unapproved decision")
		return nil // Silently ignore unapproved decisions
	}

	// Log the decision itself as a meta-event
	e.logDecisionEvent(decision)

	// Execute based on action type
	switch decision.ActionType {
	case cognition.ActionNoise:
		return e.executeNoise(decision)
	case cognition.ActionResourceCut:
		return e.executeResourceCut(decision)
	case cognition.ActionRevealSecret:
		return e.executeReveal(decision)
	case cognition.ActionReward:
		return e.executeReward(decision)
	case cognition.ActionDoNothing:
		e.logger.Info("ACTION: Los Gemelos observe silently.")
		return nil
	default:
		e.logger.Warn("ACTION: Unknown action type: " + decision.ActionType)
		return nil
	}
}

// logDecisionEvent records the Twins' decision for audit.
func (e *Executor) logDecisionEvent(decision *cognition.Decision) {
	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventType("TWINS_DECISION"),
		ActorID:   "SYSTEM_TWINS",
		TargetID:  decision.Target,
		Payload: map[string]interface{}{
			"action_type":   decision.ActionType,
			"intensity":     decision.Intensity,
			"justification": decision.Justification,
			"reason":        decision.Reason,
		},
		GameDay:    e.currentDay,
		IsRevealed: false, // Decisions are hidden until revealed
	}

	e.eventLog.Append(event)
	e.logger.Event("TWINS_ACTION", "TWINS", decision.ActionType+" -> "+decision.Target)
}

// executeNoise triggers a noise torture event.
func (e *Executor) executeNoise(decision *cognition.Decision) error {
	reason := "TWINS_AUTONOMOUS:" + decision.Justification
	e.noiseManager.TriggerPunishment(decision.Target, reason)

	// Broadcast to connected players
	e.wsHub.BroadcastToGame("", network.Message{
		Type:      network.MsgTypeEvent,
		Timestamp: time.Now().Unix(),
		Payload: map[string]interface{}{
			"event_type": "TWINS_NOISE",
			"message":    "Los Gemelos han activado una tortura de ruido...",
			"intensity":  decision.Intensity,
		},
	})

	e.logger.Info("ACTION: Noise torture executed at intensity " + string(rune('0'+decision.Intensity)))
	return nil
}

// executeResourceCut simulates cutting water/food supply.
func (e *Executor) executeResourceCut(decision *cognition.Decision) error {
	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventType("RESOURCE_CUT"),
		ActorID:   "SYSTEM_TWINS",
		TargetID:  decision.Target,
		Payload: map[string]interface{}{
			"resource_type": "WATER", // Could be WATER, FOOD, LIGHT
			"duration_hours": decision.Intensity * 2,
			"justification": decision.Justification,
		},
		GameDay: e.currentDay,
	}

	e.eventLog.Append(event)

	// Broadcast warning
	e.wsHub.BroadcastToGame("", network.Message{
		Type:      network.MsgTypeEvent,
		Timestamp: time.Now().Unix(),
		Payload: map[string]interface{}{
			"event_type": "TWINS_RESOURCE_CUT",
			"message":    "Los Gemelos han cortado el suministro de agua...",
			"target":     decision.Target,
		},
	})

	e.logger.Info("ACTION: Resource cut executed on " + decision.Target)
	return nil
}

// executeReveal reveals a hidden secret to the audience.
func (e *Executor) executeReveal(decision *cognition.Decision) error {
	// Find unrevealed events and mark one as revealed
	allEvents := e.eventLog.Replay()
	
	for _, evt := range allEvents {
		if !evt.IsRevealed && evt.Type == events.EventTypeBetrayal {
			// Mark as revealed (in production, this would update DB)
			event := events.GameEvent{
				ID:        events.GenerateEventID(),
				Timestamp: time.Now(),
				Type:      events.EventType("SECRET_REVEALED"),
				ActorID:   "SYSTEM_TWINS",
				TargetID:  evt.ActorID, // The betrayer
				Payload: map[string]interface{}{
					"original_event_id": evt.ID,
					"revealed_type":     string(evt.Type),
					"justification":     decision.Justification,
				},
				GameDay:    e.currentDay,
				IsRevealed: true,
			}

			e.eventLog.Append(event)

			// Broadcast revelation
			e.wsHub.BroadcastToGame("", network.Message{
				Type:      network.MsgTypeEvent,
				Timestamp: time.Now().Unix(),
				Payload: map[string]interface{}{
					"event_type": "TWINS_REVELATION",
					"message":    "Los Gemelos han revelado un secreto...",
					"target":     evt.ActorID,
				},
			})

			e.logger.Info("ACTION: Secret revealed about " + evt.ActorID)
			return nil
		}
	}

	e.logger.Info("ACTION: No secrets to reveal")
	return nil
}

// executeReward grants a positive effect to a prisoner.
func (e *Executor) executeReward(decision *cognition.Decision) error {
	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventType("TWINS_REWARD"),
		ActorID:   "SYSTEM_TWINS",
		TargetID:  decision.Target,
		Payload: map[string]interface{}{
			"reward_type":   "SANITY_BOOST",
			"amount":        10 * decision.Intensity,
			"justification": decision.Justification,
		},
		GameDay:    e.currentDay,
		IsRevealed: true, // Rewards are public
	}

	e.eventLog.Append(event)

	e.wsHub.BroadcastToGame("", network.Message{
		Type:      network.MsgTypeEvent,
		Timestamp: time.Now().Unix(),
		Payload: map[string]interface{}{
			"event_type": "TWINS_REWARD",
			"message":    "Los Gemelos han premiado a un prisionero...",
			"target":     decision.Target,
		},
	})

	e.logger.Info("ACTION: Reward granted to " + decision.Target)
	return nil
}
