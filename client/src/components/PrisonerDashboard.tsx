"use client";

interface Prisoner {
    id: string;
    name: string;
    archetype: string;
    sanity: number;
    loyalty: number;
    dignity: number;
    pot_contribution: number;
    online: boolean;
    stamina: number;
    hunger: number;
    thirst: number;
    states: Record<string, any>;
    inventory: any[];
}

interface PrisonerDashboardProps {
    prisoners: Prisoner[];
}

export default function PrisonerDashboard({ prisoners }: PrisonerDashboardProps) {
    const archetypeEmoji: Record<string, string> = {
        VETERAN: "🎖️",
        MYSTIC: "🔮",
        SHOWMAN: "🎭",
        REDEEMED: "✝️",
    };

    return (
        <div>
            <div style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                marginBottom: "24px"
            }}>
                <h2 style={{ fontSize: "1.25rem", fontWeight: 600 }}>
                    Estado de Prisioneros
                </h2>
                <div style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "8px",
                    color: "var(--text-secondary)",
                    fontSize: "14px"
                }}>
                    <span className="status-dot online"></span>
                    {prisoners.filter(p => p.online).length} en línea
                </div>
            </div>

            <div style={{
                display: "grid",
                gridTemplateColumns: "repeat(auto-fill, minmax(320px, 1fr))",
                gap: "16px"
            }}>
                {prisoners.map((prisoner) => (
                    <PrisonerCard key={prisoner.id} prisoner={prisoner} emoji={archetypeEmoji[prisoner.archetype] || "👤"} />
                ))}
            </div>
        </div>
    );
}

