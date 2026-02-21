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
- [x] **T013:** Implement WebSocket Hub (`internal/network/hub.go`).
- [x] **T014:** Create the `AudienceBridge` (Mobile-to-Server "Pay-to-Torture" API).
- [x] **T015:** Implement the "VAR Replay" endpoint (JSON export of betrayal history).

## Phase 5: The Twin Minds (Agentic AI)
- [x] **T016:** Implement `internal/twins/perception` (EventLog summarizer for LLM context).
- [x] **T017:** Implement `internal/twins/cognition` (MAD-BAD-SAD decision framework).
- [x] **T018:** Implement `internal/twins/action` (System event emitter for punishments/rewards).

## Phase 6: LLM Integration & Shadow Mode
- [x] **T019:** Implement `internal/infra/ai` (Agnostic LLM Provider interface + OpenAI/Anthropic adapters).
- [x] **T020:** Implement Constitutional Prompting (Inject constitution + CoT reasoning JSON format).
- [x] **T021:** Implement Shadow Mode & FinOps Budget Gate (LLMCognitor with server-side MAD validation).

## Phase 7: El Panóptico (Frontend)
- [x] **T022:** Initialize Next.js frontend with dark theme and premium design.
- [x] **T023:** Create Prisoner Dashboard (Real-time vitals, WebSocket connection).
- [x] **T024:** Create Twins Control Panel (Decision history, Shadow Mode toggle).
- [x] **T025:** Create VAR Replay viewer (Event timeline with filtering).

## Phase 8: La Jaula de Hierro (DevOps & CI/CD)
- [x] **T026:** Docker containerization (Multi-stage builds for Go/Next.js).
- [x] **T027:** Docker Compose orchestration (server, client, postgres, redis).
- [x] **T028:** GitHub Actions CI/CD pipeline (tests, lint, build).

## Phase 9: El Motín (Stress Testing & Observability)
- [x] **T029:** Load generator tool (50 concurrent WebSocket clients, 100ms actions).
- [x] **T030:** Observability metrics endpoint (Tick Lag, Event Write Latency).
- [x] **T031:** Concurrency optimization (channel buffers, DB indexes).

## Phase 10: V2.1 Hardcore Mechanics (Audit & New Cast Adaptation)
- [x] **T032:** Implement the Day 21 "Final Dilemma" resolution logic in Ticker and Engine.
- [x] **T033:** Implement Toxic Duo ("Bad Romance") proximity Sanity drain and conflict-to-Hype logic.
- [x] **T034:** Implement Chaos Agent ("Poltergeist") active lockdown noise mechanic (Aída).
- [x] **T035:** Implement Veteran ("Misanthrope") solo Sanity regen logic (Frank).
- [x] **T036:** Implement "Short Fuse" passive logic for Dakota (2x sanity drain on insult, double damage).
- [x] **T037:** Implement "Charm/Deception" mechanic for Héctor (Delay in VAR logs when stealing).
- [x] **T038:** Implement Loyalty Bar active benefits (Sanity regen buffs on high loyalty).
