# Technical Plan: "Cárcel de los Gemelos" (V1.0)

This plan defines the "How," focusing on mitigating **Architectural Technical Debt (ATD)**.

## 1. System Architecture: Authoritative Hybrid

To support the 1-week persistent reality show, we use a decentralized-authoritative model.

*   **Logic Master:** Authoritative Game Server written in **Go**.
*   **Architecture Pattern:** **Clean Architecture** (Separation of Concerns).
    - `/cmd`: Entry points for the jail-server.
    - `/internal/domain`: Pure logic (Prisoner stats, Sanity rules).
    - `/internal/engine`: Simulation ticker and event processing.
    - `/internal/events`: Event Sourcing (The VAR).
    - `/internal/network`: Agnostic WebSocket management.
*   **State Persistence:** **Event Sourcing** using a combination of **PostgreSQL** and **Redis**.
    -   **The "VAR of Betrayal":** Every interaction (resource intake, chat, skill use) is an immutable event with a timestamp. This allows the system to generate "Replays" (JSON/Video-lite) to expose lies to the audience or cellmates.
*   **Client Communication:** **WebSockets** with binary serialization (Protobuf) for movement, and a separate **Pub/Sub** channel for environmental "Noise Events".

## 2. Modular Components (Anti-God Object)

The server is divided into independent micro-services/modules to stay under the ATDI limit:

1.  **`IdentityService`:** Persistent storage of player Vitals (Hunger, Thirst) and **Sanity (Event-Sourced)**.
2.  **`SocialEngine`:** Manages "Social Debt", **Duo Loyalty Bars**, and the final **Duo Dilemma** resolution logic.
3.  **`EnvironmentService` (The Twins):** Manages "Noise Events", temperature shifts, and resource scarcity triggers. Handles NLP moderation.
4.  **`AudienceBridge`:** API for the Mobile App to trigger `EnvironmentService` events via "Sadism Points".

## 3. Data Sync Strategy (Option C)

| Layer | Technology | Frequency |
| :--- | :--- | :--- |
| **PC Client (Unreal/Unity)** | WebSockets / TCP | High (Real-time movement/voice) |
| **Mobile App (Audience)** | HTTP / SSE | Low (Voting / Event updates) |
| **Recap Engine** | Event Store | On-Login (Build state since last session) |

## 4. Mitigation of Architectural Technical Debt (ATD)

*   **Anti-Vibe Moderation:** Rules are hard-coded in the `EnvironmentService`. Consequences (Sanity drain, expulsion) are triggered by system thresholds, not human GM whims.
*   **State-Reconstruction:** Using the Event Store to recalculate Sanity/Loyalty precisely upon login. Disconnected players ("Sleepers") continue to process environmental events (Noise) via the server.
*   **Privacy-Link State:** A dedicated `VisibilityManager` tracks "Line of Sight" between cellmates and the toilet object, triggering Sanity/Stress events only when active gazing occurs.
