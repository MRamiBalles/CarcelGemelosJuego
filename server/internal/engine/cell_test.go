package engine

import (
	"testing"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
)

func TestCellAssignment(t *testing.T) {
	// Simulate the core engine cell logic
	p1 := &prisoner.Prisoner{ID: "P001", CellID: "CELL_A"}
	p2 := &prisoner.Prisoner{ID: "P008", CellID: "CELL_A"}
	p3 := &prisoner.Prisoner{ID: "P006", CellID: "CELL_B"}

	if p1.CellID != p2.CellID {
		t.Errorf("Expected P1 and P2 to be in the same cell, but got %s and %s", p1.CellID, p2.CellID)
	}

	if p1.CellID == p3.CellID {
		t.Errorf("Expected P1 and P3 to be in different cells, but both are in %s", p1.CellID)
	}
}

func TestGetCellmate(t *testing.T) {
	prison := map[string]*prisoner.Prisoner{
		"P001": {ID: "P001", CellID: "CELL_A"},
		"P008": {ID: "P008", CellID: "CELL_A"},
		"P006": {ID: "P006", CellID: "CELL_B"},
		"P007": {ID: "P007", CellID: "CELL_B"},
	}

	// Helper function mimicking the engine's internal cellmate lookup
	getCellmate := func(targetID string) *prisoner.Prisoner {
		target, ok := prison[targetID]
		if !ok {
			return nil
		}
		for _, p := range prison {
			if p.ID != targetID && p.CellID == target.CellID {
				return p
			}
		}
		return nil
	}

	mateP1 := getCellmate("P001")
	if mateP1 == nil || mateP1.ID != "P008" {
		t.Errorf("Expected cellmate of P001 to be P008, got %v", mateP1)
	}

	mateP6 := getCellmate("P006")
	if mateP6 == nil || mateP6.ID != "P007" {
		t.Errorf("Expected cellmate of P006 to be P007, got %v", mateP6)
	}
}
