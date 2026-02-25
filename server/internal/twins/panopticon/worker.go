package panopticon

import (
	"context"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/engine"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/action"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/cognition"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/twins/perception"
)

// Worker is the central nervous system of the Twins AI.
// It subscribes to the EventLog for reactive triggers (like Toilet Use)
// and runs a periodic cognitive loop for strategic LLM decisions.
type Worker struct {
	eventLog  *events.EventLog
	eng       *engine.Engine
	perceiver *perception.Perceiver
	cognitor  *cognition.Cognitor
	executor  *action.Executor
	logger    *logger.Logger

	lastProcessedEvent int
}

func NewWorker(el *events.EventLog, eng *engine.Engine, p *perception.Perceiver, c *cognition.Cognitor, a *action.Executor, log *logger.Logger) *Worker {
	return &Worker{
		eventLog:           el,
		eng:                eng,
		perceiver:          p,
		cognitor:           c,
		executor:           a,
		logger:             log,
		lastProcessedEvent: 0,
	}
}

// Start begins the reactive and proactive loops of the Panóptico.
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("Panóptico Worker initialized. The Twins are watching.")

	go w.eventListenerLoop(ctx)
	go w.cognitiveLoop(ctx)
}

// eventListenerLoop reacts to specific game events instantly.
func (w *Worker) eventListenerLoop(ctx context.Context) {
	pollInterval := time.NewTicker(200 * time.Millisecond)
	defer pollInterval.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Panóptico EventListener stopped.")
			return
		case <-pollInterval.C:
			allEvents := w.eventLog.Replay()
			newEventsCount := len(allEvents) - w.lastProcessedEvent

			if newEventsCount > 0 {
				newEvents := allEvents[w.lastProcessedEvent:]
				for _, event := range newEvents {
					w.reactToEvent(event)
				}
				w.lastProcessedEvent = len(allEvents)
			}
		}
	}
}

// reactToEvent evaluates deterministic Rules-as-Code tortures.
func (w *Worker) reactToEvent(event events.GameEvent) {
	if event.Type == events.EventTypeToiletUse {
		payload, ok := event.Payload.(events.ToiletUsePayload)
		if ok && payload.IsObserved {
			w.logger.Event("PANOPTICON_JUDGMENT", event.ActorID, "Observed toilet use. Penalizing dignity.")
			// Inyectar penalización al usar el inodoro a la vista
			w.eng.GetPrisoners()[event.ActorID].Dignity -= 10
			w.eng.GetPrisoners()[event.ActorID].Sanity -= 5

			// Log the punishment
			w.eventLog.Append(events.GameEvent{
				ID:        events.GenerateEventID(),
				Timestamp: time.Now(),
				Type:      events.EventTypeSanityChange,
				ActorID:   "SYSTEM_PANOPTICON",
				TargetID:  event.ActorID,
				Payload: map[string]interface{}{
					"reason": "Toilet use observed by cellmate",
					"amount": -5,
				},
				GameDay: event.GameDay,
			})
		}
	}
}

// cognitiveLoop runs the LLM strategy evaluation periodically.
func (w *Worker) cognitiveLoop(ctx context.Context) {
	// Evaluates the prison every 5 real minutes (approx. 10 game hours)
	evalInterval := time.NewTicker(5 * time.Minute)
	defer evalInterval.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Panóptico CognitiveLoop stopped.")
			return
		case <-evalInterval.C:
			w.logger.Info("Panóptico initiating cognitive evaluation...")

			// Reconstruct state from internal perceiver
			state, err := w.perceiver.BuildPrisonState(ctx, "GAME_1")
			if err != nil {
				w.logger.Error("Panóptico perception failed: " + err.Error())
				continue
			}

			// 2. Decide
			decision, err := w.cognitor.Decide(ctx, state)
			if err != nil {
				w.logger.Error("Panóptico cognition failed: " + err.Error())
				continue
			}

			// 3. Act
			w.logger.Info("Panóptico executing decision...")
			err = w.executor.Execute(ctx, decision)
			if err != nil {
				w.logger.Error("Failed to execute action: " + err.Error())
			}
		}
	}
}
