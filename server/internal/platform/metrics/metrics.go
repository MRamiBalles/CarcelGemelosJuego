// Package metrics provides observability for the game server.
// T030: Metrics collection for stress testing analysis.
package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Collector gathers performance metrics.
type Collector struct {
	// Tick metrics
	TickCount      int64
	TickLatencySum int64 // nanoseconds
	TickLatencyMax int64
	LastTickTime   time.Time

	// Event metrics
	EventsWritten    int64
	EventWriteLatSum int64
	EventWriteLatMax int64
	EventWriteErrors int64

	// WebSocket metrics
	WSConnectionsActive int64
	WSMessagesIn        int64
	WSMessagesOut       int64
	WSErrors            int64

	// LLM metrics
	LLMRequests   int64
	LLMTokensUsed int64
	LLMCostUSD    float64
	LLMLatencySum int64

	// System
	StartTime time.Time
	mu        sync.RWMutex
}

// Global collector instance
var collector = &Collector{
	StartTime: time.Now(),
}

// Get returns the global collector.
func Get() *Collector {
	return collector
}

// RecordTick records a tick cycle completion.
func (c *Collector) RecordTick(latency time.Duration) {
	atomic.AddInt64(&c.TickCount, 1)
	atomic.AddInt64(&c.TickLatencySum, int64(latency))

	// Update max (non-atomic but acceptable for metrics)
	if int64(latency) > atomic.LoadInt64(&c.TickLatencyMax) {
		atomic.StoreInt64(&c.TickLatencyMax, int64(latency))
	}

	c.mu.Lock()
	c.LastTickTime = time.Now()
	c.mu.Unlock()
}

// RecordEventWrite records an event write to the database.
func (c *Collector) RecordEventWrite(latency time.Duration, err error) {
	atomic.AddInt64(&c.EventsWritten, 1)
	atomic.AddInt64(&c.EventWriteLatSum, int64(latency))

	if int64(latency) > atomic.LoadInt64(&c.EventWriteLatMax) {
		atomic.StoreInt64(&c.EventWriteLatMax, int64(latency))
	}

	if err != nil {
		atomic.AddInt64(&c.EventWriteErrors, 1)
	}
}

// RecordWSConnection records WebSocket connection changes.
func (c *Collector) RecordWSConnection(delta int64) {
	atomic.AddInt64(&c.WSConnectionsActive, delta)
}

// RecordWSMessage records WebSocket messages.
func (c *Collector) RecordWSMessage(incoming bool) {
	if incoming {
		atomic.AddInt64(&c.WSMessagesIn, 1)
	} else {
		atomic.AddInt64(&c.WSMessagesOut, 1)
	}
}

// RecordWSError records a WebSocket error.
func (c *Collector) RecordWSError() {
	atomic.AddInt64(&c.WSErrors, 1)
}

// RecordLLMCall records an LLM API call.
func (c *Collector) RecordLLMCall(tokens int, cost float64, latency time.Duration) {
	atomic.AddInt64(&c.LLMRequests, 1)
	atomic.AddInt64(&c.LLMTokensUsed, int64(tokens))
	atomic.AddInt64(&c.LLMLatencySum, int64(latency))

	c.mu.Lock()
	c.LLMCostUSD += cost
	c.mu.Unlock()
}

// Snapshot returns current metrics as a map.
func (c *Collector) Snapshot() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tickCount := atomic.LoadInt64(&c.TickCount)
	eventsWritten := atomic.LoadInt64(&c.EventsWritten)
	llmRequests := atomic.LoadInt64(&c.LLMRequests)

	// Calculate averages
	var tickAvg, eventAvg, llmAvg float64
	if tickCount > 0 {
		tickAvg = float64(atomic.LoadInt64(&c.TickLatencySum)) / float64(tickCount) / 1e6 // ms
	}
	if eventsWritten > 0 {
		eventAvg = float64(atomic.LoadInt64(&c.EventWriteLatSum)) / float64(eventsWritten) / 1e6
	}
	if llmRequests > 0 {
		llmAvg = float64(atomic.LoadInt64(&c.LLMLatencySum)) / float64(llmRequests) / 1e9 // seconds
	}

	return map[string]interface{}{
		"uptime_seconds": time.Since(c.StartTime).Seconds(),

		"tick": map[string]interface{}{
			"count":          tickCount,
			"avg_latency_ms": tickAvg,
			"max_latency_ms": float64(atomic.LoadInt64(&c.TickLatencyMax)) / 1e6,
			"last_tick":      c.LastTickTime.Format(time.RFC3339),
		},

		"events": map[string]interface{}{
			"written":          eventsWritten,
			"avg_write_lat_ms": eventAvg,
			"max_write_lat_ms": float64(atomic.LoadInt64(&c.EventWriteLatMax)) / 1e6,
			"errors":           atomic.LoadInt64(&c.EventWriteErrors),
		},

		"websocket": map[string]interface{}{
			"active_connections": atomic.LoadInt64(&c.WSConnectionsActive),
			"messages_in":        atomic.LoadInt64(&c.WSMessagesIn),
			"messages_out":       atomic.LoadInt64(&c.WSMessagesOut),
			"errors":             atomic.LoadInt64(&c.WSErrors),
		},

		"llm": map[string]interface{}{
			"requests":        llmRequests,
			"tokens_used":     atomic.LoadInt64(&c.LLMTokensUsed),
			"cost_usd":        c.LLMCostUSD,
			"avg_latency_sec": llmAvg,
		},
	}
}

