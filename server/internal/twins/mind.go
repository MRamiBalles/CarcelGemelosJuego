// Package twins provides the orchestrator for Los Gemelos AI.
// This is the main entry point that coordinates Perception, Cognition, and Action.
package twins

import (
	"context"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/engine"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/infra/ai"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/network"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/action"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/cognition"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/perception"
)

// TwinsMind is the orchestrator for the AI overlords.
// It runs the Perception-Cognition-Action loop at regular intervals.
type TwinsMind struct {
	perceiver *perception.Perceiver
	cognitor  *cognition.Cognitor
	executor  *action.Executor
	logger    *logger.Logger
	gameID    string
	interval  time.Duration
	stopChan  chan struct{}
}

// NewTwinsMind creates the AI orchestrator.
func NewTwinsMind(
	eventLog *events.EventLog,
	stateProvider perception.StateProvider,
	llmProvider ai.LLMProvider,
	noiseManager *engine.NoiseManager,
	wsHub *network.Hub,
	log *logger.Logger,
) *TwinsMind {
	return &TwinsMind{
		perceiver: perception.NewPerceiver(eventLog, stateProvider, log),
		cognitor:  cognition.NewCognitor(llmProvider, log),
		executor:  action.NewExecutor(noiseManager, eventLog, wsHub, log),
		logger:    log,
		interval:  5 * time.Minute, // Evaluate every 5 real minutes
		stopChan:  make(chan struct{}),
	}
}

// SetGameID sets the current game for the AI to monitor.
func (tm *TwinsMind) SetGameID(gameID string) {
	tm.gameID = gameID
}

// Start begins the autonomous decision loop.
func (tm *TwinsMind) Start(ctx context.Context) {
	tm.logger.Info("Los Gemelos have awakened. They are watching...")

	ticker := time.NewTicker(tm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			tm.logger.Info("Los Gemelos enter dormancy.")
			return
		case <-tm.stopChan:
			tm.logger.Info("Los Gemelos manually silenced.")
			return
		case <-ticker.C:
			tm.runCycle(ctx)
		}
	}
}

// Stop halts the AI loop.
func (tm *TwinsMind) Stop() {
	close(tm.stopChan)
}

// runCycle executes one Perception-Cognition-Action cycle.
func (tm *TwinsMind) runCycle(ctx context.Context) {
	tm.logger.Info("Los Gemelos cycle: PERCEIVE -> DECIDE -> ACT")

	// 1. PERCEIVE: Build current state
	state, err := tm.perceiver.BuildPrisonState(ctx, tm.gameID, 1) // TODO: Get actual day
	if err != nil {
		tm.logger.Error("Perception failed: " + err.Error())
		return
	}

	// 2. COGNITION: Make a decision
	decision, err := tm.cognitor.Decide(ctx, state)
	if err != nil {
		tm.logger.Error("Cognition failed: " + err.Error())
		return
	}

	// 3. ACTION: Execute the decision
	if err := tm.executor.Execute(ctx, decision); err != nil {
		tm.logger.Error("Action failed: " + err.Error())
		return
	}

	tm.logger.Info("Los Gemelos cycle complete.")
}

// ForceDecision triggers an immediate evaluation (for testing/admin).
func (tm *TwinsMind) ForceDecision(ctx context.Context, currentDay int) (*cognition.Decision, error) {
	tm.executor.SetCurrentDay(currentDay)

	state, err := tm.perceiver.BuildPrisonState(ctx, tm.gameID, currentDay)
	if err != nil {
		return nil, err
	}

	decision, err := tm.cognitor.Decide(ctx, state)
	if err != nil {
		return nil, err
	}

	if err := tm.executor.Execute(ctx, decision); err != nil {
		return nil, err
	}

	return decision, nil
}
