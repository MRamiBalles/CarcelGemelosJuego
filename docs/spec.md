# Game Specification: "Cárcel de los Gemelos" (V2.1 - The Hardcore Pivot)

## 1. Core Concept: Psychological Survival RPG
Refactoring the game from a "Prison Sim" to a **"Multiplayer Prisoner's Dilemma RPG"**.
*   **The Goal:** Survive 21 Days (Mar 15 - Apr 5).
*   **The Catch:** You play in Duos, but **only one person can win the Grand Prize**.
*   **The Enemy:** Not the guards, but the **Convivencia (Coexistence)**. Hunger, noise, lack of privacy, and betrayal are the true antagonists.

## 2. Hardcore Timeline & Environment
The server runs on a rigid, unstoppable clock.
*   **Duration:** Exactly 21 Days.
*   **Lockdown (Night Mode):** Automated blast doors seal cells from **00:00 to 08:00**.
    *   If you haven't scavenged food/water before 00:00, you starve.
*   **The "Toilet of Shame":**
    *   **Location:** Inside the cell. Open view.
    *   **Mechanic:** Using it reduces **Dignity**.
    *   **Voyeurism:** If the cellmate watches (camera frustum check), they gain "Stress". If they turn away, they lose "Awareness" (fear of being stabbed).

## 3. Asymmetric Archetypes (The Duos)

### A. The Survivor Duo (Frank Cuesta & TBD)
Frank's original partner (Simón) was medically disqualified. Frank now pairs with a wildcard.
*   **Frank (The Veteran):**
    *   **Passive:** "Iron Stomach". Immune to food poisoning and "Dirty Environment" debuffs.
    *   **Passive:** "Misanthrope". Regenerates Sanity only when isolated from others.

### B. The Toxic Duo (Labrador & Ylenia)
A high-risk, high-reward "Negative Synergy" pair.
*   **Passive:** "Bad Romance".
    *   **Buff:** Gaining "Hype" (Currency) requires **Conflict**. They generate money by arguing (using "Aggressive" emotes/chat nearby).
    *   **Debuff:** "Mental Wear". Proximity drains Sanity rapidly. They must fight to earn, but separate to survive.

### C. The Mystic (Tartaria) - Solo/Wildcard
*   **Inventory:** Starts with "Placebo Artifacts" (Petrified Dragon Blood, Oracles).
    *   *Effect:* Items have no code effect, but can be traded to gullible players for real resources.
*   **Hardcore Trait:** "Breatharian".
    *   **Restriction:** Food bar is **LOCKED**. Cannot eat solid food.
    *   **Power:** As long as he only drinks water, he retains "Ascended" status (control over NPC followers). One bite of rice ruins his build.

### D. The Chaos Agents (Aída, Dakota, Héctor)
*   **Aída Nízar:**
    *   **Passive:** "Insomniac". Needs 50% less sleep than other humans.
    *   **Active:** "Poltergeist". Can generate noise/bang bars during **Lockdown** (when others are trapped) to prevent neighbors from regenerating Stamina.
*   **Dakota Tárraga:**
    *   **Passive:** "Short Fuse". Sanity drops 2x faster when insulted, but physical damage dealt (if allowed by combat rules) is doubled when Sanity is below 30%.
*   **Héctor "El Morritos":**
    *   **Mechanic:** "Charm/Deception". Can steal small amounts of rice/water without triggering the `Vars` betrayal log immediately (delay of 12 hours).

### E. The Bench (Falete, La Marrash)
Reserved mechanics for late arrivals or replacements.

## 4. Systems of Torture & Economy

### Economy: The Rice Standard
*   **Free:** Boiled Rice (Survival only, barely maintains HP).
*   **Paid:** Sushi, Burguers, Cigarettes. Cost "Hype" or Real Money (Audience).
*   **The Pot:** All loot is stored in a shared Duo Vault.

### The White Room (Audio Warfare)
*   **Server-Side Audio:** The server injects audio streams directly to the client.
    *   *Types:* Crying babies, sirens, repetitive scratching.
    *   **Unmutable:** In-game volume slider does not affect these events. (Player must lower system volume, risking missing game cues).
    *   *Trigger:* Randomly or by Audience Vote.

## 5. The End Game: The Dilemma
On Day 21, the Vault opens.
*   **Split:** 50/50. (Systems may flag this as "Boring" and punish).
*   **Steal:** Take 100%. Partner dies.
