// Package main - agitator
// T029: "El Agitador" - Load generator for stress testing
// Simulates 50+ concurrent prisoners spamming WebSocket actions
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// Config for the agitator
type Config struct {
	ServerURL      string
	NumClients     int
	ActionInterval time.Duration
	TestDuration   time.Duration
	GameID         string
}

// Stats tracks performance metrics
type Stats struct {
	MessagesSent     int64
	MessagesReceived int64
	Errors           int64
	Latencies        []time.Duration
	mu               sync.Mutex
}

// Action types for simulation
var actionTypes = []string{
	"MOVE",
	"SPEAK",
	"VOTE_TORTURE",
	"VOTE_MERCY",
	"WHISPER",
	"BETRAY",
	"COLLABORATE",
	"REQUEST_RESOURCE",
}

func main() {
	// Parse flags
	serverURL := flag.String("url", "ws://localhost:8080/ws", "WebSocket server URL")
	numClients := flag.Int("clients", 50, "Number of concurrent clients")
	interval := flag.Duration("interval", 100*time.Millisecond, "Action interval per client")
	duration := flag.Duration("duration", 60*time.Second, "Test duration")
	gameID := flag.String("game", "STRESS_TEST_001", "Game ID")
	flag.Parse()

	config := Config{
		ServerURL:      *serverURL,
		NumClients:     *numClients,
		ActionInterval: *interval,
		TestDuration:   *duration,
		GameID:         *gameID,
	}

	fmt.Println("=========================================")
	fmt.Println("🔥 EL AGITADOR - Stress Test Tool")
	fmt.Println("=========================================")
	fmt.Printf("Server: %s\n", config.ServerURL)
	fmt.Printf("Clients: %d\n", config.NumClients)
	fmt.Printf("Interval: %v\n", config.ActionInterval)
	fmt.Printf("Duration: %v\n", config.TestDuration)
	fmt.Println("=========================================")

	// Setup graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), config.TestDuration)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		fmt.Println("\n⚠️ Interrupt received, stopping...")
		cancel()
	}()

	// Run the stress test
	stats := runStressTest(ctx, config)

	// Print results
	printResults(stats, config)
}

func runStressTest(ctx context.Context, config Config) *Stats {
	stats := &Stats{
		Latencies: make([]time.Duration, 0, 10000),
	}

	var wg sync.WaitGroup

	fmt.Println("\n🚀 Starting clients...")

	for i := 0; i < config.NumClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			runClient(ctx, clientID, config, stats)
		}(i)

		// Stagger client starts to avoid thundering herd
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Printf("✅ All %d clients started\n\n", config.NumClients)

	// Progress updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				sent := atomic.LoadInt64(&stats.MessagesSent)
				recv := atomic.LoadInt64(&stats.MessagesReceived)
				errs := atomic.LoadInt64(&stats.Errors)
				fmt.Printf("📊 Progress: Sent=%d Recv=%d Errors=%d\n", sent, recv, errs)
			}
		}
	}()

	wg.Wait()
	return stats
}

func runClient(ctx context.Context, clientID int, config Config, stats *Stats) {
	prisonerID := fmt.Sprintf("PRISONER_%03d", clientID)

	// Parse URL and add query params
	u, err := url.Parse(config.ServerURL)
	if err != nil {
		log.Printf("Client %d: URL parse error: %v", clientID, err)
		atomic.AddInt64(&stats.Errors, 1)
		return
	}

	q := u.Query()
	q.Set("prisoner_id", prisonerID)
	q.Set("game_id", config.GameID)
	u.RawQuery = q.Encode()

	// Connect
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		log.Printf("Client %d: Connection failed: %v", clientID, err)
		atomic.AddInt64(&stats.Errors, 1)
		return
	}
	defer conn.Close()

	// Start receiver goroutine
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			atomic.AddInt64(&stats.MessagesReceived, 1)
		}
	}()

	// Send actions at configured interval
	ticker := time.NewTicker(config.ActionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			action := generateRandomAction(prisonerID, config.GameID)
			start := time.Now()

			if err := conn.WriteJSON(action); err != nil {
				atomic.AddInt64(&stats.Errors, 1)
				return
			}

			latency := time.Since(start)
			atomic.AddInt64(&stats.MessagesSent, 1)

			stats.mu.Lock()
			stats.Latencies = append(stats.Latencies, latency)
			stats.mu.Unlock()
		}
	}
}

