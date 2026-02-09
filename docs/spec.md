# Game Specification: "Cárcel de los Gemelos" (V2.1)

This spec defines the "What" and the "Why" of the game mechanics, serving as the live GDD under SDD protocol.

## 1. User Journey: Day 1 (The Intake)

1.  **Selection (Out-of-Game):** Player selects an Archetype.
2.  **Awakening:** Cinematic of waking up in a cramped, dark cell. The door is locked. A screen displays "Subject ID" and a 21-day countdown.
3.  **Role Assignment:** "The Twins" (voice/visuals) announce the day's first objective: "Earn your calories."
4.  **First Friction:** The cell door opens. To get water, you must interact with a neighbor. The water dispenser only gives enough for 2 people, but there are 4 in the block.
5.  **First Bond:** The Veteran must use "Survival Extreme" to show others how to drink from a leaky pipe if they fail the water trial, establishing dependency.

## 2. Comprehensive Class Mechanics

| Class | Primary Goal | Active Input | Passive Output |
| :--- | :--- | :--- | :--- |
| **Veteran (Frank)** | Sanctuary Funding | **Extreme Survival:** 24h immunity to tainted water/food. | **Asocial:** Sanity drains faster when socializing; regens when alone or doing physical tasks. |
| **Mystic (Tartaria)** | Reality Distortion | **Reality Debate:** Lower an opponent's Sanity in a debate mini-game. | **Distortion Aura:** Regen Sanity to allies; if Sanity < 20%, allies suffer Confusion. |
| **Showman (Marrash)** | Content/Attention | **Smoke Screen:** Hide cell activity from cameras for 15m. | **Attention Vampire:** Stats buffs for generating "Viral" conflict/attention. |
| **Redeemed (Simón)** | Survival/Expiation | **Cold Turkey (Day 1-5):** Surviving lock-in without "Consuming" grants massive Clarity buffs. | **Withdrawal (Day 1-5):** Blurred vision, shaking aim, reduced stamina. |

## 3. The "Anti-Goals" (Scope Control)

To ensure delivery of V1, the following features are **OUT OF SCOPE**:
*   **Complex Melee Combat:** No block/parry/stamina-based duel system. Combat is resolved via quick-time events or stealth traps.
*   **Procedural World Generation:** The prison layout is fixed to maximize social engineering and "bottleneck" encounters.
*   **Open World Exploration:** The game is restricted to the Ship/Prison. No "landing on planets" in V1.
*   **Player-Driven Construction:** No base building. You modify your cell (scrappy furniture), but you cannot build structures.

## 4. The 15-Minute Loop (Operational Detail)

*   **Minute 0-2 (Triage):** Auto-check of HUD (persistent stats). Read the "Reality Recap" (Events that happened while offline).
*   **Minute 2-10 (Foraging/Social):** Move to common areas. Use archetype skill. Socialize to hide true status.
## 5. The "Sanity & Noise" System

*   **Sanity Stats:** A fluctuating bar (0-100). 0 Sanity causes "Breakdown" (Loss of character control/forced actions).
*   **The "Twins" Noise Events:** Audio tortures (sirens, crying babies, scratching) triggered by Admins/Audience.
    *   *Effect:* Rapid Sanity drain. 
    *   *Counter:* Mystic's Aura or Veteran's "Solitude Focus".
*   **Zero Privacy (The Visible Toilet):** Toilets are located inside cells and are fully visible.
    *   *Debuff:* Using the toilet drains "Dignity/Sanity" from the user.
    *   *Gaze Stress:* If the cell partner is online/looking, they suffer "Stress" drain.
*   **The Loyalty Bar (Duo System):**
    *   Tracks the "Technical Trust" between cellmates.
    *   Higher Loyalty = Faster resource sharing & Sanity regen bonuses.
    *   Lower Loyalty = Higher chance of "Traition Criticals" (bonus for stealing from partner).

## 6. The "Duo Dilemma" (Endgame)

*   **Rule of Gold:** "Joint effort, Single payout."
*   **Day 21 Decision:** If a duo reaches the end, the system forces a prisoners' dilemma.
    *   *Cold Blood Check:** The "Betray" option requires a server-side check. If Sanity is too low or "Empathy" (hidden stat) is too high, the action might fail or trigger a mini-game.
    *   *Both Cooperate:* Split reward (50% each - *Warning: May be prohibited by The Twins*).
    *   *One Betrays:* Betrayer gets 100%, Victim gets 0%.
    *   *Both Betray:* Both lose everything. The Twins win.
