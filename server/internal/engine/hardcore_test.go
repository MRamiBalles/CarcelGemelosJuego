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
	el := events.NewEventLog(nil)
	log := logger.NewLogger()
	ss := NewSanitySystem(el, log)

	// Create Actors
	user := prisoner.NewPrisoner("P1", "User", prisoner.ArchetypeVeteran, "CELL_101")
	// cell assignment done by constructor

	witness := prisoner.NewPrisoner("P2", "Witness", prisoner.ArchetypeDeceiver, "CELL_101")
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
	el := events.NewEventLog(nil)
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
	el := events.NewEventLog(nil)
	log := logger.NewLogger()
	ms := NewMetabolismSystem(el, log)

	// Create Mystic
	mystic := prisoner.NewPrisoner("MYSTIC_1", "Tartaria", prisoner.ArchetypeMystic, "CELL_201")
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

func TestShortFuse(t *testing.T) {
	// Setup
	el := events.NewEventLog(nil)
	log := logger.NewLogger()
	ss := NewSanitySystem(el, log)

	// Create Actors
	toxic := prisoner.NewPrisoner("TOX_1", "Dakota", prisoner.ArchetypeExplosive, "CELL_301")
	normal := prisoner.NewPrisoner("NRM_1", "Marco", prisoner.ArchetypeVeteran, "CELL_301")

	// Verify trait integration
	if !toxic.HasTrait(prisoner.TraitShortFuse) {
		t.Fatalf("Toxic archetype should have ShortFuse trait")
	}

	ss.RegisterPrisoner(toxic)
	ss.RegisterPrisoner(normal)

	// Act: Trigger an Insult (Baseline damage modeled as NoiseEventPayload)
	// Intensity 2 * BaseMultiplier(5) = 10 base damage
	payload1 := NoiseEventPayload{
		NoiseType:   NoiseSiren,
		TargetZone:  "TOX_1",
		Intensity:   2,
		DurationSec: 60,
		Reason:      "VERBAL_ABUSE",
	}
	payload2 := NoiseEventPayload{
		NoiseType:   NoiseSiren,
		TargetZone:  "NRM_1",
		Intensity:   2,
		DurationSec: 60,
		Reason:      "VERBAL_ABUSE",
	}

	evt1 := events.GameEvent{Type: events.EventTypeNoiseEvent, Payload: payload1}
	evt2 := events.GameEvent{Type: events.EventTypeNoiseEvent, Payload: payload2}

	ss.OnNoiseEvent(evt1)
	ss.OnNoiseEvent(evt2)

	// Assert: Toxic takes double damage (20 vs 10)
	// Base Sanity is 100.
	expectedToxicSanity := 100 - (10 * 2)
	expectedNormalSanity := 100 - 10

	if toxic.Sanity != expectedToxicSanity {
		t.Errorf("Toxic Sanity mismatch: Expected %v, got %v", expectedToxicSanity, toxic.Sanity)
	}

	if normal.Sanity != expectedNormalSanity {
		t.Errorf("Normal Sanity mismatch: Expected %v, got %v", expectedNormalSanity, normal.Sanity)
	}
}
