// Package cell defines the domain entity for a prison cell.
// This package is PURE and must NOT import any infrastructure packages.
package cell

// Cell represents a physical room where a duo of prisoners is assigned.
type Cell struct {
	ID        string    `json:"id"`
	Occupants [2]string `json:"occupants"` // Max 2 prisoner IDs
	IsLocked  bool      `json:"is_locked"`
}

// NewCell creates a new cell with empty occupants.
func NewCell(id string) *Cell {
	return &Cell{
		ID:        id,
		Occupants: [2]string{"", ""},
		IsLocked:  false,
	}
}

// AddOccupant attempts to add a prisoner ID to the cell.
// Returns false if the cell is full.
func (c *Cell) AddOccupant(prisonerID string) bool {
	if c.Occupants[0] == "" {
		c.Occupants[0] = prisonerID
		return true
	}
	if c.Occupants[1] == "" {
		c.Occupants[1] = prisonerID
		return true
	}
	return false
}

// RemoveOccupant removes a prisoner from the cell.
func (c *Cell) RemoveOccupant(prisonerID string) {
	if c.Occupants[0] == prisonerID {
		c.Occupants[0] = ""
	} else if c.Occupants[1] == prisonerID {
		c.Occupants[1] = ""
	}
}
