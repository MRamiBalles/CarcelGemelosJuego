# Phase 18 Tasks: Reality Show Mechanics

- [x] **T063:** Add `IsIsolated` boolean to `Prisoner`, `PrisonerSnapshot`, and SQLite schema in `sqlite_repo.go`.
- [x] **T064:** Create `isolation_system.go` to handle the 24h solitary confinement rules, specifically buffing Frank and debuffing Toxic archetypes.
- [x] **T065:** Implement `api/poll/create` and `api/poll/vote` endpoints in `router.go` to simulate Twitch audience interventions, generating a `PollResolvedEvent` to reward/punish players.
- [x] **T066:** Build `PollWidget.tsx` in the Next.js frontend to visualize ongoing audience polls.
- [ ] **T067:** Create `patio_system.go` to trigger daily risk-reward fitness tests (trading Stamina for Pot).
- [ ] **T068:** Create `contraband_system.go` to manage hidden loot and the `ActionSnitch` betrayal mechanic.
- [ ] **T069:** Write unit/integration tests for Isolation and Snitching logic.
