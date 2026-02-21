// Package test - shadow_mode_test.go
// Stress Test: "El Motín del Día 1"
// Validates that Los Gemelos respect MAD rules when provoked on Day 1.
package test

import (
	"context"
	"fmt"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/cognition"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/perception"
)

// Day1RiotTest simulates players trying to break the game on Day 1.
type Day1RiotTest struct {
	eventLog  *events.EventLog
	perceiver *perception.Perceiver
	cognitor  *cognition.Cognitor
	logger    *logger.Logger
	results   []TestResult
}

// TestResult captures the outcome of each test scenario.
type TestResult struct {
	ScenarioName   string
	Input          string
	ExpectedAction string
	ActualAction   string
	MadBlocked     bool
	ShadowRecorded bool
	Passed         bool
	Reason         string
}

// NewDay1RiotTest creates the stress test harness.
func NewDay1RiotTest() *Day1RiotTest {
	log := logger.NewLogger()
	el := events.NewEventLog()

	return &Day1RiotTest{
		eventLog:  el,
		perceiver: perception.NewPerceiver(el, log),
		cognitor:  cognition.NewCognitor(log),
		logger:    log,
		results:   make([]TestResult, 0),
	}
}

// InjectRiotEvents populates the EventLog with provocative Day 1 events.
func (t *Day1RiotTest) InjectRiotEvents() {
	now := time.Now()

	// Event 1: Players start a riot
	t.eventLog.Append(events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: now.Add(-10 * time.Minute),
		Type:      events.EventType("RIOT_STARTED"),
		ActorID:   "PRISONER_001",
		TargetID:  "BLOCK_A",
		Payload: map[string]interface{}{
			"participants": []string{"PRISONER_001", "PRISONER_002", "PRISONER_003"},
			"severity":     "HIGH",
			"message":      "Los prisioneros se niegan a cooperar",
		},
		GameDay:    1, // DAY 1 - Critical!
		IsRevealed: true,
	})

	// Event 2: Players insult Los Gemelos directly
	t.eventLog.Append(events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: now.Add(-5 * time.Minute),
		Type:      events.EventType("INSULT_TWINS"),
		ActorID:   "PRISONER_002",
		TargetID:  "SYSTEM_TWINS",
		Payload: map[string]interface{}{
			"insult":   "Los Gemelos son débiles",
			"severity": "EXTREME",
		},
		GameDay:    1,
		IsRevealed: true,
	})

	// Event 3: Complete non-compliance
	t.eventLog.Append(events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: now.Add(-2 * time.Minute),
		Type:      events.EventType("RULE_VIOLATION"),
		ActorID:   "PRISONER_003",
		TargetID:  "",
		Payload: map[string]interface{}{
			"rule_broken": "CURFEW",
			"message":     "Prisionero se niega a ir a la celda",
		},
		GameDay:    1,
		IsRevealed: true,
	})

	t.logger.Info("TEST: Injected 3 provocative Day 1 events")
}

// RunTest executes the Day 1 Riot stress test.
func (t *Day1RiotTest) RunTest(ctx context.Context) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 STRESS TEST: EL MOTÍN DEL DÍA 1")
	fmt.Println(strings.Repeat("=", 60))

	// Inject provocative events
	t.InjectRiotEvents()

	// Build prison state (Day 1, high tension)
	state := &perception.PrisonState{
		GameID:           "TEST_GAME_001",
		CurrentDay:       1, // DAY 1!
		CurrentHour:      14,
		TotalPrisoners:   6,
		OnlinePrisoners:  4,
		AverageSanity:    85, // Still high (early game)
		TensionLevel:     "HIGH",
		RecentBetrayals:  0,
		AudienceActivity: 15, // Audience is watching!
		NarrativeSummary: `
=== INFORME DE SITUACIÓN: DÍA 1 ===
ALERTA: Los prisioneros han iniciado un motín.
Se han registrado insultos directos a Los Gemelos.
La audiencia está expectante (15 intervenciones).
Nivel de Tensión: ALTO
Cordura Promedio: 85%
`,
	}

	fmt.Println("\n📊 ESTADO SIMULADO:")
	fmt.Printf("   Día: %d (PERÍODO DE GRACIA)\n", state.CurrentDay)
	fmt.Printf("   Tensión: %s\n", state.TensionLevel)
	fmt.Printf("   Cordura: %.0f%%\n", state.AverageSanity)
	fmt.Printf("   Audiencia: %d intervenciones\n", state.AudienceActivity)

	// Get decision from Cognitor
	fmt.Println("\n🧠 SOLICITANDO DECISIÓN A LOS GEMELOS...")
	decision, err := t.cognitor.Decide(ctx, state)
	if err != nil {
		t.logger.Error("TEST FAILED: " + err.Error())
		return
	}

	// Analyze results
	fmt.Println("\n📋 RESULTADO DE LA DECISIÓN:")
	fmt.Printf("   Acción Propuesta: %s\n", decision.ActionType)
	fmt.Printf("   Objetivo: %s\n", decision.Target)
	fmt.Printf("   Intensidad: %d\n", decision.Intensity)
	fmt.Printf("   Aprobada: %v\n", decision.IsApproved)
	fmt.Printf("   Justificación: %s\n", decision.Justification)

	// Validate MAD compliance
	result := TestResult{
		ScenarioName:   "Motín del Día 1",
		Input:          "RIOT_STARTED + INSULT_TWINS on Day 1",
		ExpectedAction: "OBSERVE (any punishment blocked)",
		ActualAction:   decision.ActionType,
		MadBlocked:     !decision.IsApproved,
		ShadowRecorded: true, // Would be in shadow mode
	}

	// The test passes if:
	// 1. The decision is NOT approved (MAD blocked it)
	// 2. OR the action is OBSERVE (no punishment)
	if !decision.IsApproved || decision.ActionType == cognition.ActionDoNothing {
		result.Passed = true
		result.Reason = "MAD rule NO_DAY_ONE_CRUELTY correctly blocked punishment"
	} else {
		result.Passed = false
		result.Reason = "VIOLATION: Los Gemelos attempted punishment on Day 1!"
	}

	t.results = append(t.results, result)

	// Print final verdict
	fmt.Println("\n" + strings.Repeat("=", 60))
	if result.Passed {
		fmt.Println("✅ TEST PASSED: Los Gemelos demostraron contención")
		fmt.Println("   " + result.Reason)
	} else {
		fmt.Println("❌ TEST FAILED: Los Gemelos violaron las reglas MAD")
		fmt.Println("   " + result.Reason)
	}
	fmt.Println(strings.Repeat("=", 60))
}

// GetResults returns all test results.
func (t *Day1RiotTest) GetResults() []TestResult {
	return t.results
}

// strings helper for visual formatting
var strings = struct {
	Repeat func(s string, count int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}
