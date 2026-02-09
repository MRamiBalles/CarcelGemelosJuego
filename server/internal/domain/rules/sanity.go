// Package rules contains the pure calculation logic for game mechanics.
// This package is PURE and must NOT import any infrastructure packages.
package rules

import "github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"

// SanityDrainParams holds the parameters for a sanity drain event.
type SanityDrainParams struct {
	IsMystic     bool
	NoiseLevel   int // 1-3 (intensity of the Twins' torture)
	IsPrivacyHit bool
}

// CalculateSanityDrain computes the sanity loss based on events.
// This is the core of the "The Twins' Noise" system.
func CalculateSanityDrain(p *prisoner.Prisoner, params SanityDrainParams) int {
	baseDrain := params.NoiseLevel * 5 // Base drain is 5-15 per event

	// Mystic mitigation
	if p.Archetype == prisoner.ArchetypeMystic && p.Sanity > 20 {
		baseDrain /= 2
	}

	// Privacy hit (Visible Toilet)
	if params.IsPrivacyHit {
		baseDrain += 5
	}

	// Veteran bonus when alone (negative drain = regen)
	// This would be triggered from the engine with context, not here directly.

	return baseDrain
}

// CheckBetrayalSuccess determines if a betrayal attempt succeeds on Day 21.
// Returns true if the player has sufficient "cold blood" to betray.
func CheckBetrayalSuccess(p *prisoner.Prisoner, rollValue int) bool {
	// Threshold: lower Sanity + lower Empathy = easier to betray
	threshold := p.Empathy + (100 - p.Sanity)
	return rollValue > threshold
}

// ProcessWithdrawal applies the "Cold Turkey" logic for the Redeemed class.
// Returns buffs/debuffs in the form of (sanityMod, staminaMod).
func ProcessWithdrawal(p *prisoner.Prisoner) (sanityMod int, staminaMod int) {
	if p.Archetype != prisoner.ArchetypeRedeemed {
		return 0, 0
	}

	if p.DayInGame <= 5 && p.IsWithdraw {
		// Debuff phase
		return -5, -10
	}

	if p.DayInGame > 5 && !p.IsWithdraw { // Survived withdrawal
		// CLARITY buff
		return 20, 10
	}

	return 0, 0
}
