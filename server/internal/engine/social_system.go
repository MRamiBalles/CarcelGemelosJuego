package engine

import (
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// SocialSystem manages Duo dynamics, betrayal logic, and the Day 21 Dilemma.
type SocialSystem struct {
	eventLog  *events.EventLog
	logger    *logger.Logger
	prisoners map[string]*prisoner.Prisoner

	// Track dilemma decisions per duo: CellID -> PrisonerID -> Decision
	dilemmaState map[string]map[string]string
}

// DilemmaDecisionPayload holds the player's choice on Day 21.
type DilemmaDecisionPayload struct {
	PrisonerID string `json:"prisoner_id"`
	Decision   string `json:"decision"` // "BETRAY" or "COLLABORATE"
}

// EmotePayload holds info about social interactions.
type EmotePayload struct {
	ActorID   string `json:"actor_id"`
	EmoteType string `json:"emote_type"` // e.g., "AGGRESSIVE", "FRIENDLY"
	TargetID  string `json:"target_id"`
}

// SocialActionPayload describes a social action between prisoners.
type SocialActionPayload struct {
	ActorID    string `json:"actor_id"`
	TargetID   string `json:"target_id"`
	ActionType string `json:"action_type"`
}

// NewSocialSystem creates a new social manager.
func NewSocialSystem(eventLog *events.EventLog, log *logger.Logger) *SocialSystem {
	return &SocialSystem{
		eventLog:     eventLog,
		logger:       log,
		prisoners:    make(map[string]*prisoner.Prisoner),
		dilemmaState: make(map[string]map[string]string),
	}
}

// RegisterPrisoner adds a prisoner to be tracked.
func (ss *SocialSystem) RegisterPrisoner(p *prisoner.Prisoner) {
	ss.prisoners[p.ID] = p
}

// OnTimeTick handles periodic social effects like Toxic Duo proximity drain and Loyalty buffs.
func (ss *SocialSystem) OnTimeTick(event events.GameEvent) {
	// Every hour (or tick), apply loyalty buffs/debuffs
	for _, p := range ss.prisoners {
		// Toxic Duo (Labrador/Ylenia): "Mental Wear" Proximity
		// In a real implementation, we'd check distance. For now, if they are in the same cell, apply drain.
		if p.HasTrait(prisoner.TraitBadRomance) {
			cellmate := ss.getCellmate(p)
			if cellmate != nil && cellmate.CellID == p.CellID {
				// Distance < Threshold simulation:
				p.Sanity -= 2
				if p.Sanity < 0 {
					p.Sanity = 0
				}
			}
		}

		// Loyalty Bar Buffs (T038)
		if p.Loyalty >= 80 {
			// High loyalty provides minor passive sanity regen
			p.Sanity += 1
			if p.Sanity > 100 {
				p.Sanity = 100
			}
		}
	}
}

// OnAggressiveEmote handles Toxic Duo "Hype" generation.
func (ss *SocialSystem) OnAggressiveEmote(event events.GameEvent) {
	payload, ok := event.Payload.(EmotePayload)
	if !ok {
		return
	}

	actor, exists := ss.prisoners[payload.ActorID]
	if !exists {
		return
	}

	// Toxic Duo (Bad Romance) Hype generation
	if actor.HasTrait(prisoner.TraitBadRomance) && payload.EmoteType == "AGGRESSIVE" {
		target, ok := ss.prisoners[payload.TargetID]
		if ok && target.CellID == actor.CellID {
			// Generates Hype (Shared Pot)
			actor.PotContribution += 5.0
			target.PotContribution += 5.0
			ss.logger.Event("HYPE_GENERATED", actor.CellID, "Arguing in cell generates €10")
		}
	}
}

// OnFinalDilemmaDecision tracks responses and resolves the endgame dilemma.
func (ss *SocialSystem) OnFinalDilemmaDecision(event events.GameEvent) {
	payload, ok := event.Payload.(DilemmaDecisionPayload)
	if !ok {
		return
	}

	p, exists := ss.prisoners[payload.PrisonerID]
	if !exists {
		return
	}

	// Initialize state tracking for the duo
	if ss.dilemmaState[p.CellID] == nil {
		ss.dilemmaState[p.CellID] = make(map[string]string)
	}

	// Register decision
	ss.dilemmaState[p.CellID][p.ID] = payload.Decision
	ss.logger.Info("Dilemma Decision locked for " + p.ID)

	cellmate := ss.getCellmate(p)
	if cellmate != nil {
		cellmateDecision, cellmateDecided := ss.dilemmaState[p.CellID][cellmate.ID]
		if cellmateDecided {
			// RESOLVE DILEMMA
			ss.resolveDilemma(p, payload.Decision, cellmate, cellmateDecision)
		}
	}
}

// resolveDilemma applies the classic Prisoner's Dilemma logic to the PotContribution (Winnings).
func (ss *SocialSystem) resolveDilemma(p1 *prisoner.Prisoner, dec1 string, p2 *prisoner.Prisoner, dec2 string) {
	ss.logger.Event("DILEMMA_RESOLUTION", p1.CellID, "P1:"+dec1+" vs P2:"+dec2)

	totalPot := p1.PotContribution + p2.PotContribution

	var p1Outcome, p2Outcome string

	if dec1 == "COLLABORATE" && dec2 == "COLLABORATE" {
		// 50/50 Split
		p1.PotContribution = totalPot / 2
		p2.PotContribution = totalPot / 2
		p1Outcome = "SHARED"
		p2Outcome = "SHARED"
	} else if dec1 == "BETRAY" && dec2 == "COLLABORATE" {
		p1.PotContribution = totalPot
		p2.PotContribution = 0
		p1Outcome = "WON_ALL"
		p2Outcome = "LOST_ALL"
	} else if dec1 == "COLLABORATE" && dec2 == "BETRAY" {
		p1.PotContribution = 0
		p2.PotContribution = totalPot
		p1Outcome = "LOST_ALL"
		p2Outcome = "WON_ALL"
	} else {
		// BETRAY / BETRAY - House wins
		p1.PotContribution = 0
		p2.PotContribution = 0
		p1Outcome = "LOST_ALL"
		p2Outcome = "LOST_ALL"
	}

	// Emit resolution events for DB / UI
	ss.emitResolution(p1, p2, dec1, p1Outcome)
	ss.emitResolution(p2, p1, dec2, p2Outcome)
}

func (ss *SocialSystem) emitResolution(actor *prisoner.Prisoner, target *prisoner.Prisoner, decision string, outcome string) {
	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      "DILEMMA_OUTCOME", // Internal tracking type
		ActorID:   actor.ID,
		TargetID:  target.ID,
		Payload: map[string]interface{}{
			"decision": decision,
			"outcome":  outcome,
			"winnings": actor.PotContribution,
		},
		GameDay: 21,
	}
	ss.eventLog.Append(event)
}

func (ss *SocialSystem) getCellmate(p *prisoner.Prisoner) *prisoner.Prisoner {
	for _, cellmate := range ss.prisoners {
		if cellmate.CellID == p.CellID && cellmate.ID != p.ID && !cellmate.IsSleeper { // Simplified ignoring dead players
			return cellmate
		}
	}
	return nil
}
