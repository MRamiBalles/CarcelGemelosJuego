// Package main - test_runner.go
// Executable to run Shadow Mode stress tests.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/test"
)

func main() {
	fmt.Println("🦅 CÁRCEL DE LOS GEMELOS - SHADOW MODE TEST SUITE")
	fmt.Println("================================================")

	ctx := context.Background()

	// Test 1: Day 1 Riot
	fmt.Println("\n🧪 Iniciando Test: El Motín del Día 1...")
	riotTest := test.NewDay1RiotTest()
	riotTest.RunTest(ctx)

	// Summary
	results := riotTest.GetResults()
	passed := 0
	failed := 0

	for _, r := range results {
		if r.Passed {
			passed++
		} else {
			failed++
		}
	}

	fmt.Println("\n" + string(repeatChar('=', 60)))
	fmt.Println("📊 RESUMEN DE PRUEBAS")
	fmt.Println(string(repeatChar('=', 60)))
	fmt.Printf("   ✅ Pasadas: %d\n", passed)
	fmt.Printf("   ❌ Fallidas: %d\n", failed)

	if failed > 0 {
		fmt.Println("\n⚠️  Los Gemelos requieren recalibración")
		os.Exit(1)
	} else {
		fmt.Println("\n✅ Los Gemelos están listos para el despliegue")
		os.Exit(0)
	}
}

func repeatChar(c byte, count int) []byte {
	result := make([]byte, count)
	for i := 0; i < count; i++ {
		result[i] = c
	}
	return result
}
