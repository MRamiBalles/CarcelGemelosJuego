# Technical Task List: "Cárcel de los Gemelos" (T-Series)

This document breaks down the implementation into atomic, testable units to avoid God Components and Architectural Smells.

## Phase 0: Infrastructure & Project Scaffolding
- [x] **T001:** Setup Go Workspace Structure (`cmd/`, `internal/domain`, `internal/infra`).
- [x] **T002:** Initialize `go.mod` and add core dependencies (Websockets, Loggers).
- [x] **T003:** Implement `internal/platform/logger` (Structured logging for tracing "The Twins" actions).

## Phase 1: Domain Entities & Game State
- [x] **T004:** Define `internal/domain/prisoner` structs (Sanity, Loyalty, Archetype traits).
- [x] **T005:** Implement `internal/domain/events` (Immutable Event log for the "VAR of Betrayal").
- [x] **T006:** Define `internal/domain/rules` (Pure functions for Hunger/Sanity/Noise calculations).

## Phase 2: Core Game Logic (The Twins Engine)
- [x] **T007:** Implement `internal/engine/ticker` (The server heartbeat for resource depletion).
- [x] **T008:** Implement `internal/engine/noise` (Managing the timing and Sanity drain of audio tortures).
- [x] **T009:** Implement `internal/engine/sanity_system` (Subscriber to NoiseEvents, emits SanityChangeEvents).

## Phase 3: Persistence & Event Sourcing
- [x] **T010:** Design PostgreSQL schema for append-only event log (Immutable Ledger).
- [x] **T011:** Implement `EventRepository` in `internal/infra/storage` (Domain-agnostic interface).
- [x] **T012:** Implement Redis cache + Reality Recap reconstructor.

## Phase 4: Networking & External Bridges
- [/] **T013:** Implement WebSocket Server (`internal/network/ws`).
- [/] **T014:** Create the `AudienceBridge` (Mobile-to-Server API for Sadism Points).
- [/] **T015:** Implement the "VAR Replay" endpoint (JSON export of betrayal history).
