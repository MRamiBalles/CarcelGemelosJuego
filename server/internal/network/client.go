package network

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// NewClient creates a new WebSocket client and returns it.
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

// PlayerAction represents an incoming command from the frontend.
type PlayerAction struct {
	Type       string          `json:"type"`        // "EAT", "STEAL", "TOILET", "SNITCH", etc.
	PrisonerID string          `json:"prisoner_id"` // Who triggered the action
	Payload    json.RawMessage `json:"payload"`     // Action-specific data
}

// Client object to hold connection status. Added Hub ref to allow unregister.
type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	send           chan []byte
	lastActionTime time.Time
}

// Register adds the client to the hub.
func (c *Client) Register() {
	c.hub.register <- c
}

// ReadPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var action PlayerAction
		if err := json.Unmarshal(message, &action); err != nil {
			c.hub.logger.Error("Failed to parse PlayerAction from WebSocket. err: " + err.Error())
			continue
		}

		c.handlePlayerAction(action)
	}
}

func (c *Client) handlePlayerAction(action PlayerAction) {
	// 1. Rate Limiting Check
	if time.Since(c.lastActionTime) < 15*time.Second {
		c.hub.logger.Warn("Rate limit exceeded for client action from " + action.PrisonerID)
		return
	}
	c.lastActionTime = time.Now()

	// 2. Fetch Actor
	eng := c.hub.engine
	prisoners := eng.GetPrisoners()
	actor, exists := prisoners[action.PrisonerID]
	if !exists {
		c.hub.logger.Error("PlayerAction from unknown prisoner: " + action.PrisonerID)
		return
	}

	// 3. Global Validations (e.g. Isolation, Death)
	if actor.HP <= 0 {
		return
	}
	if actor.HasState(prisoner.StateIsolated) && action.Type != "DILEMMA" {
		c.hub.logger.Warn("Isolated prisoner " + action.PrisonerID + " attempted action " + action.Type)
		return
	}

	// Route to specific action endpoints based on Phase F5 requirements
	switch action.Type {
	case "EAT", "DRINK":
		c.handleConsume(actor, action.Type, action.Payload)
	case "TOILET":
		c.handleToilet(actor)
	case "STEAL":
		c.handleSteal(actor, action.Payload)
	case "SNITCH":
		c.handleSnitch(actor, action.Payload)
	case "USE_RED_PHONE":
		c.handleRedPhone(actor)
	case "MEDITATE":
		c.handleMeditate(actor)
	case "USE_ORACLE":
		c.handleOracle(actor, action.Payload)
	case "GIVE_ELIXIR":
		// TODO: Tartaria Elixir logic
	default:
		c.hub.logger.Warn("Unknown PlayerAction type: " + action.Type)
	}
}

func (c *Client) handleConsume(actor *prisoner.Prisoner, actionType string, rawPayload []byte) {
	var parsed struct {
		ItemType string `json:"item_type"`
	}
	if err := json.Unmarshal(rawPayload, &parsed); err != nil {
		c.hub.logger.Warn("Failed to parse consume payload for " + actor.ID)
		return
	}

	// Fetch latest inventory from engine
	// Note: since F2, InventorySystem manages items. However, since we might read it directly from prisoner state,
	// let's assume we can emit an intent and let engine validate, or validate here. The plan says "Verificación estricta de ítem en inventario antes de aplicar el buff."
	// To do strict validation, we check the global state.
	hasItem := false
	for _, item := range actor.Inventory {
		if string(item.Type) == parsed.ItemType {
			hasItem = true
			break
		}
	}

	if !hasItem {
		c.hub.logger.Warn("Attempted to consume item not in inventory: " + parsed.ItemType)
		return
	}

	// Item exists, emit the event!
	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeItemConsumed,
		ActorID:   actor.ID,
		TargetID:  "",
		Payload:   parsed.ItemType,
		GameDay:   actor.DayInGame,
	}
	c.hub.engine.GetEventLog().Append(event)
	c.hub.logger.Event("PLAYER_ACTION_CONSUME", actor.ID, "Consumed "+parsed.ItemType)
}

func (c *Client) handleToilet(actor *prisoner.Prisoner) {
	// Emit Toilet Use Event
	payload := events.ToiletUsePayload{
		PrisonerID: actor.ID,
		IsObserved: false, // Updated by SanitySystem later
		WarningMsg: "¡Me estoy meando encima, date la vuelta a la pared!",
	}

	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeToiletUse,
		ActorID:   actor.ID,
		TargetID:  "CELL_" + actor.CellID, // Broadcasts warning to cell
		Payload:   payload,
		GameDay:   actor.DayInGame,
	}
	c.hub.engine.GetEventLog().Append(event)
	c.hub.logger.Event("PLAYER_ACTION_TOILET", actor.ID, "Used toilet in cell")
}