// Handler returns an HTTP handler for the /metrics endpoint.
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")

		snapshot := collector.Snapshot()
		json.NewEncoder(w).Encode(snapshot)
	}
}

// PrometheusHandler returns metrics in Prometheus format.
func PrometheusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		c := collector

		// Tick metrics
		fmt.Fprintf(w, "# HELP carcel_tick_count Total tick cycles\n")
		fmt.Fprintf(w, "# TYPE carcel_tick_count counter\n")
		fmt.Fprintf(w, "carcel_tick_count %d\n\n", atomic.LoadInt64(&c.TickCount))

		fmt.Fprintf(w, "# HELP carcel_tick_latency_max_ms Maximum tick latency\n")
		fmt.Fprintf(w, "# TYPE carcel_tick_latency_max_ms gauge\n")
		fmt.Fprintf(w, "carcel_tick_latency_max_ms %.2f\n\n", float64(atomic.LoadInt64(&c.TickLatencyMax))/1e6)

		// Event metrics
		fmt.Fprintf(w, "# HELP carcel_events_written Total events written\n")
		fmt.Fprintf(w, "# TYPE carcel_events_written counter\n")
		fmt.Fprintf(w, "carcel_events_written %d\n\n", atomic.LoadInt64(&c.EventsWritten))

		fmt.Fprintf(w, "# HELP carcel_event_write_errors Total event write errors\n")
		fmt.Fprintf(w, "# TYPE carcel_event_write_errors counter\n")
		fmt.Fprintf(w, "carcel_event_write_errors %d\n\n", atomic.LoadInt64(&c.EventWriteErrors))

		// WebSocket metrics
		fmt.Fprintf(w, "# HELP carcel_ws_connections Active WebSocket connections\n")
		fmt.Fprintf(w, "# TYPE carcel_ws_connections gauge\n")
		fmt.Fprintf(w, "carcel_ws_connections %d\n\n", atomic.LoadInt64(&c.WSConnectionsActive))

		fmt.Fprintf(w, "# HELP carcel_ws_messages_total Total WebSocket messages\n")
		fmt.Fprintf(w, "# TYPE carcel_ws_messages_total counter\n")
		fmt.Fprintf(w, "carcel_ws_messages_total{direction=\"in\"} %d\n", atomic.LoadInt64(&c.WSMessagesIn))
		fmt.Fprintf(w, "carcel_ws_messages_total{direction=\"out\"} %d\n\n", atomic.LoadInt64(&c.WSMessagesOut))

		// LLM metrics
		fmt.Fprintf(w, "# HELP carcel_llm_requests Total LLM API requests\n")
		fmt.Fprintf(w, "# TYPE carcel_llm_requests counter\n")
		fmt.Fprintf(w, "carcel_llm_requests %d\n\n", atomic.LoadInt64(&c.LLMRequests))

		fmt.Fprintf(w, "# HELP carcel_llm_tokens_used Total tokens consumed\n")
		fmt.Fprintf(w, "# TYPE carcel_llm_tokens_used counter\n")
		fmt.Fprintf(w, "carcel_llm_tokens_used %d\n\n", atomic.LoadInt64(&c.LLMTokensUsed))

		c.mu.RLock()
		fmt.Fprintf(w, "# HELP carcel_llm_cost_usd Total LLM cost in USD\n")
		fmt.Fprintf(w, "# TYPE carcel_llm_cost_usd counter\n")
		fmt.Fprintf(w, "carcel_llm_cost_usd %.4f\n", c.LLMCostUSD)
		c.mu.RUnlock()
	}
}
