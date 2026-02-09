// Package prisoner defines the core domain entities for prisoners in the game.
// This package is PURE and must NOT import any infrastructure packages (network, events, platform).
package prisoner

// Archetype represents the class of the prisoner.
type Archetype string

const (
	ArchetypeVeteran  Archetype = "Veteran"  // Frank - Sanctuary Funding
	ArchetypeMystic   Archetype = "Mystic"   // Tartaria - Reality Distortion
	ArchetypeShowman  Archetype = "Showman"  // Marrash - Content/Attention
	ArchetypeRedeemed Archetype = "Redeemed" // Simón - Survival/Expiation
)

// Prisoner represents the state of a participant in the game.
type Prisoner struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Archetype Archetype `json:"archetype"`

	// Vitals
	Hunger  int `json:"hunger"`  // 0-100 (0 = starving)
	Thirst  int `json:"thirst"`  // 0-100
	Sanity  int `json:"sanity"`  // 0-100 (0 = breakdown)
	Dignity int `json:"dignity"` // 0-100 (privacy-based)

	// Social
	Loyalty int `json:"loyalty"` // -100 to 100 (towards cellmate)
	Empathy int `json:"empathy"` // Hidden stat for betrayal checks

	// State
	IsSleeper  bool `json:"is_sleeper"`   // Disconnected player
	DayInGame  int  `json:"day_in_game"`  // 1-21
	IsWithdraw bool `json:"is_withdraw"`  // Simón's 5-day debuff
}

// NewPrisoner creates a fresh prisoner with default starting stats.
func NewPrisoner(id, name string, archetype Archetype) *Prisoner {
	return &Prisoner{
		ID:        id,
		Name:      name,
		Archetype: archetype,
		Hunger:    100,
		Thirst:    100,
		Sanity:    100,
		Dignity:   100,
		Loyalty:   50,
		Empathy:   50,
		IsSleeper: false,
		DayInGame: 1,
		IsWithdraw: archetype == ArchetypeRedeemed, // Simón starts in withdrawal
	}
}
