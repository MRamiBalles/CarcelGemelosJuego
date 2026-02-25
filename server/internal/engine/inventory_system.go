package engine

import (
	"fmt"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/item"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

// ItemTransferPayload represents an item moving between prisoners or from system.
type ItemTransferPayload struct {
	FromPrisonerID string        `json:"from_id"`
	ToPrisonerID   string        `json:"to_id"`
	ItemType       item.ItemType `json:"item_type"`
	Quantity       int           `json:"quantity"`
	Reason         string        `json:"reason"`
}

// ItemConsumedPayload represents an item being used.
type ItemConsumedPayload struct {
	PrisonerID string        `json:"prisoner_id"`
	ItemType   item.ItemType `json:"item_type"`
	Quantity   int           `json:"quantity"`
}

// InventorySystem handles the logic for item transfers, consumption, and limits.
type InventorySystem struct {
	eventLog  *events.EventLog
	prisoners map[string]*prisoner.Prisoner
	logger    *logger.Logger
}

func NewInventorySystem(el *events.EventLog, log *logger.Logger) *InventorySystem {
	return &InventorySystem{
		eventLog:  el,
		prisoners: make(map[string]*prisoner.Prisoner),
		logger:    log,
	}
}

func (is *InventorySystem) RegisterPrisoner(p *prisoner.Prisoner) {
	is.prisoners[p.ID] = p
}

// TransferItem moves an item from one prisoner to another.
// Use "SYSTEM" as FromPrisonerID for game rewards (e.g. Patio success).
func (is *InventorySystem) TransferItem(fromID, toID string, itemType item.ItemType, quantity int, reason string) error {
	receiver, ok := is.prisoners[toID]
	if !ok {
		return fmt.Errorf("receiver %s not found", toID)
	}

	if fromID != "SYSTEM" {
		sender, ok := is.prisoners[fromID]
		if !ok {
			return fmt.Errorf("sender %s not found", fromID)
		}
		if !sender.RemoveItem(itemType, quantity) {
			return fmt.Errorf("sender %s does not have %d of %s", fromID, quantity, itemType)
		}
	}

	receiver.AddItem(itemType, quantity)

	// Emit event
	payload := ItemTransferPayload{
		FromPrisonerID: fromID,
		ToPrisonerID:   toID,
		ItemType:       itemType,
		Quantity:       quantity,
		Reason:         reason,
	}
	event := events.GameEvent{
		ID:         events.GenerateEventID(),
		Timestamp:  time.Now(),
		Type:       events.EventTypeItemTransfer,
		ActorID:    fromID,
		TargetID:   toID,
		Payload:    payload,
		GameDay:    receiver.DayInGame,
		IsRevealed: true, // Trades are public to the audience
	}
	is.eventLog.Append(event)
	is.logger.Info(fmt.Sprintf("[INVENTORY] Transfer: %s gave %d %s to %s (%s)", fromID, quantity, itemType, toID, reason))

	return nil
}

// ConsumeItem uses an item from a prisoner's inventory and applies its effects (if any internal to inventory logic).
// Most effects (Hunger/Thirst) will be handled by MetabolismSystem reacting to the ITEM_CONSUMED event.
func (is *InventorySystem) ConsumeItem(prisonerID string, itemType item.ItemType) error {
	p, ok := is.prisoners[prisonerID]
	if !ok {
		return fmt.Errorf("prisoner %s not found", prisonerID)
	}

	if !p.RemoveItem(itemType, 1) {
		return fmt.Errorf("prisoner %s does not have a %s to consume", prisonerID, itemType)
	}

	// Emit event so other systems (like Metabolism) can apply effects
	payload := ItemConsumedPayload{
		PrisonerID: prisonerID,
		ItemType:   itemType,
		Quantity:   1,
	}
	event := events.GameEvent{
		ID:         events.GenerateEventID(),
		Timestamp:  time.Now(),
		Type:       events.EventTypeItemConsumed,
		ActorID:    prisonerID,
		TargetID:   "",
		Payload:    payload,
		GameDay:    p.DayInGame,
		IsRevealed: true,
	}
	is.eventLog.Append(event)
	is.logger.Info(fmt.Sprintf("[INVENTORY] Consumed: %s used 1 %s", prisonerID, itemType))

	return nil
}
