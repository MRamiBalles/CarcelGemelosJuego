import random
    print("Enhanced Python prototype for 'Cárcel de los Gemelos'.")

class Player:
    def __init__(self, name, archetype):
        self.name = name
        self.archetype = archetype
        self.sanity = 100
        self.loyalty = 50
        self.day = 1
        self.has_consumed = False # Simon's specific check
        self.empathy = random.randint(30, 70) # Hidden stat for betrayal check

    def __repr__(self):
        return f"{self.name} ({self.archetype}) | Sanity: {self.sanity} | Loyalty: {self.loyalty} | Day: {self.day}"

def process_day_logic(player):
    # Simon's 5-Day Cold Turkey
    if player.archetype == "Redeemed":
        if player.day <= 5:
            print(f"> [DEBUFF] {player.name} is in WITHDRAWAL (Day {player.day}). Sanity -5, Stamina low.")
            player.sanity -= 5
        elif not player.has_consumed:
            print(f"> [BUFF] {player.name} attained CLARITY. Sanity +20 and Stat multiplier active.")
            player.sanity += 20

    # Frank's Asocial Trait
    if player.archetype == "Veteran":
        socializing = random.choice([True, False])
        if socializing:
            print(f"> [TRAIT] {player.name} (Veteran) hated the small talk. Sanity -10.")
            player.sanity -= 10
        else:
            print(f"> [TRAIT] {player.name} (Veteran) enjoyed the silence. Sanity +5.")
            player.sanity += 5

def simulate_noise_event(players, event_type="Siren"):
    print(f"\n--- [EVENT] THE TWINS ACTIVATE: {event_type.upper()} TORTURE ---")
    for p in players:
        drain = random.randint(10, 20)
        if p.archetype == "Mystic":
            drain //= 2
            print(f"> {p.name} (Mystic) uses 'Reality Distortion' to doubt the noise. Drain halved.")
        p.sanity -= drain
        print(f"> {p.name} Sanity drained to {p.sanity}.")

def simulate_privacy_drain(p1, p2):
    print("\n--- [MECHANIC] PRIVACY BREACH: TOILET USE ---")
    print(f"{p1.name} is using the open-air toilet.")
    p1.sanity -= 5
    p2.sanity -= 3 # Stress from witness
    print(f"> {p1.name} (User) Dignity loss. {p2.name} (Witness) Stress increase.")

def betrayal_check(player):
    # Cold Blood / Sanity / Empathy Check
    threshold = player.empathy + (100 - player.sanity)
    roll = random.randint(0, 150)
    print(f"> {player.name} Betrayal Roll: {roll} vs Threshold: {threshold}")
    return roll > threshold

def simulate_duo_dilemma(p1, p2):
    print("\n--- [ENDGAME] DAY 21: THE DUO DILEMMA ---")
    p1_betray = betrayal_check(p1)
    p2_betray = betrayal_check(p2)
    
    choice1 = "Betray" if p1_betray else "Cooperate"
    choice2 = "Betray" if p2_betray else "Cooperate"
    
    print(f"{p1.name} chooses to {choice1}")
    print(f"{p2.name} chooses to {choice2}")
    
    if choice1 == "Cooperate" and choice2 == "Cooperate":
        print("RESULT: Both split the reward. VICTORY (Unless Twins prohibit).")
    elif choice1 == "Betray" and choice2 == "Cooperate":
        print(f"RESULT: {p1.name} BETRAYS! {p1.name} wins 100%.")
    elif choice1 == "Cooperate" and choice2 == "Betray":
        print(f"RESULT: {p2.name} BETRAYS! {p2.name} wins 100%.")
    else:
        print("RESULT: BOTH BETRAY. Total rewards lost.")

# Run Simulation
if __name__ == "__main__":
    frank = Player("Frank", "Veteran")
    simon = Player("Simón", "Redeemed")
    
    players = [frank, simon]
    
    for d in range(1, 7): # Simulate first 6 days
        print(f"\n--- DAY {d} ---")
        for p in players:
            p.day = d
            process_day_logic(p)
        
        if d == 3:
            simulate_noise_event(players, "Crying Baby")
            simulate_privacy_drain(frank, simon)

    simulate_duo_dilemma(frank, simon)
