// Package network provides WebSocket server functionality.
// T013: WebSocket server for real-time player-server communication.
//
// ARCHITECTURAL RULE: This package is AGNOSTIC to game logic.
// It only knows how to route messages; game logic lives in domain/engine.
package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// MessageType defines the category of WebSocket messages.
type MessageType string

const (
	MsgTypeAuth       MessageType = "AUTH"
	MsgTypeGameState  MessageType = "GAME_STATE"
	MsgTypeAction     MessageType = "ACTION"
	MsgTypeEvent      MessageType = "EVENT"
	MsgTypePing       MessageType = "PING"
	MsgTypePong       MessageType = "PONG"
	MsgTypeError      MessageType = "ERROR"
	MsgTypeRecap      MessageType = "RECAP"
)

// Message is the standard WebSocket message envelope.
type Message struct {
	Type      MessageType            `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// Client represents a connected player.
type Client struct {
	ID         string
	PrisonerID string
	GameID     string
	Send       chan []byte
	IsAuth     bool
	LastPing   time.Time
}

// Hub manages all WebSocket connections.
type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMessage
	mu         sync.RWMutex
	logger     *logger.Logger
}

// BroadcastMessage targets specific clients or all clients in a game.
type BroadcastMessage struct {
	GameID  string  // Target game (empty = all)
	Targets []string // Specific client IDs (empty = all in game)
	Message Message
}

// NewHub creates a new WebSocket hub.
func NewHub(log *logger.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMessage, 256),
		logger:     log,
	}
}

// Run starts the hub's main loop. Call in a goroutine.
func (h *Hub) Run(ctx context.Context) {
	h.logger.Info("WebSocket Hub started")
	
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("WebSocket Hub shutting down")
			return
			
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			h.logger.Event("WS_CONNECT", client.ID, "Client connected")
			
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			h.logger.Event("WS_DISCONNECT", client.ID, "Client disconnected")
			
		case msg := <-h.broadcast:
			h.handleBroadcast(msg)
		}
	}
}

// handleBroadcast sends a message to targeted clients.
func (h *Hub) handleBroadcast(bm BroadcastMessage) {
	data, err := json.Marshal(bm.Message)
	if err != nil {
		h.logger.Error("Failed to marshal broadcast: " + err.Error())
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		// Filter by game
		if bm.GameID != "" && client.GameID != bm.GameID {
			continue
		}
		
		// Filter by specific targets
		if len(bm.Targets) > 0 {
			found := false
			for _, t := range bm.Targets {
				if t == client.ID || t == client.PrisonerID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		
		select {
		case client.Send <- data:
		default:
			// Client buffer full, close connection
			close(client.Send)
			delete(h.clients, client.ID)
		}
	}
}

// BroadcastToGame sends a message to all players in a game.
func (h *Hub) BroadcastToGame(gameID string, msg Message) {
	h.broadcast <- BroadcastMessage{
		GameID:  gameID,
		Message: msg,
	}
}

// SendToClient sends a message to a specific client.
func (h *Hub) SendToClient(clientID string, msg Message) error {
	h.mu.RLock()
	client, ok := h.clients[clientID]
	h.mu.RUnlock()
	
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}
	
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	client.Send <- data
	return nil
}

// GetConnectedPrisoners returns IDs of all connected prisoners in a game.
func (h *Hub) GetConnectedPrisoners(gameID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	var ids []string
	for _, c := range h.clients {
		if c.GameID == gameID && c.IsAuth {
			ids = append(ids, c.PrisonerID)
		}
	}
	return ids
}

// IsClientOnline checks if a prisoner is currently connected.
func (h *Hub) IsClientOnline(prisonerID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	for _, c := range h.clients {
		if c.PrisonerID == prisonerID {
			return true
		}
	}
	return false
}

// ServeHTTP is the HTTP handler for WebSocket upgrades.
// NOTE: Requires gorilla/websocket in production. This is a stub.
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Integrate gorilla/websocket upgrade
	// For now, return upgrade required response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUpgradeRequired)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "WebSocket upgrade required",
		"hint":  "Connect via ws:// protocol",
	})
}