function PrisonerCard({ prisoner, emoji }: { prisoner: Prisoner; emoji: string }) {
    const sanityColor = prisoner.sanity < 30 ? "var(--twins-red)" :
        prisoner.sanity < 60 ? "var(--warning-amber)" :
            "var(--sanity-green)";

    const loyaltyColor = prisoner.loyalty < 30 ? "var(--danger-crimson)" :
        prisoner.loyalty < 60 ? "var(--warning-amber)" :
            "var(--loyalty-blue)";

    const dignityColor = prisoner.dignity < 30 ? "var(--twins-red)" :
        prisoner.dignity < 60 ? "var(--warning-amber)" :
            "var(--sanity-green)";

    const staminaColor = prisoner.stamina < 30 ? "var(--twins-red)" :
        prisoner.stamina < 60 ? "var(--warning-amber)" :
            "var(--loyalty-blue)";

    // Determine warnings based on archetype traits
    let traitWarning = null;
    if (prisoner.archetype === "VETERAN" && prisoner.sanity < 40) {
        traitWarning = "⚠️ Aislar para regenerar (Misántropo)";
    } else if (prisoner.archetype === "DECEIVER" && prisoner.sanity < 30) {
        traitWarning = "⚠️ Peligro de explosión corta (Short Fuse)";
    } else if (prisoner.archetype === "TOXIC" && prisoner.sanity < 30) {
        traitWarning = "⚠️ Relación Tóxica inminente (Bad Romance)";
    }

    // Convert states object keys to an array for easy rendering
    const activeStates = prisoner.states ? Object.keys(prisoner.states) : [];

    return (
        <div className="card" style={{ position: "relative" }}>
            {/* Online Indicator */}
            <div
                className={`status-dot ${prisoner.online ? "online" : "offline"}`}
                style={{ position: "absolute", top: "16px", right: "16px" }}
            ></div>

            {/* Header */}
            <div style={{ display: "flex", alignItems: "center", gap: "12px", marginBottom: "16px" }}>
                <div style={{
                    width: "48px",
                    height: "48px",
                    borderRadius: "50%",
                    background: "var(--bg-surface)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "24px",
                }}>
                    {emoji}
                </div>
                <div>
                    <div style={{ fontWeight: 600, fontSize: "16px" }}>{prisoner.name}</div>
                    <div style={{
                        fontSize: "12px",
                        color: "var(--text-muted)",
                        textTransform: "uppercase",
                        letterSpacing: "0.05em"
                    }}>
                        {prisoner.archetype}
                    </div>
                </div>
            </div>

            {/* Sanity Bar */}
            <div style={{ marginBottom: "12px" }}>
                <div style={{
                    display: "flex",
                    justifyContent: "space-between",
                    marginBottom: "4px",
                    fontSize: "12px"
                }}>
                    <span style={{ color: "var(--text-muted)" }}>Cordura</span>
                    <span style={{ color: sanityColor, fontWeight: 600 }}>{prisoner.sanity}%</span>
                </div>
                <div className="stat-bar">
                    <div
                        className="stat-bar-fill"
                        style={{
                            width: `${prisoner.sanity}%`,
                            background: sanityColor,
                        }}
                    ></div>
                </div>
            </div>

            {/* Loyalty and Pot Container */}
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "12px", marginBottom: "12px" }}>
                {/* Loyalty Bar */}
                <div>
                    <div style={{
                        display: "flex",
                        justifyContent: "space-between",
                        marginBottom: "4px",
                        fontSize: "12px"
                    }}>
                        <span style={{ color: "var(--text-muted)" }}>Lealtad</span>
                        <span style={{ color: loyaltyColor, fontWeight: 600 }}>{prisoner.loyalty}%</span>
                    </div>
                    <div className="stat-bar">
                        <div
                            className="stat-bar-fill"
                            style={{
                                width: `${prisoner.loyalty}%`,
                                background: loyaltyColor,
                            }}
                        ></div>
                    </div>
                </div>

                {/* Pot/Money Contribution */}
                <div>
                    <div style={{
                        display: "flex",
                        justifyContent: "space-between",
                        marginBottom: "4px",
                        fontSize: "12px"
                    }}>
                        <span style={{ color: "var(--text-muted)" }}>Bote (Hype)</span>
                        <span style={{ color: "var(--warning-amber)", fontWeight: 600 }}>${prisoner.pot_contribution.toFixed(2)}</span>
                    </div>
                    <div style={{
                        background: "rgba(245, 158, 11, 0.1)",
                        height: "8px",
                        borderRadius: "4px",
                        border: "1px solid rgba(245, 158, 11, 0.3)",
                        display: "flex",
                        alignItems: "center",
                        padding: "0 4px"
                    }}>
                        {/* Static visual representation for money */}
                        <div style={{ width: "100%", height: "2px", background: "var(--warning-amber)", opacity: 0.5 }}></div>
                    </div>
                </div>
            </div>

            {/* Dignity/Toilet Tracker */}
            <div>
                <div style={{
                    display: "flex",
                    justifyContent: "space-between",
                    marginBottom: "4px",
                    fontSize: "12px"
                }}>
                    <span style={{ color: "var(--text-muted)" }}>Dignidad</span>
                    <span style={{ color: dignityColor, fontWeight: 600 }}>{prisoner.dignity}%</span>
                </div>
                <div className="stat-bar">
                    <div
                        className="stat-bar-fill"
                        style={{
                            width: `${prisoner.dignity}%`,
                            background: dignityColor,
                        }}
                    ></div>
                </div>
            </div>

            {/* Stamina / Thirst / Hunger Layer (F4/F5) */}
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: "8px", marginBottom: "12px" }}>
                <div>
                    <div style={{ fontSize: "10px", color: "var(--text-muted)", marginBottom: "4px" }}>Stamina</div>
                    <div className="stat-bar" style={{ height: "4px" }}>
                        <div className="stat-bar-fill" style={{ width: `${prisoner.stamina}%`, background: staminaColor }}></div>
                    </div>
                </div>
                <div>
                    <div style={{ fontSize: "10px", color: "var(--text-muted)", marginBottom: "4px" }}>Thirst</div>
                    <div className="stat-bar" style={{ height: "4px" }}>
                        <div className="stat-bar-fill" style={{ width: `${prisoner.thirst}%`, background: "var(--loyalty-blue)" }}></div>
                    </div>
                </div>
                <div>
                    <div style={{ fontSize: "10px", color: "var(--text-muted)", marginBottom: "4px" }}>Hunger</div>
                    <div className="stat-bar" style={{ height: "4px" }}>
                        <div className="stat-bar-fill" style={{ width: `${prisoner.hunger}%`, background: "var(--warning-amber)" }}></div>
                    </div>
                </div>
            </div>

            {/* Inventory Box (F2) */}
            {prisoner.inventory && prisoner.inventory.length > 0 && (
                <div style={{ marginBottom: "12px", background: "var(--bg-card)", padding: "8px", borderRadius: "6px", display: "flex", flexWrap: "wrap", gap: "4px" }}>
                    <div style={{ fontSize: "10px", color: "var(--text-muted)", width: "100%", marginBottom: "4px" }}>Inventario Físico</div>
                    {prisoner.inventory.map((item, idx) => (
                        <div key={idx} style={{
                            fontSize: "12px",
                            background: "var(--bg-surface)",
                            padding: "2px 6px",
                            borderRadius: "4px",
                            border: item.is_contraband ? "1px solid var(--twins-red)" : "1px solid var(--border-light)"
                        }} title={String(item.type)}>
                            {item.quantity}x {item.type === "RICE" ? "🍚" : item.type === "WATER" ? "💧" : item.type === "ELIXIR" ? "🧪" : item.type === "SUSHI" ? "🍣" : "📦"}
                        </div>
                    ))}
                </div>
            )}

            {/* F4/F5 States Container */}
            {activeStates.length > 0 && (
                <div style={{ display: "flex", flexWrap: "wrap", gap: "4px", marginBottom: "12px" }}>
                    {activeStates.map(state => (
                        <span key={state} style={{
                            fontSize: "10px",
                            padding: "2px 6px",
                            borderRadius: "12px",
                            background: state === "Dead" ? "var(--twins-red)" : state === "Exhausted" ? "var(--warning-amber)" : "var(--bg-surface)",
                            color: state === "Dead" || state === "Exhausted" ? "white" : "var(--text-secondary)",
                            fontWeight: state === "Dead" ? "bold" : "normal"
                        }}>
                            {state === "Asleep" ? "💤 Durmiendo" : state === "Exhausted" ? "😰 Exhausto" : state === "Isolated" ? "🧊 Aislado" : state === "Dead" ? "💀 Evacuado" : state === "Meditating" ? "🧘 Meditando" : state}
                        </span>
                    ))}
                </div>
            )}

            {/* Warnings container */}
            <div style={{ marginTop: "12px", display: "flex", flexDirection: "column", gap: "6px" }}>
                {prisoner.sanity < 30 && (
                    <div style={{
                        padding: "8px",
                        background: "rgba(220, 38, 38, 0.1)",
                        borderRadius: "6px",
                        fontSize: "12px",
                        color: "var(--twins-red)",
                        display: "flex",
                        alignItems: "center",
                        gap: "6px"
                    }}>
                        ⚠️ Cordura crítica - MAD protegido
                    </div>
                )}
                {traitWarning && (
                    <div style={{
                        padding: "8px",
                        background: "rgba(245, 158, 11, 0.1)",
                        borderRadius: "6px",
                        fontSize: "12px",
                        color: "var(--warning-amber)",
                        display: "flex",
                        alignItems: "center",
                        gap: "6px"
                    }}>
                        {traitWarning}
                    </div>
                )}
            </div>
        </div>
    );
}
