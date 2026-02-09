// Package optimization provides concurrency tuning for high load.
// T031: Optimized channel buffers and connection pool settings.
package optimization

import (
	"runtime"
)

// Config holds tuned parameters for high-load scenarios.
type Config struct {
	// Channel buffer sizes
	EventChannelBuffer     int
	BroadcastChannelBuffer int
	ClientSendBuffer       int

	// Connection pools
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	RedisPoolSize     int

	// Worker pools
	EventWorkers     int
	BroadcastWorkers int

	// Rate limiting
	MaxMessagesPerSecond int
	MaxClientsPerGame    int
}

// DefaultConfig returns sensible defaults for production.
func DefaultConfig() *Config {
	numCPU := runtime.NumCPU()

	return &Config{
		// Channel buffers - larger = more memory, less blocking
		EventChannelBuffer:     1024,  // Handle bursts
		BroadcastChannelBuffer: 256,   // Per client
		ClientSendBuffer:       64,    // Per WebSocket

		// Database connections
		DBMaxOpenConns: numCPU * 4,    // 4 connections per CPU
		DBMaxIdleConns: numCPU * 2,    // Keep half warm

		// Redis
		RedisPoolSize: numCPU * 2,

		// Workers
		EventWorkers:     numCPU,       // One per CPU for event processing
		BroadcastWorkers: numCPU * 2,   // Two per CPU for I/O bound work

		// Rate limits
		MaxMessagesPerSecond: 100,      // Per client
		MaxClientsPerGame:    200,      // Per game instance
	}
}

// StressTestConfig returns aggressive settings for stress testing.
func StressTestConfig() *Config {
	numCPU := runtime.NumCPU()

	return &Config{
		EventChannelBuffer:     4096,
		BroadcastChannelBuffer: 512,
		ClientSendBuffer:       128,

		DBMaxOpenConns: numCPU * 8,
		DBMaxIdleConns: numCPU * 4,
		RedisPoolSize:  numCPU * 4,

		EventWorkers:     numCPU * 2,
		BroadcastWorkers: numCPU * 4,

		MaxMessagesPerSecond: 500,
		MaxClientsPerGame:    500,
	}
}

// LowResourceConfig returns minimal settings for development.
func LowResourceConfig() *Config {
	return &Config{
		EventChannelBuffer:     64,
		BroadcastChannelBuffer: 16,
		ClientSendBuffer:       8,

		DBMaxOpenConns: 5,
		DBMaxIdleConns: 2,
		RedisPoolSize:  5,

		EventWorkers:     2,
		BroadcastWorkers: 2,

		MaxMessagesPerSecond: 10,
		MaxClientsPerGame:    20,
	}
}

// Recommendations provides suggestions based on observed metrics.
type Recommendations struct {
	IncreaseEventBuffer     bool
	IncreaseBroadcastBuffer bool
	IncreaseDBConnections   bool
	IncreaseWorkers         bool
	Notes                   []string
}

// Analyze examines current metrics and returns optimization recommendations.
func Analyze(metrics map[string]interface{}) *Recommendations {
	rec := &Recommendations{
		Notes: make([]string, 0),
	}

	// Check tick latency
	if tick, ok := metrics["tick"].(map[string]interface{}); ok {
		if maxLat, ok := tick["max_latency_ms"].(float64); ok && maxLat > 100 {
			rec.IncreaseEventBuffer = true
			rec.IncreaseWorkers = true
			rec.Notes = append(rec.Notes, "Tick latency exceeds 100ms - increase event workers")
		}
	}

	// Check event write latency
	if events, ok := metrics["events"].(map[string]interface{}); ok {
		if maxLat, ok := events["max_write_lat_ms"].(float64); ok && maxLat > 50 {
			rec.IncreaseDBConnections = true
			rec.Notes = append(rec.Notes, "Event write latency exceeds 50ms - increase DB connections")
		}
		if errors, ok := events["errors"].(int64); ok && errors > 0 {
			rec.IncreaseDBConnections = true
			rec.Notes = append(rec.Notes, "Event write errors detected - check DB connection pool")
		}
	}

	// Check WebSocket backpressure
	if ws, ok := metrics["websocket"].(map[string]interface{}); ok {
		if errors, ok := ws["errors"].(int64); ok && errors > 0 {
			rec.IncreaseBroadcastBuffer = true
			rec.Notes = append(rec.Notes, "WebSocket errors detected - increase client send buffer")
		}
	}

	return rec
}

// ApplyRecommendations modifies config based on recommendations.
func ApplyRecommendations(config *Config, rec *Recommendations) *Config {
	if rec.IncreaseEventBuffer {
		config.EventChannelBuffer *= 2
	}
	if rec.IncreaseBroadcastBuffer {
		config.BroadcastChannelBuffer *= 2
		config.ClientSendBuffer *= 2
	}
	if rec.IncreaseDBConnections {
		config.DBMaxOpenConns = int(float64(config.DBMaxOpenConns) * 1.5)
		config.DBMaxIdleConns = int(float64(config.DBMaxIdleConns) * 1.5)
	}
	if rec.IncreaseWorkers {
		config.EventWorkers *= 2
		config.BroadcastWorkers = int(float64(config.BroadcastWorkers) * 1.5)
	}
	return config
}
