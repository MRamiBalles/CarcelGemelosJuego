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
	ArchetypeChaos    Archetype = "Chaos"    // Aída - Energy/Noise
)

// TraitID identifies a specific passive or active ability.
type TraitID string

const (
	TraitIronStomach   TraitID = "IronStomach"   // Immune to filth/food poisoning
	TraitMisanthrope   TraitID = "Misanthrope"   // Sanity regen when alone
	TraitBreatharian   TraitID = "Breatharian"   // Cannot eat solids, hunger frozen
	TraitContentCreator TraitID = "ContentCreator" // Bonus interaction with cameras
	TraitInsomniac     TraitID = "Insomniac"     // Needs less sleep
)

// StateID identifies a temporary status effect.
type StateID string

const (
	StateWithdrawal StateID = "Withdrawal" // Stat reduction
	StateMeditating StateID = "Meditating" // Stunned but regenerating
	StateStarvation StateID = "Starvation" // HP Drain
	StateFacingWall StateID = "FacingWall" // Avoiding eye contact
)

// Prisoner represents the state of a participant in the game.
type Prisoner struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Archetype Archetype `json:"archetype"`
	CellID    string    `json:"cell_id"` // Grouping identifier

	// Vitals
	Hunger  int `json:"hunger"`  // 0-100 (0 = starving)
	Thirst  int `json:"thirst"`  // 0-100
	Sanity  int `json:"sanity"`  // 0-100 (0 = breakdown)
	Dignity int `json:"dignity"` // 0-100 (privacy-based)
	HP      int `json:"hp"`      // 0-100 (Physical Health)
	Stamina int `json:"stamina"` // 0-100 (Energy/Fatigue)

	// Economics
	PotContribution float64 `json:"pot_contribution"` // Individual winnings

	// Social
	Loyalty int `json:"loyalty"` // -100 to 100 (towards cellmate)
	Empathy int `json:"empathy"` // Hidden stat for betrayal checks

	// State
	IsSleeper bool            `json:"is_sleeper"`  // Disconnected player
	DayInGame int             `json:"day_in_game"` // 1-21
	Traits    []TraitID       `json:"traits"`      // Active passive abilities
	States    map[StateID]int `json:"states"`      // Temporary effects (ID -> Duration Ticks)
}

// NewPrisoner creates a fresh prisoner with default starting stats based on Archetype.
func NewPrisoner(id, name string, archetype Archetype) *Prisoner {
	p := &Prisoner{
		ID:              id,
		Name:            name,
		Archetype:       archetype,
		Hunger:          100,
		Thirst:          100,
		Sanity:          100,
		Dignity:         100,
		HP:              100,
		Loyalty:         50,
		Empathy:         50,
		PotContribution: 0,
		IsSleeper:       false,
		DayInGame:       1,
		Traits:          []TraitID{},
		States:          make(map[StateID]int),
	}

	// Apply Archetype defaults
	switch archetype {
	case ArchetypeVeteran:
		p.Traits = append(p.Traits, TraitIronStomach, TraitMisanthrope)
	case ArchetypeMystic:
		p.Traits = append(p.Traits, TraitBreatharian)
	case ArchetypeRedeemed:
		p.AddState(StateWithdrawal, 5*24*60) // 5 Days in minutes (assuming 1 tick = 1 min)
	case ArchetypeChaos:
		p.Traits = append(p.Traits, TraitInsomniac)
	}

	return p
}

// Helper methods

func (p *Prisoner) AddState(state StateID, duration int) {
	if p.States == nil {
		p.States = make(map[StateID]int)
	}
	p.States[state] = duration
}

func (p *Prisoner) HasState(state StateID) bool {
	_, ok := p.States[state]
	return ok
}

func (p *Prisoner) HasTrait(trait TraitID) bool {
	for _, t := range p.Traits {
		if t == trait {
			return true
		}
	}
	return false
}

func (p *Prisoner) TickStates() {
	for id, duration := range p.States {
		if duration > 0 {
			p.States[id] = duration - 1
		} else {
			delete(p.States, id)
		}
	}
}

// Legacy compatibility (to be refactored)
func (p *Prisoner) IsWithdraw() bool {
	return p.HasState(StateWithdrawal)
}
