# Implementation Plan - Phase 18 (Reality Show Mechanics)

Since the C: drive is full, we are tracking Phase 18 directly in the D: drive.
This phase aims to adapt the hardcore survival mechanics to a Twitch-friendly reality show format.

## Goal Description
Implement specific Meta-mechanics (Isolation, Audience Polling, Patio Exams, Contraband) to make the psychological survival aspect heavily influenced by "El Gran Hermano" style audience interaction and asymmetric warfare.

## Proposed Changes

### 1. Isolation Module (`internal/engine/isolation_system.go`)
- **Mechanic:** The "Celda de Castigo". A 24h debuff/buff room.
- **Prisoner State:** Add `IsIsolated bool` to `Prisoner` struct and `PrisonerSnapshot`.
- **Logic:** 
  - Isolated prisoners cannot interact with others (Social queries skip them).
  - Frank (Misanthrope) regenerates double Sanity here.
  - Toxic (Ylenia/Labrador) lose "Hype" generation capabilities and take heavy Sanity damage.

### 2. Audience Polling (`internal/network/router.go` & `frontend`)
- **Backend API:**
  - `POST /api/poll/create`: Admin endpoint to start a poll (e.g. "Who gets Sushi vs Torture?").
  - `POST /api/poll/vote`: For simulated Twitch chat or frontend users to cast votes.
  - Generates `PollResolvedEvent` which automatically distributes the reward (Stamina/Sanity boost) and punishment (Audio/Sanity drop).
- **Frontend Widget:** Build a React component (`PollWidget.tsx`) to show real-time voting percentages.

### 3. Patio Challenges (`internal/engine/patio_system.go`)
- **Mechanic:** Daily high-risk event at 12:00 PM.
- **Logic:** CronJob triggers `PatioEvent`. A duo must choose one member to participate.
  - Active participant loses 80% Stamina (high risk of death if they were already starving).
  - Reward: Grants massive amount of `PotContribution` (Money) or resets Hunger to 100 for both.

### 4. Contraband & Snitching (`internal/engine/contraband_system.go`)
- **Mechanic:** Hidden random loot items and a new action `ActionSnitch`.
- **System:**
  - `LootEvent` randomly assigns "Contraband" (Cigarettes, Phone) to a cell.
  - Using it gives a massive Sanity buff but lowers Loyalty.
  - Any prisoner in line-of-sight can call `ActionSnitch` on the target.
  - If target had contraband: Snitch gains 50% of target's Pot. Target is sent to Isolation.
  - If target was clean: Snitch takes massive Sanity damage for lying to "Los Gemelos".

## Verification Plan
- Unit tests for Isolation logic (checking Frank's buff vs Toxic's debuff).
- Manual test of `POST /api/poll/create` and checking frontend React state updates.
