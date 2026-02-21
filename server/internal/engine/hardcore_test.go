package engine

import (
	"testing"
	"time"

	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/domain/prisoner"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/events"
	"github.com/MRamiBalles/CarcelGemelosJuego/server/internal/platform/logger"
)

func TestToiletShame(t *testing.T) {
	// Setup
	el := events.NewEventLog()
	log := logger.NewLogger()
	ss := NewSanitySystem(el, log)

	// Create Actors
	user := prisoner.NewPrisoner("P1", "User", prisoner.ArchetypeVeteran)
	user.CellID = "CELL_101"

	witness := prisoner.NewPrisoner("P2", "Witness", prisoner.ArchetypeDeceiver)
	witness.CellID = "CELL_101"
	// Witness is NOT facing wall (default state)

	ss.RegisterPrisoner(user)
	ss.RegisterPrisoner(witness)

	// Act: User uses toilet
	payload := ToiletUsePayload{
		ActorID: "P1",
		CellID:  "CELL_101",
	}

	event := events.GameEvent{
		ID:        "TEST_EVT_1",
		Type:      events.EventTypeToiletUse,
		ActorID:   "P1",
		Payload:   payload,
		Timestamp: time.Now(),
	}

	ss.OnToiletUseEvent(event)

	// Assert: Both should lose Sanity
	if user.Dignity == 100 {
		t.Errorf("Expected User Dignity drop, got 100")
	}
	if user.Sanity == 100 {
		t.Errorf("Expected User Sanity drop (Shame), got 100")
	}
	if witness.Sanity == 100 {
		t.Errorf("Expected Witness Sanity drop (Privacy Violation), got 100")
	}
}

func TestLockdownSchedule(t *testing.T) {
	// Setup
	el := events.NewEventLog()
	log := logger.NewLogger()
	ls := NewLockdownSystem(el, log)

	// Act: Tick at 22:00
	tickPayload := TimeTickPayload{
		GameDay:  1,
		GameHour: 22,
	}
	event := events.GameEvent{
		Type:    events.EventTypeTimeTick,
		Payload: tickPayload,
	}

	ls.OnTimeTick(event)

	// Assert: Check event log for DOOR_LOCK
	found := false
	for _, e := range el.Replay() {
		if e.Type == events.EventTypeDoorLock {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected DOOR_LOCK event at 22:00")
	}
}

func TestMysticDiet(t *testing.T) {
	// Setup
	el := events.NewEventLog()
	log := logger.NewLogger()
	ms := NewMetabolismSystem(el, log)

	// Create Mystic
	mystic := prisoner.NewPrisoner("MYSTIC_1", "Tartaria", prisoner.ArchetypeMystic)
	// Verify trait
	if !mystic.HasTrait(prisoner.TraitBreatharian) {
		t.Fatalf("Mystic should have Breatharian trait")
	}
	ms.RegisterPrisoner(mystic)

	// Act: Eat Rice
	payload := ResourceIntakePayload{
		PrisonerID: "MYSTIC_1",
		ItemType:   "RICE",
		Amount:     10,
	}
	event := events.GameEvent{
		Type:    events.EventTypeResourceIntake,
		Payload: payload,
	}

	ms.OnResourceIntake(event)

	// Assert: Punishment
	if mystic.Sanity == 100 {
		t.Errorf("Mystic should lose Sanity for eating solids")
	}
	if mystic.HP == 100 {
		t.Errorf("Mystic should lose HP for eating solids")
	}
}
