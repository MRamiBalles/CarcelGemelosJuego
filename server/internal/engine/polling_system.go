package engine

import (
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// PollPayload is the data for EventTypePollCreated.
type PollPayload struct {
	PollID      string   `json:"poll_id"`
	Question    string   `json:"question"`
	Options     []string `json:"options"` // Usually Prisoner IDs or Names
	DurationSec int      `json:"duration_sec"`
	RewardType  string   `json:"reward_type"` // SUSHI, TORTURE, ISOLATION
}

// PollResolvedPayload is the data for EventTypePollResolved.
type PollResolvedPayload struct {
	PollID       string         `json:"poll_id"`
	WinnerOption string         `json:"winner_option"`
	Results      map[string]int `json:"results"`
	RewardType   string         `json:"reward_type"`
}

type PlayingPoll struct {
	Payload PollPayload
	Votes   map[string]int
	GameDay int
}

// PollingSystem handles Twitch-style audience interventions.
type PollingSystem struct {
	prisoners   map[string]*prisoner.Prisoner
	eventLog    *events.EventLog
	logger      *logger.Logger
	activePolls map[string]*PlayingPoll
}

func NewPollingSystem(eventLog *events.EventLog, log *logger.Logger) *PollingSystem {
	ps := &PollingSystem{
		prisoners:   make(map[string]*prisoner.Prisoner),
		eventLog:    eventLog,
		logger:      log,
		activePolls: make(map[string]*PlayingPoll),
	}
	return ps
}

func (ps *PollingSystem) RegisterPrisoner(p *prisoner.Prisoner) {
	ps.prisoners[p.ID] = p
}

// OnPollCreated starts tracking a poll and spins up a timer to resolve it.
func (ps *PollingSystem) OnPollCreated(event events.GameEvent) {
	payload, ok := event.Payload.(PollPayload)
	if !ok {
		// Map fallback for DB recovery
		m, mapped := event.Payload.(map[string]interface{})
		if mapped {
			payload.PollID, _ = m["poll_id"].(string)
			payload.Question, _ = m["question"].(string)
			payload.RewardType, _ = m["reward_type"].(string)
			durFloat, _ := m["duration_sec"].(float64)
			payload.DurationSec = int(durFloat)

			opts, _ := m["options"].([]interface{})
			for _, o := range opts {
				if str, isStr := o.(string); isStr {
					payload.Options = append(payload.Options, str)
				}
			}
		} else {
			return
		}
	}

	poll := &PlayingPoll{
		Payload: payload,
		Votes:   make(map[string]int),
		GameDay: event.GameDay,
	}

	// Initialize vote counts to 0
	for _, opt := range payload.Options {
		poll.Votes[opt] = 0
	}

	ps.activePolls[payload.PollID] = poll
	ps.logger.Info("Audience Poll Started: " + payload.Question)

	// Spin off asynchronous resolver
	time.AfterFunc(time.Duration(payload.DurationSec)*time.Second, func() {
		ps.resolvePoll(payload.PollID)
	})
}

// CastVote manually accepts a vote directly (not event-sourced to prevent spam).
func (ps *PollingSystem) CastVote(pollID, option string) {
	poll, exists := ps.activePolls[pollID]
	if !exists {
		return
	}

	// Ensure option is valid
	valid := false
	for _, o := range poll.Payload.Options {
		if o == option {
			valid = true
			break
		}
	}
	if !valid {
		return
	}

	poll.Votes[option]++
}

func (ps *PollingSystem) resolvePoll(pollID string) {
	poll, exists := ps.activePolls[pollID]
	if !exists {
		return
	}

	// Calculate winner
	winner := ""
	maxVotes := -1
	for opt, count := range poll.Votes {
		if count > maxVotes {
			maxVotes = count
			winner = opt
		}
	}

	// Emit Resolution
	resolveEvent := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypePollResolved,
		ActorID:   "SYSTEM_AUDIENCE",
		GameDay:   poll.GameDay,
		Payload: PollResolvedPayload{
			PollID:       pollID,
			WinnerOption: winner,
			Results:      poll.Votes,
			RewardType:   poll.Payload.RewardType,
		},
	}

	ps.eventLog.Append(resolveEvent)
	delete(ps.activePolls, pollID)
	ps.logger.Info("Poll " + pollID + " Resolved. Winner: " + winner)
}

// OnPollResolved applies the game mechanics to the winning prisoner.
func (ps *PollingSystem) OnPollResolved(event events.GameEvent) {
	payload, ok := event.Payload.(PollResolvedPayload)
	if !ok {
		// Map fallback for DB recovery
		m, mapped := event.Payload.(map[string]interface{})
		if mapped {
			payload.PollID, _ = m["poll_id"].(string)
			payload.WinnerOption, _ = m["winner_option"].(string)
			payload.RewardType, _ = m["reward_type"].(string)
		} else {
			return
		}
	}

	// Usually WinnerOption is the PrisonerID (or Name, simplify to ID for now)
	target, exists := ps.prisoners[payload.WinnerOption]
	if !exists {
		return
	}

	switch payload.RewardType {
	case "SUSHI":
		// Restore Sanity and Dignity
		target.Sanity += 40
		if target.Sanity > 100 {
			target.Sanity = 100
		}
		target.Dignity += 20
		if target.Dignity > 100 {
			target.Dignity = 100
		}
		ps.logger.Event("AUDIENCE_REWARD", target.ID, "Received Premium Meal (Sushi)")

	case "TORTURE":
		// Destroy sanity
		target.Sanity -= 50
		if target.Sanity < 0 {
			target.Sanity = 0
		}
		ps.logger.Event("AUDIENCE_PUNISH", target.ID, "Selected for intense torture")

	case "ISOLATION":
		// Send to isolation via event
		isoEvent := events.GameEvent{
			ID:        events.GenerateEventID(),
			Timestamp: time.Now(),
			Type:      events.EventTypeIsolationChanged,
			ActorID:   "SYSTEM_AUDIENCE",
			TargetID:  target.ID,
			GameDay:   event.GameDay,
			Payload: IsolationChangePayload{
				TargetID:   target.ID,
				IsIsolated: true,
			},
		}
		ps.eventLog.Append(isoEvent)
	}
}
