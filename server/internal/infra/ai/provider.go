// Package ai provides the LLM integration layer for Los Gemelos.
// T019: Agnostic LLM Provider interface that allows swapping between
// OpenAI, Anthropic Claude, DeepSeek, or local models.
package ai

import (
	"context"
	"time"
)

// Message represents a chat message for the LLM.
type Message struct {
	Role    string `json:"role"`    // "system", "user", "assistant"
	Content string `json:"content"`
}

// CompletionRequest is the input for LLM inference.
type CompletionRequest struct {
	Messages      []Message `json:"messages"`
	MaxTokens     int       `json:"max_tokens"`
	Temperature   float64   `json:"temperature"`
	Model         string    `json:"model,omitempty"`         // Override default model
	ResponseFormat string   `json:"response_format,omitempty"` // "json" for structured output
}

// CompletionResponse is the output from LLM inference.
type CompletionResponse struct {
	Content      string        `json:"content"`
	Model        string        `json:"model"`
	PromptTokens int           `json:"prompt_tokens"`
	OutputTokens int           `json:"output_tokens"`
	TotalTokens  int           `json:"total_tokens"`
	Latency      time.Duration `json:"latency"`
	FinishReason string        `json:"finish_reason"`
}

// UsageStats tracks API usage for FinOps.
type UsageStats struct {
	TotalRequests   int     `json:"total_requests"`
	TotalTokens     int     `json:"total_tokens"`
	TotalCostUSD    float64 `json:"total_cost_usd"`
	BudgetRemaining float64 `json:"budget_remaining"`
	LastReset       time.Time `json:"last_reset"`
}

// LLMProvider is the agnostic interface for LLM backends.
// The Cognitor uses this interface without knowing which provider is behind it.
type LLMProvider interface {
	// Complete sends a prompt and returns the LLM response.
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

	// GetUsageStats returns current API usage for FinOps monitoring.
	GetUsageStats() UsageStats

	// ResetUsage resets the usage counters (e.g., monthly reset).
	ResetUsage()

	// Name returns the provider name (for logging).
	Name() string

	// IsAvailable checks if the provider is configured and reachable.
	IsAvailable() bool
}

// BudgetGate controls spending limits for LLM calls.
type BudgetGate struct {
	DailyLimitUSD   float64
	MonthlyLimitUSD float64
	CurrentDaySpend float64
	CurrentMonthSpend float64
	LastDayReset    time.Time
	LastMonthReset  time.Time
}

// NewBudgetGate creates a new budget controller.
func NewBudgetGate(dailyLimit, monthlyLimit float64) *BudgetGate {
	now := time.Now()
	return &BudgetGate{
		DailyLimitUSD:     dailyLimit,
		MonthlyLimitUSD:   monthlyLimit,
		CurrentDaySpend:   0,
		CurrentMonthSpend: 0,
		LastDayReset:      now,
		LastMonthReset:    now,
	}
}

// CanSpend checks if a cost is within budget.
func (bg *BudgetGate) CanSpend(costUSD float64) bool {
	bg.maybeReset()
	return (bg.CurrentDaySpend+costUSD <= bg.DailyLimitUSD) &&
		(bg.CurrentMonthSpend+costUSD <= bg.MonthlyLimitUSD)
}

// RecordSpend logs a cost.
func (bg *BudgetGate) RecordSpend(costUSD float64) {
	bg.maybeReset()
	bg.CurrentDaySpend += costUSD
	bg.CurrentMonthSpend += costUSD
}

// maybeReset resets counters if day/month has changed.
func (bg *BudgetGate) maybeReset() {
	now := time.Now()
	
	// Daily reset
	if now.YearDay() != bg.LastDayReset.YearDay() || now.Year() != bg.LastDayReset.Year() {
		bg.CurrentDaySpend = 0
		bg.LastDayReset = now
	}
	
	// Monthly reset
	if now.Month() != bg.LastMonthReset.Month() || now.Year() != bg.LastMonthReset.Year() {
		bg.CurrentMonthSpend = 0
		bg.LastMonthReset = now
	}
}

// GetStatus returns a human-readable budget status.
func (bg *BudgetGate) GetStatus() string {
	return "Day: $" + formatFloat(bg.CurrentDaySpend) + "/" + formatFloat(bg.DailyLimitUSD) +
		" | Month: $" + formatFloat(bg.CurrentMonthSpend) + "/" + formatFloat(bg.MonthlyLimitUSD)
}

func formatFloat(f float64) string {
	// Simple formatting without importing fmt
	return string(rune('0'+int(f))) + "." + string(rune('0'+int(f*10)%10)) + string(rune('0'+int(f*100)%10))
}
