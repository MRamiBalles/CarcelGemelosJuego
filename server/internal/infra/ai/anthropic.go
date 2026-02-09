// Package ai - anthropic.go
// T019: Anthropic Claude adapter implementing the LLMProvider interface.
// Claude 3.5 Sonnet recommended for superior reasoning and larger context.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// AnthropicProvider implements LLMProvider for Anthropic Claude API.
type AnthropicProvider struct {
	apiKey      string
	baseURL     string
	model       string
	httpClient  *http.Client
	usageStats  UsageStats
	budgetGate  *BudgetGate
}

// Anthropic API structures
type anthropicRequest struct {
	Model       string               `json:"model"`
	MaxTokens   int                  `json:"max_tokens"`
	System      string               `json:"system,omitempty"`
	Messages    []anthropicMessage   `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model      string `json:"model"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// NewAnthropicProvider creates a new Claude adapter.
func NewAnthropicProvider(budgetGate *BudgetGate) *AnthropicProvider {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	
	return &AnthropicProvider{
		apiKey:     apiKey,
		baseURL:    "https://api.anthropic.com/v1/messages",
		model:      "claude-3-5-sonnet-20241022", // Best reasoning
		httpClient: &http.Client{Timeout: 120 * time.Second},
		budgetGate: budgetGate,
	}
}

// Name returns the provider name.
func (p *AnthropicProvider) Name() string {
	return "Anthropic Claude"
}

// IsAvailable checks if the API key is configured.
func (p *AnthropicProvider) IsAvailable() bool {
	return p.apiKey != ""
}

// Complete sends a completion request to Claude.
func (p *AnthropicProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	if !p.IsAvailable() {
		return nil, fmt.Errorf("Anthropic API key not configured")
	}

	// Check budget
	estimatedCost := p.estimateCost(req)
	if !p.budgetGate.CanSpend(estimatedCost) {
		return nil, fmt.Errorf("budget limit exceeded: %s", p.budgetGate.GetStatus())
	}

	// Extract system message and build user messages
	var systemMsg string
	var messages []anthropicMessage
	
	for _, m := range req.Messages {
		if m.Role == "system" {
			systemMsg = m.Content
		} else {
			messages = append(messages, anthropicMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}
	}

	model := p.model
	if req.Model != "" {
		model = req.Model
	}

	anthReq := anthropicRequest{
		Model:     model,
		MaxTokens: req.MaxTokens,
		System:    systemMsg,
		Messages:  messages,
	}

	body, err := json.Marshal(anthReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send request
	start := time.Now()
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Anthropic error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var anthResp anthropicResponse
	if err := json.Unmarshal(respBody, &anthResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(anthResp.Content) == 0 {
		return nil, fmt.Errorf("no response content returned")
	}

	totalTokens := anthResp.Usage.InputTokens + anthResp.Usage.OutputTokens
	actualCost := p.calculateCost(totalTokens, model)
	p.budgetGate.RecordSpend(actualCost)
	p.usageStats.TotalRequests++
	p.usageStats.TotalTokens += totalTokens
	p.usageStats.TotalCostUSD += actualCost

	return &CompletionResponse{
		Content:      anthResp.Content[0].Text,
		Model:        anthResp.Model,
		PromptTokens: anthResp.Usage.InputTokens,
		OutputTokens: anthResp.Usage.OutputTokens,
		TotalTokens:  totalTokens,
		Latency:      latency,
		FinishReason: anthResp.StopReason,
	}, nil
}

// estimateCost estimates cost before making a request.
func (p *AnthropicProvider) estimateCost(req CompletionRequest) float64 {
	estimatedTokens := 2000 + req.MaxTokens
	return p.calculateCost(estimatedTokens, p.model)
}

// calculateCost computes actual cost based on tokens.
func (p *AnthropicProvider) calculateCost(tokens int, model string) float64 {
	// Claude 3.5 Sonnet: ~$3/1M input, ~$15/1M output
	// Averaged: ~$0.009 per 1K tokens
	switch model {
	case "claude-3-5-sonnet-20241022":
		return float64(tokens) * 0.000009
	case "claude-3-haiku-20240307":
		return float64(tokens) * 0.0000005 // Much cheaper
	default:
		return float64(tokens) * 0.00001
	}
}

// GetUsageStats returns current usage statistics.
func (p *AnthropicProvider) GetUsageStats() UsageStats {
	p.usageStats.BudgetRemaining = p.budgetGate.MonthlyLimitUSD - p.budgetGate.CurrentMonthSpend
	return p.usageStats
}

// ResetUsage resets all usage counters.
func (p *AnthropicProvider) ResetUsage() {
	p.usageStats = UsageStats{LastReset: time.Now()}
}

// Ensure AnthropicProvider implements LLMProvider
var _ LLMProvider = (*AnthropicProvider)(nil)
