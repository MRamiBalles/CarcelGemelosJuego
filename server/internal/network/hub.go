package network

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
	"github.com/gorilla/websocket"
)

// Client represents an active WebSocket connection.
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
	logger     *logger.Logger
}

// NewHub initializes a new WebSocket Hub.
func NewHub(log *logger.Logger) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		logger:     log,
	}
}

// Run starts the Hub's main loop to handle client connections and broadcasts.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("WebSocket Hub shutting down.")
			return
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Info("New WebSocket client connected")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.logger.Info("WebSocket client disconnected")
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastEvent takes a GameEvent, serializes it to JSON, and sends it to all connected clients.
func (h *Hub) BroadcastEvent(event events.GameEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		h.logger.Error("Failed to serialize GameEvent for WebSocket broadcast: %v", err)
		return
	}
	h.broadcast <- payload
}

// StartEventPoller spans a goroutine to poll the EventLog and push new events to the Hub.
// This allows the Hub to run independently from the Engine's dispatch loop while picking up the same events.
func (h *Hub) StartEventPoller(ctx context.Context, eventLog *events.EventLog) {
	go func() {
		pollInterval := time.NewTicker(200 * time.Millisecond)
		defer pollInterval.Stop()

		lastProcessedEvent := 0

		for {
			select {
			case <-ctx.Done():
				return
			case <-pollInterval.C:
				allEvents := eventLog.Replay()
				newEventsCount := len(allEvents) - lastProcessedEvent

				if newEventsCount > 0 {
					newEvents := allEvents[lastProcessedEvent:]
					for _, event := range newEvents {
						// Only broadcast events that are relevant for the frontend (or all of them depending on need)
						h.BroadcastEvent(event)
					}
					lastProcessedEvent = len(allEvents)
				}
			}
		}
	}()
}
