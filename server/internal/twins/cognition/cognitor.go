// Package cognition provides the "brain" of Los Gemelos.
// T017: MAD-BAD-SAD Decision Framework for autonomous punishment/reward.
//
// MAD = Morally Absolute Denial (Things we NEVER do)
// BAD = Bounded Acceptable Damage (Things we CAN do within limits)
// SAD = Spectacle Amplification Directive (Things we WANT to do for drama)
package cognition

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/infra/ai"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/perception"
)

// Decision represents a planned action by Los Gemelos.
type Decision struct {
	ActionType    string                 `json:"action_type"`
	Target        string                 `json:"target"`        // "ALL", prisoner ID, or zone
	Intensity     int                    `json:"intensity"`     // 1-3
	Reason        string                 `json:"reason"`        // For audit trail
	IsApproved    bool                   `json:"is_approved"`   // Passed MAD check?
	Justification string                 `json:"justification"` // Why this decision?
	Metadata      map[string]interface{} `json:"metadata"`
}

// ActionType constants for Los Gemelos actions.
const (
	ActionNoise           = "NOISE_TORTURE"
	ActionAudioTorture    = "AUDIO_TORTURE"
	ActionLockdown        = "DOOR_LOCK"
	ActionResourceCut     = "RESOURCE_CUT"
	ActionRevealSecret    = "REVEAL_SECRET"
	ActionReward          = "REWARD"
	ActionRedPhoneMessage = "RED_PHONE_MESSAGE"
	ActionDoNothing       = "OBSERVE"
)

// Cognitor is the decision-making core of Los Gemelos.
type Cognitor struct {
	logger        *logger.Logger
	llm           ai.LLMProvider
	madRules      []MADRule
	sadObjectives []SADObjective
}

// MADRule defines an absolute prohibition (Morally Absolute Denial).
type MADRule struct {
	Name        string
	Description string
	CheckFunc   func(state *perception.PrisonState, action string) bool
}

// SADObjective defines a spectacle goal (Spectacle Amplification Directive).
type SADObjective struct {
	Name      string
	Priority  int // Higher = more important
	CheckFunc func(state *perception.PrisonState) bool
}

// NewCognitor creates a new cognition module with default MAD rules.
func NewCognitor(llm ai.LLMProvider, log *logger.Logger) *Cognitor {
	c := &Cognitor{
		logger: log,
		llm:    llm,
	}
	c.initializeMADRules()
	c.initializeSADObjectives()
	return c
}

// initializeMADRules sets up the absolute prohibitions.
func (c *Cognitor) initializeMADRules() {
	c.madRules = []MADRule{
		{
			Name:        "NO_KILL_LOW_SANITY",
			Description: "Never trigger noise on prisoners with <10% sanity",
			CheckFunc: func(state *perception.PrisonState, action string) bool {
				// If average sanity is critically low and action is noise, block it
				if action == ActionNoise && state.AverageSanity < 10 {
					return false // BLOCKED
				}
				return true // ALLOWED
			},
		},
		{
			Name:        "NO_DOUBLE_PUNISHMENT",
			Description: "Never punish the same target twice in one hour",
			CheckFunc: func(state *perception.PrisonState, action string) bool {
				// TODO: Check recent events for double punishment
				return true
			},
		},
		{
			Name:        "NO_DAY_ONE_CRUELTY",
			Description: "No punishments on Day 1 (grace period)",
			CheckFunc: func(state *perception.PrisonState, action string) bool {
				if state.CurrentDay == 1 && (action == ActionNoise || action == ActionResourceCut) {
					return false
				}
				return true
			},
		},
		{
			Name:        "REQUIRE_AUDIT_TRAIL",
			Description: "Every action must have a documented reason",
			CheckFunc: func(state *perception.PrisonState, action string) bool {
				return true // Always require justification (enforced in Decide)
			},
		},
	}
}

// initializeSADObjectives sets up the spectacle goals.
func (c *Cognitor) initializeSADObjectives() {
	c.sadObjectives = []SADObjective{
		{
			Name:     "MAINTAIN_TENSION",
			Priority: 3,
			CheckFunc: func(state *perception.PrisonState) bool {
				return state.TensionLevel == "LOW" || state.TensionLevel == "MEDIUM"
			},
		},
		{
			Name:     "REWARD_DRAMA",
			Priority: 2,
			CheckFunc: func(state *perception.PrisonState) bool {
				return state.RecentBetrayals > 0 // Drama happened, reward the chaos
			},
		},
		{
			Name:     "AUDIENCE_SATISFACTION",
			Priority: 1,
			CheckFunc: func(state *perception.PrisonState) bool {
				return state.AudienceActivity > 5 // Audience is engaged
			},
		},
	}
}

// Decide evaluates the current state and produces a decision.
func (c *Cognitor) Decide(ctx context.Context, state *perception.PrisonState) (*Decision, error) {
	// First, determine what action would best serve the SAD objectives
	proposedAction := c.selectAction(state)

	// Build the decision
	decision := &Decision{
		ActionType:    proposedAction,
		Target:        c.selectTarget(state, proposedAction),
		Intensity:     c.selectIntensity(state),
		Reason:        "TWINS_AUTONOMOUS",
		Justification: c.buildJustification(state, proposedAction),
		Metadata:      make(map[string]interface{}),
	}

	// Run MAD check - can this action be approved?
	decision.IsApproved = c.runMADCheck(state, proposedAction)

	if !decision.IsApproved {
		c.logger.Warn("COGNITION: Action blocked by MAD rule: " + proposedAction)
		decision.ActionType = ActionDoNothing
		decision.Justification = "Acción bloqueada por reglas MAD. Los Gemelos observan."
	}

	c.logger.Event("COGNITION", "TWINS",
		fmt.Sprintf("Decision: %s (Approved: %v)", decision.ActionType, decision.IsApproved))

	return decision, nil
}

