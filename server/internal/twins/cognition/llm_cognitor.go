// Package cognition - llm_cognitor.go
// T020/T021: LLM-powered Cognitor with Shadow Mode.
// This replaces the rule-based Cognitor when LLM is available.
package cognition

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/infra/ai"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/perception"
)

// LLMCognitor uses an LLM for decision-making instead of rules.
type LLMCognitor struct {
	provider    ai.LLMProvider
	logger      *logger.Logger
	shadowMode  bool
	fallback    *Cognitor // Rule-based fallback
	maxRetries  int
}

// NewLLMCognitor creates an LLM-powered cognitor.
func NewLLMCognitor(provider ai.LLMProvider, log *logger.Logger, fallback *Cognitor) *LLMCognitor {
	return &LLMCognitor{
		provider:   provider,
		logger:     log,
		shadowMode: true, // Default to shadow mode for safety
		fallback:   fallback,
		maxRetries: 2,
	}
}

// SetShadowMode enables/disables shadow mode.
func (lc *LLMCognitor) SetShadowMode(enabled bool) {
	lc.shadowMode = enabled
	if enabled {
		lc.logger.Info("LLM Cognitor: Shadow Mode ENABLED (decisions not executed)")
	} else {
		lc.logger.Warn("LLM Cognitor: Shadow Mode DISABLED (decisions WILL execute)")
	}
}

// IsShadowMode returns the current shadow mode state.
func (lc *LLMCognitor) IsShadowMode() bool {
	return lc.shadowMode
}

// Decide uses the LLM to make a decision, with MAD validation.
func (lc *LLMCognitor) Decide(ctx context.Context, state *perception.PrisonState) (*Decision, error) {
	if !lc.provider.IsAvailable() {
		lc.logger.Warn("LLM provider unavailable, falling back to rule-based")
		return lc.fallback.Decide(ctx, state)
	}

	// Build the prompt
	contextPrompt := ai.BuildContextPrompt(state.NarrativeSummary, lc.extractRecentEvents(state))

	req := ai.CompletionRequest{
		Messages: []ai.Message{
			{Role: "system", Content: ai.TwinsSystemPrompt},
			{Role: "user", Content: contextPrompt},
		},
		MaxTokens:   1000,
		Temperature: 0.7, // Some creativity for drama
	}

	// Try LLM with retries
	var llmResp *ai.CompletionResponse
	var err error
	
	for attempt := 0; attempt <= lc.maxRetries; attempt++ {
		llmResp, err = lc.provider.Complete(ctx, req)
		if err == nil {
			break
		}
		lc.logger.Warn(fmt.Sprintf("LLM attempt %d failed: %v", attempt+1, err))
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}

	if err != nil {
		lc.logger.Error("LLM failed after retries, using fallback: " + err.Error())
		return lc.fallback.Decide(ctx, state)
	}

	// Parse the JSON response
	var twinsResp ai.TwinsDecisionResponse
	if err := json.Unmarshal([]byte(llmResp.Content), &twinsResp); err != nil {
		lc.logger.Error("Failed to parse LLM response: " + err.Error())
		lc.logger.Event("LLM_RAW", "TWINS", llmResp.Content)
		return lc.fallback.Decide(ctx, state)
	}

	// Validate the response
	if err := ai.ValidateDecisionResponse(&twinsResp); err != nil {
		lc.logger.Error("LLM response validation failed: " + err.Error())
		return lc.fallback.Decide(ctx, state)
	}

	// Log the full reasoning for audit
	lc.logger.Event("LLM_REASONING", "TWINS", twinsResp.Reasoning)

	// Convert to internal Decision format
	decision := &Decision{
		ActionType:    twinsResp.Decision.ActionType,
		Target:        twinsResp.Decision.Target,
		Intensity:     twinsResp.Decision.Intensity,
		Reason:        "LLM_AUTONOMOUS",
		Justification: twinsResp.Decision.Justification,
		IsApproved:    twinsResp.MADCheck.Passed,
		Metadata: map[string]interface{}{
			"llm_reasoning":  twinsResp.Reasoning,
			"shadow_mode":    lc.shadowMode,
			"tokens_used":    llmResp.TotalTokens,
			"latency_ms":     llmResp.Latency.Milliseconds(),
			"provider":       lc.provider.Name(),
		},
	}

	// If in shadow mode, mark as not approved for execution
	if lc.shadowMode {
		decision.IsApproved = false
		decision.Metadata["shadow_blocked"] = true
		lc.logger.Info("Shadow Mode: Decision recorded but NOT executed")
	}

	// Additional MAD validation on our side (LLM might hallucinate)
	if !lc.validateMAD(state, decision) {
		lc.logger.Warn("Server-side MAD check failed, overriding LLM decision")
		decision.IsApproved = false
		decision.ActionType = ActionDoNothing
	}

	lc.logUsageStats()

	return decision, nil
}

// validateMAD performs server-side MAD validation (never trust the LLM completely).
func (lc *LLMCognitor) validateMAD(state *perception.PrisonState, decision *Decision) bool {
	// Rule 1: No killing low sanity prisoners
	if decision.ActionType == ActionNoise && state.AverageSanity < 10 {
		lc.logger.Warn("MAD: Blocked noise on low sanity population")
		return false
	}

	// Rule 2: No day one cruelty
	if state.CurrentDay == 1 && (decision.ActionType == ActionNoise || decision.ActionType == ActionResourceCut) {
		lc.logger.Warn("MAD: Blocked cruelty on Day 1")
		return false
	}

	// Rule 3: Require justification
	if decision.Justification == "" {
		lc.logger.Warn("MAD: No justification provided")
		return false
	}

	return true
}

// extractRecentEvents converts state data to event strings for prompt.
func (lc *LLMCognitor) extractRecentEvents(state *perception.PrisonState) []string {
	events := []string{}
	
	if state.RecentBetrayals > 0 {
		events = append(events, fmt.Sprintf("Traiciones recientes: %d", state.RecentBetrayals))
	}
	if state.AudienceActivity > 0 {
		events = append(events, fmt.Sprintf("Intervenciones de audiencia: %d", state.AudienceActivity))
	}
	events = append(events, fmt.Sprintf("Nivel de tensión: %s", state.TensionLevel))
	events = append(events, fmt.Sprintf("Cordura promedio: %.0f%%", state.AverageSanity))
	
	return events
}

// logUsageStats logs current API usage for monitoring.
func (lc *LLMCognitor) logUsageStats() {
	stats := lc.provider.GetUsageStats()
	lc.logger.Info(fmt.Sprintf("LLM Usage: %d requests, %d tokens, $%.4f spent", 
		stats.TotalRequests, stats.TotalTokens, stats.TotalCostUSD))
}
