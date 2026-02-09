"use client";

interface Prisoner {
    id: string;
    name: string;
    archetype: string;
    sanity: number;
    loyalty: number;
    online: boolean;
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

            {/* Warning */}
            {prisoner.sanity < 30 && (
                <div style={{
                    marginTop: "12px",
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
        </div>
    );
}
