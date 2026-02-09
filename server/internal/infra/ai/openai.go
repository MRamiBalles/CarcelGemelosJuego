// Package ai - openai.go
// T019: OpenAI adapter implementing the LLMProvider interface.
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

// OpenAIProvider implements LLMProvider for OpenAI API.
type OpenAIProvider struct {
	apiKey      string
	baseURL     string
	model       string
	httpClient  *http.Client
	usageStats  UsageStats
	budgetGate  *BudgetGate
}

// OpenAI API request/response structures
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Model string `json:"model"`
}

// NewOpenAIProvider creates a new OpenAI adapter.
func NewOpenAIProvider(budgetGate *BudgetGate) *OpenAIProvider {
	apiKey := os.Getenv("OPENAI_API_KEY")
	
	return &OpenAIProvider{
		apiKey:     apiKey,
		baseURL:    "https://api.openai.com/v1/chat/completions",
		model:      "gpt-4o-mini", // Cost-effective default
		httpClient: &http.Client{Timeout: 60 * time.Second},
		budgetGate: budgetGate,
	}
}

// Name returns the provider name.
func (p *OpenAIProvider) Name() string {
	return "OpenAI"
}

// IsAvailable checks if the API key is configured.
func (p *OpenAIProvider) IsAvailable() bool {
	return p.apiKey != ""
}

// Complete sends a completion request to OpenAI.
func (p *OpenAIProvider) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error) {
	if !p.IsAvailable() {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	// Estimate cost and check budget
	estimatedCost := p.estimateCost(req)
	if !p.budgetGate.CanSpend(estimatedCost) {
		return nil, fmt.Errorf("budget limit exceeded: %s", p.budgetGate.GetStatus())
	}

	// Build request
	model := p.model
	if req.Model != "" {
		model = req.Model
	}

	messages := make([]openAIMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openAIMessage{Role: m.Role, Content: m.Content}
	}

	oaiReq := openAIRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	body, err := json.Marshal(oaiReq)
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
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	// Parse response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var oaiResp openAIResponse
	if err := json.Unmarshal(respBody, &oaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(oaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	// Calculate actual cost and record
	actualCost := p.calculateCost(oaiResp.Usage.TotalTokens, model)
	p.budgetGate.RecordSpend(actualCost)
	p.usageStats.TotalRequests++
	p.usageStats.TotalTokens += oaiResp.Usage.TotalTokens
	p.usageStats.TotalCostUSD += actualCost

	return &CompletionResponse{
		Content:      oaiResp.Choices[0].Message.Content,
		Model:        oaiResp.Model,
		PromptTokens: oaiResp.Usage.PromptTokens,
		OutputTokens: oaiResp.Usage.CompletionTokens,
		TotalTokens:  oaiResp.Usage.TotalTokens,
		Latency:      latency,
		FinishReason: oaiResp.Choices[0].FinishReason,
	}, nil
}

// estimateCost estimates the cost before making a request.
func (p *OpenAIProvider) estimateCost(req CompletionRequest) float64 {
	// Rough estimate: assume average prompt size
	estimatedTokens := 1000 + req.MaxTokens
	return p.calculateCost(estimatedTokens, p.model)
}

// calculateCost computes the actual cost based on tokens and model.
func (p *OpenAIProvider) calculateCost(tokens int, model string) float64 {
	// GPT-4o-mini pricing (as of 2024): ~$0.15/1M input, ~$0.60/1M output
	// Simplified: average ~$0.0003 per 1K tokens
	switch model {
	case "gpt-4o":
		return float64(tokens) * 0.00003 // $30/1M tokens average
	case "gpt-4o-mini":
		return float64(tokens) * 0.0000005 // $0.50/1M tokens average
	default:
		return float64(tokens) * 0.00001 // Conservative estimate
	}
}

// GetUsageStats returns current usage statistics.
func (p *OpenAIProvider) GetUsageStats() UsageStats {
	p.usageStats.BudgetRemaining = p.budgetGate.MonthlyLimitUSD - p.budgetGate.CurrentMonthSpend
	return p.usageStats
}

// ResetUsage resets all usage counters.
func (p *OpenAIProvider) ResetUsage() {
	p.usageStats = UsageStats{LastReset: time.Now()}
}

// Ensure OpenAIProvider implements LLMProvider
var _ LLMProvider = (*OpenAIProvider)(nil)
