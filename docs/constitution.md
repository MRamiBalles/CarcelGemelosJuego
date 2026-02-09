# Constitution of "Cárcel de los Gemelos"

This document defines the immutable laws of the project. Every design and technical decision must adhere to these principles.

## 1. Principle of Persistence
> **"The state of the game is sacred."**
*   When a player disconnects, their character does not vanish. They remain in the world as a **"Sleeper"**.
*   Sleepers are vulnerable to environmental hazards and player interactions (stealing, voting, or protection).
*   The game world is governed by an **Authoritative Server**. No critical state is stored or calculated solely on the client.

## 2. Principle of Asynchronous Interaction
> **"Life goes on, even when you are away."**
*   All critical social mechanics (Votations, Resource Trades, Trials) must be **Window-Based**.
*   A window of resolution (e.g., 4 to 8 hours) allows players in different time zones or with varying schedules to participate asynchronously.
*   The outcome is projected and resolved only after the window closes.

## 3. Principle of Diegetic Moderation
> **"The Twins are watching, and they are cruel."**
*   Moderation is a gameplay mechanic, not an administrative task.
*   **NLP Verification:** Automated analysis of player interactions checks for extreme toxicity (OOC harassment, bigotry).
*   **Lore Punishment:** Sanctions are applied as in-game events (Sanity debuffs, Water cutoffs, Solitary confinement) triggered by "The Twins" (AI Overlords).
*   Administrative bans are the absolute last resort; most friction is resolved through behavioral mechanics and social debt.

## 4. Principle of Scarcity (The Moral Crumple Zone)
> **"Conflict is a systemic design, not a player flaw."**
*   The game must always provide fewer resources than the total needed for optimal survival.
*   This systemic scarcity is designed to force moral choices, shifting the blame of cruelty from the developers/system to the player decisions.

## 5. Principle of Narrative Sovereignty
> **"The Audience is the Third Twin."**
*   The audience (espectadores) has agency through the "Pay-to-Torture" system.
*   Their interaction must directly impact the environment or resource availability, creating a loop between the reality show players and their observers.