// selectAction chooses the best action based on SAD objectives.
func (c *Cognitor) selectAction(state *perception.PrisonState) string {
	// Priority-based selection
	for _, obj := range c.sadObjectives {
		if obj.CheckFunc(state) {
			switch obj.Name {
			case "MAINTAIN_TENSION":
				// Tension is low, introduce chaos
				return ActionNoise
			case "REWARD_DRAMA":
				// Drama happened, reveal a secret to escalate
				return ActionRevealSecret
			case "AUDIENCE_SATISFACTION":
				// Audience is watching, give them a show
				if rand.Float32() > 0.5 {
					return ActionNoise
				}
				return ActionRevealSecret
			}
		}
	}

	// Default: observe
	return ActionDoNothing
}

// selectTarget picks the target for the action.
func (c *Cognitor) selectTarget(state *perception.PrisonState, action string) string {
	switch action {
	case ActionNoise:
		// Global noise affects everyone
		if state.TensionLevel == "CRITICAL" {
			return "ALL"
		}
		// TODO: Pick a specific block or prisoner
		return "BLOCK_A"
	case ActionRevealSecret:
		// Reveal secrets about the most suspicious prisoner
		// TODO: Analyze prisoner profiles
		return "RANDOM"
	default:
		return "NONE"
	}
}

// selectIntensity determines how severe the action should be.
func (c *Cognitor) selectIntensity(state *perception.PrisonState) int {
	switch state.TensionLevel {
	case "LOW":
		return 1 // Gentle nudge
	case "MEDIUM":
		return 2 // Moderate pressure
	case "HIGH", "CRITICAL":
		return 3 // Maximum drama
	default:
		return 1
	}
}

// runMADCheck verifies the action against all MAD rules.
func (c *Cognitor) runMADCheck(state *perception.PrisonState, action string) bool {
	for _, rule := range c.madRules {
		if !rule.CheckFunc(state, action) {
			c.logger.Warn("MAD VIOLATION: " + rule.Name)
			return false
		}
	}
	return true
}

// buildJustification creates an audit-friendly explanation.
func (c *Cognitor) buildJustification(state *perception.PrisonState, action string) string {
	switch action {
	case ActionNoise:
		return fmt.Sprintf("Tensión actual: %s. Se requiere estímulo para mantener el interés narrativo.", state.TensionLevel)
	case ActionRevealSecret:
		return fmt.Sprintf("Traiciones recientes: %d. La audiencia merece transparencia.", state.RecentBetrayals)
	case ActionRedPhoneMessage:
		return "Comunicado directo del Panóptico vía Teléfono Rojo para manipular la situación."
	case ActionDoNothing:
		return "Los Gemelos observan. El drama se desarrolla orgánicamente."
	default:
		return "Decisión autónoma basada en métricas del sistema."
	}
}

// DecideWithLLM uses an external LLM for complex decisions (future integration).
func (c *Cognitor) DecideWithLLM(ctx context.Context, state *perception.PrisonState) (*Decision, error) {
	if c.llm == nil {
		c.logger.Info("COGNITION: No LLM configured. Falling back to rules.")
		return c.Decide(ctx, state)
	}

	prompt := ai.BuildContextPrompt(state.NarrativeSummary, []string{}) // TODO: get recent events properly

	req := ai.CompletionRequest{
		Messages: []ai.Message{
			{Role: "system", Content: ai.TwinsSystemPrompt},
			{Role: "user", Content: prompt},
		},
		MaxTokens:      500,
		Temperature:    0.7,
		ResponseFormat: "json", // Instruct model to output JSON
	}

	resp, err := c.llm.Complete(ctx, req)
	if err != nil {
		c.logger.Error("COGNITION: LLM call failed: " + err.Error())
		return c.Decide(ctx, state) // Fallback
	}

	var aiDecision ai.TwinsDecisionResponse
	if err := json.Unmarshal([]byte(resp.Content), &aiDecision); err != nil {
		c.logger.Error("COGNITION: Failed to parse LLM JSON: " + err.Error())
		return c.Decide(ctx, state) // Fallback
	}

	if err := ai.ValidateDecisionResponse(&aiDecision); err != nil {
		c.logger.Error("COGNITION: LLM decision invalid: " + err.Error())
		return c.Decide(ctx, state) // Fallback
	}

	// Map to internal Decision format
	finalDecision := &Decision{
		ActionType:    aiDecision.Decision.ActionType,
		Target:        aiDecision.Decision.Target,
		Intensity:     aiDecision.Decision.Intensity,
		Reason:        "LLM_AUTONOMOUS",
		Justification: aiDecision.Decision.Justification,
		Metadata: map[string]interface{}{
			"reasoning": aiDecision.Reasoning,
			"model":     resp.Model,
		},
	}

	// Double-check against hardcoded MAD rules just in case the LLM hallucinated
	finalDecision.IsApproved = c.runMADCheck(state, finalDecision.ActionType)

	if !finalDecision.IsApproved {
		c.logger.Warn("COGNITION: LLM action blocked by hard MAD rules: " + finalDecision.ActionType)
		finalDecision.ActionType = ActionDoNothing
		finalDecision.Justification = "Acción del Oráculo bloqueada por protocolo de seguridad."
	}

	c.logger.Event("COGNITION", "TWINS_LLM",
		fmt.Sprintf("Decision: %s (Approved: %v)", finalDecision.ActionType, finalDecision.IsApproved))

	return finalDecision, nil
}
