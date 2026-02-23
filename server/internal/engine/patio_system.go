package engine

import (
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// PatioChallengePayload tracks the outcome of the daily challenge.
type PatioChallengePayload struct {
	ParticipantID    string `json:"participant_id"`
	StaminaSpent     int    `json:"stamina_spent"`
	PotContribution  int    `json:"pot_contribution"`
	SanityMultiplier int    `json:"sanity_multiplier"`
}

// PatioSystem triggers a daily high-risk event at 12:00 testing resource management.
type PatioSystem struct {
	prisoners map[string]*prisoner.Prisoner
	eventLog  *events.EventLog
	logger    *logger.Logger
}

func NewPatioSystem(eventLog *events.EventLog, log *logger.Logger) *PatioSystem {
	ps := &PatioSystem{
		prisoners: make(map[string]*prisoner.Prisoner),
		eventLog:  eventLog,
		logger:    log,
	}

	return ps
}

func (ps *PatioSystem) RegisterPrisoner(p *prisoner.Prisoner) {
	ps.prisoners[p.ID] = p
}

// OnTimeTick listens for 12:00 PM daily to trigger the event
func (ps *PatioSystem) OnTimeTick(event events.GameEvent) {
	tickPayload, ok := event.Payload.(TimeTickPayload)
	if !ok {
		// handle DB map string recovery
		m, mapped := event.Payload.(map[string]interface{})
		if mapped {
			dayFloat, _ := m["game_day"].(float64)
			tickPayload.GameDay = int(dayFloat)
			hourFloat, _ := m["game_hour"].(float64)
			tickPayload.GameHour = int(hourFloat)
		} else {
			return
		}
	}

	ps.logger.Info("Starting Daily Patio Challenge (Game Day: %v)", tickPayload.GameDay)

	// Pick a random participant who is NOT isolated and NOT a sleeper (keep it simple for now)
	var validParticipants []*prisoner.Prisoner
	for _, p := range ps.prisoners {
		if !p.IsSleeper && !p.IsIsolated {
			validParticipants = append(validParticipants, p)
		}
	}

	if len(validParticipants) == 0 {
		ps.logger.Info("No valid participants for Patio Challenge.")
		return
	}

	// Simple pseudo-random selection based on GameDay logic
	// (usually we'd use math/rand, but this is deterministic enough for demo purposes)
	selectedIndex := tickPayload.GameDay % len(validParticipants)
	chosenOne := validParticipants[selectedIndex]

	// Apply Cost: Huge Stamina Drop (Simulated as Hunger here)
	// If Hunger hits 100, they start dying.
	staminaCost := 80
	chosenOne.Hunger += staminaCost
	if chosenOne.Hunger > 100 {
		chosenOne.Hunger = 100
	}

	// Calculate Reward: High risk, high reward
	// E.g., add to PotContribution (requires adding PotContribution to Prisoner struct or handled via traits)
	// We'll simulate a random multiplier based on their Dignity
	multiplier := 1
	if chosenOne.Dignity > 80 {
		multiplier = 2
	}

	rewardAmount := 500 * multiplier // $500 to the central jackpot
	// Modify pot contribution (we need a global pot, but for now we grant it directly to the player's stat)

	payload := PatioChallengePayload{
		ParticipantID:    chosenOne.ID,
		StaminaSpent:     staminaCost,
		PotContribution:  rewardAmount,
		SanityMultiplier: multiplier,
	}

	// We append a Social Action for the challenge record as it affects group dynamics
	ps.eventLog.Append(events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeSocialAction, // Reusing Social Action as the framework for "Doing something for the group"
		ActorID:   "SYSTEM_PATIO",
		TargetID:  chosenOne.ID,
		GameDay:   tickPayload.GameDay,
		Payload:   payload,
	})

	ps.logger.Event("PATIO_CHALLENGE", chosenOne.ID, "Survived the patio challenge. Cost: 80 Stamina.")
}