func generateRandomAction(prisonerID, gameID string) map[string]interface{} {
	actionType := actionTypes[rand.Intn(len(actionTypes))]

	action := map[string]interface{}{
		"type":        actionType,
		"prisoner_id": prisonerID,
		"game_id":     gameID,
		"timestamp":   time.Now().UnixMilli(),
	}

	// Add type-specific payloads
	switch actionType {
	case "SPEAK":
		messages := []string{
			"Los Gemelos nos escuchan...",
			"¿Quién votó por la tortura?",
			"Necesitamos un plan de escape.",
			"No confío en nadie aquí.",
			"El ruido me está volviendo loco.",
		}
		action["message"] = messages[rand.Intn(len(messages))]

	case "VOTE_TORTURE", "VOTE_MERCY":
		targets := []string{"PRISONER_001", "PRISONER_002", "PRISONER_003"}
		action["target"] = targets[rand.Intn(len(targets))]

	case "WHISPER":
		action["target"] = fmt.Sprintf("PRISONER_%03d", rand.Intn(50))
		action["message"] = "Mensaje secreto..."

	case "BETRAY":
		action["target"] = fmt.Sprintf("PRISONER_%03d", rand.Intn(50))
		action["evidence"] = "Vi algo sospechoso"

	case "MOVE":
		locations := []string{"CELL", "COURTYARD", "CAFETERIA", "CORRIDOR"}
		action["destination"] = locations[rand.Intn(len(locations))]
	}

	return action
}

func printResults(stats *Stats, config Config) {
	fmt.Println("\n=========================================")
	fmt.Println("📊 STRESS TEST RESULTS")
	fmt.Println("=========================================")

	sent := atomic.LoadInt64(&stats.MessagesSent)
	recv := atomic.LoadInt64(&stats.MessagesReceived)
	errs := atomic.LoadInt64(&stats.Errors)

	fmt.Printf("Messages Sent:     %d\n", sent)
	fmt.Printf("Messages Received: %d\n", recv)
	fmt.Printf("Errors:            %d\n", errs)
	fmt.Printf("Error Rate:        %.2f%%\n", float64(errs)/float64(sent+1)*100)

	// Calculate throughput
	throughput := float64(sent) / config.TestDuration.Seconds()
	fmt.Printf("Throughput:        %.2f msg/sec\n", throughput)

	// Latency stats
	if len(stats.Latencies) > 0 {
		var total time.Duration
		var min, max time.Duration = stats.Latencies[0], stats.Latencies[0]

		for _, l := range stats.Latencies {
			total += l
			if l < min {
				min = l
			}
			if l > max {
				max = l
			}
		}

		avg := total / time.Duration(len(stats.Latencies))

		fmt.Printf("\nLatency:\n")
		fmt.Printf("  Min: %v\n", min)
		fmt.Printf("  Avg: %v\n", avg)
		fmt.Printf("  Max: %v\n", max)
	}

	// Verdict
	fmt.Println("\n-----------------------------------------")
	if errs == 0 && float64(sent) > float64(config.NumClients)*config.TestDuration.Seconds()*5 {
		fmt.Println("✅ TEST PASSED: System handled the load")
	} else if float64(errs)/float64(sent+1) < 0.05 {
		fmt.Println("⚠️ TEST WARNING: Some errors detected")
	} else {
		fmt.Println("❌ TEST FAILED: High error rate")
	}
	fmt.Println("=========================================")

	// Export results as JSON
	results := map[string]interface{}{
		"messages_sent":     sent,
		"messages_received": recv,
		"errors":            errs,
		"throughput_per_sec": throughput,
		"config": map[string]interface{}{
			"clients":  config.NumClients,
			"interval": config.ActionInterval.String(),
			"duration": config.TestDuration.String(),
		},
	}

	jsonData, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("stress_test_results.json", jsonData, 0644)
	fmt.Println("\n📁 Results saved to stress_test_results.json")
}
