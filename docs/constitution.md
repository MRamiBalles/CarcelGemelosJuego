# Constitution of "Cárcel de los Gemelos"

This document defines the immutable laws of the project. Every design and technical decision must adhere to these principles.

## 1. Principle of The "Single Survivor" (The Prisoner's Dilemma)
> **"We are born alone, we die alone."**
*   **The Duo Paradox:** Players compete in pairs (Duos), but **only one person can claim the final cash prize**.
*   **The Judgment Day:** At the end of the 21-day cycle, surviving duos face the "Final Dilemma":
    *   **Betray:** You claim 100% of the pot. Your partner gets 0.
    *   **Collaborate:** You *attempt* to split the pot 50/50. (WARNING: The Twins may declare this illegal or tax it heavily).
    *   **Mutual Betrayal:** Both lose everything. The House wins.

## 2. Principle of Hardcore Chronology
> **"Time is the only true currency."**
*   **The 21-Day Cycle:** The game runs on a rigid, server-authoritative schedule from **March 15th to April 5th**.
*   **Global Clock:** There are no "quick matches" or "pauses". The server clock never stops.
*   **Consequences:** If a player fails to eat by the server-enforced "Dinner Time", they wake up the next day with critical debuffs. Late login = Starvation.

## 3. Principle of Persistence
> **"The state of the game is sacred."**
*   When a player disconnects, their character does not vanish. They remain in the world as a **"Sleeper"**.
*   Sleepers are vulnerable to environmental hazards and player interactions (stealing, voting, or protection).
*   The game world is governed by an **Authoritative Server**. No critical state is stored or calculated solely on the client.

## 4. Principle of Asynchronous Interaction
> **"Life goes on, even when you are away."**
*   All critical social mechanics (Votations, Resource Trades, Trials) must be **Window-Based**.
*   A window of resolution (e.g., 4 to 8 hours) allows players in different time zones or with varying schedules to participate asynchronously.
*   The outcome is projected and resolved only after the window closes.

## 5. Principle of Diegetic Moderation
> **"The Twins are watching, and they are cruel."**
*   Moderation is a gameplay mechanic, not an administrative task.
*   **NLP Verification:** Automated analysis of player interactions checks for extreme toxicity (OOC harassment, bigotry).
*   **Lore Punishment:** Sanctions are applied as in-game events (Sanity debuffs, Water cutoffs, Solitary confinement) triggered by "The Twins" (AI Overlords).
*   Administrative bans are the absolute last resort.

## 6. Principle of Scarcity (The Moral Crumple Zone)
> **"Conflict is a systemic design, not a player flaw."**
*   The game must always provide fewer resources than the total needed for optimal survival.
*   This systemic scarcity is designed to force moral choices, shifting the blame of cruelty from the developers/system to the player decisions.

## 7. Principle of Narrative Sovereignty
> **"The Audience is the Third Twin."**
*   The audience (espectadores) has agency through the "Pay-to-Torture" system.
*   Their interaction must directly impact the environment or resource availability, creating a loop between the reality show players and their observers.