func (c *Client) handleSteal(actor *prisoner.Prisoner, rawPayload []byte) {
	// 1. Fatigue validation (F5 Lore)
	if actor.HasState(prisoner.StateExhausted) {
		// 80% chance to fail
		if rand.Float32() < 0.8 {
			c.hub.logger.Warn("STEAL FAILED: " + actor.Name + " was too exhausted and got caught!")

			// Directly punish loyalty or emit a betrayal event
			failEvent := events.GameEvent{
				ID:        events.GenerateEventID(),
				Timestamp: time.Now(),
				Type:      events.EventTypeBetrayal,
				ActorID:   actor.ID,
				TargetID:  "CELLMATE", // The logic resolves this later in SocialSystem
				Payload:   "FAILED_STEAL_EXHAUSTED",
				GameDay:   actor.DayInGame,
			}
			c.hub.engine.GetEventLog().Append(failEvent)
			return // Cannot proceed with steal
		}
	}

	// 2. Parse target from rawPayload and emit STEAL event for ChaosSystem
	var target struct {
		TargetID string `json:"target_id"`
	}
	if err := json.Unmarshal(rawPayload, &target); err == nil {
		event := events.GameEvent{
			ID:        events.GenerateEventID(),
			Timestamp: time.Now(),
			Type:      events.EventTypeSteal,
			ActorID:   actor.ID,
			TargetID:  target.TargetID,
			Payload:   nil, // Handled by chaos_system
			GameDay:   actor.DayInGame,
		}
		c.hub.engine.GetEventLog().Append(event)
	}
}

func (c *Client) handleSnitch(actor *prisoner.Prisoner, rawPayload []byte) {
	// 1. Lore: Panopticon cross-vision. Celdas enfrentadas.
	var snitch events.SnitchPayload
	if err := json.Unmarshal(rawPayload, &snitch); err != nil {
		return
	}

	// Emit SNITCH event
	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeSnitch,
		ActorID:   actor.ID,
		TargetID:  snitch.TargetID,
		Payload:   snitch,
		GameDay:   actor.DayInGame,
	}
	c.hub.engine.GetEventLog().Append(event)
	c.hub.logger.Event("PLAYER_ACTION_SNITCH", actor.ID, "Snitched on "+snitch.TargetID+" for "+snitch.Action)
}

func (c *Client) handleRedPhone(actor *prisoner.Prisoner) {
	// The phone rang (Panopticon action `EventTypeRedPhoneMessage`)
	// The first to answer gets a buff or debuff randomly.
	if rand.Float32() < 0.5 {
		actor.HP += 20 // Buff
		actor.Sanity += 20
		c.hub.logger.Event("RED_PHONE_BUFF", actor.ID, "Answered and received a buff")
	} else {
		actor.Sanity -= 30 // Debuff
		c.hub.logger.Event("RED_PHONE_DEBUFF", actor.ID, "Answered and received a sanity hit")
	}

	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeRedPhoneAnswer,
		ActorID:   actor.ID,
		TargetID:  "SYSTEM",
		Payload:   "ANSWERED_RED_PHONE",
		GameDay:   actor.DayInGame,
	}
	c.hub.engine.GetEventLog().Append(event)
}

func (c *Client) handleMeditate(actor *prisoner.Prisoner) {
	// Only Tartaria (ArchetypeMystic) can do this effectively
	if actor.Archetype != prisoner.ArchetypeMystic {
		return
	}

	// Lore: "Llevarlos a la Antártida"
	// Freezes own stamina drain (handled by adding Meditating state)
	actor.AddState(prisoner.StateMeditating, 60) // Lasts 1 hour (60 ticks)

	// Debuff contiguous cells (Simulate "Cold")
	// For simplicity, target all other cells except own
	cells := []string{}
	for _, p := range c.hub.engine.GetPrisoners() {
		if p.CellID != actor.CellID && p.CellID != "" {
			// Apply cold debuff (drain stamina to simulate shivering)
			p.Stamina -= 10
			if p.Stamina < 0 {
				p.Stamina = 0
			}

			// Append unique cell IDs to payload
			found := false
			for _, cid := range cells {
				if cid == p.CellID {
					found = true
				}
			}
			if !found {
				cells = append(cells, p.CellID)
			}
		}
	}

	payload := events.MeditatePayload{
		TargetCells: cells,
	}

	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeMeditate,
		ActorID:   actor.ID,
		TargetID:  "CONTIGUOUS_CELLS",
		Payload:   payload,
		GameDay:   actor.DayInGame,
	}
	c.hub.engine.GetEventLog().Append(event)
	c.hub.logger.Event("TARTARIA_ANTARCTICA", actor.ID, "Cast AoE Cold Meditation")
}

func (c *Client) handleOracle(actor *prisoner.Prisoner, rawPayload []byte) {
	// Tartaria reads painful secrets to cut heads
	if actor.Archetype != prisoner.ArchetypeMystic {
		return
	}

	var parsed struct {
		TargetID string `json:"target_id"`
	}
	if err := json.Unmarshal(rawPayload, &parsed); err != nil {
		return
	}

	payload := events.OracleUsePayload{
		TargetID: parsed.TargetID,
		Secret:   "Su compañero planea robarle esta noche", // Hardcoded painful truth
	}

	event := events.GameEvent{
		ID:        events.GenerateEventID(),
		Timestamp: time.Now(),
		Type:      events.EventTypeOracleUse,
		ActorID:   actor.ID,
		TargetID:  parsed.TargetID,
		Payload:   payload,
		GameDay:   actor.DayInGame,
	}
	c.hub.engine.GetEventLog().Append(event)
	c.hub.logger.Event("ORACLE_USED", actor.ID, "Revealed painful truth to "+parsed.TargetID)
	// Loyalty drop will be processed by Chaos/Social System
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
