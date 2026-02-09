"use client";

interface HeaderProps {
    gameDay: number;
    tensionLevel: string;
    shadowMode: boolean;
    onToggleShadowMode: () => void;
}

export default function Header({ gameDay, tensionLevel, shadowMode, onToggleShadowMode }: HeaderProps) {
    const tensionColor = {
        LOW: "var(--sanity-green)",
        MEDIUM: "var(--warning-amber)",
        HIGH: "var(--twins-red)",
        CRITICAL: "var(--danger-crimson)",
    }[tensionLevel] || "var(--text-muted)";

    return (
        <header style={{
            background: "var(--bg-cell)",
            borderBottom: "1px solid var(--border-subtle)",
            padding: "16px 24px",
        }}>
            <div className="container" style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between"
            }}>
                {/* Logo & Title */}
                <div style={{ display: "flex", alignItems: "center", gap: "16px" }}>
                    <div className="twins-eye">
                        <span style={{ fontSize: "24px" }}>👁️</span>
                    </div>
                    <div>
                        <h1 style={{ fontSize: "1.5rem", margin: 0 }}>El Panóptico</h1>
                        <p style={{
                            color: "var(--text-muted)",
                            fontSize: "12px",
                            marginTop: "4px"
                        }}>
                            Centro de Control | Cárcel de los Gemelos
                        </p>
                    </div>
                </div>

                {/* Status Indicators */}
                <div style={{ display: "flex", alignItems: "center", gap: "24px" }}>
                    {/* Game Day */}
                    <div style={{ textAlign: "center" }}>
                        <div style={{
                            fontSize: "24px",
                            fontWeight: 700,
                            color: "var(--text-primary)"
                        }}>
                            DÍA {gameDay}
                        </div>
                        <div style={{ fontSize: "11px", color: "var(--text-muted)" }}>
                            de 21
                        </div>
                    </div>

                    {/* Tension Level */}
                    <div style={{
                        padding: "8px 16px",
                        background: "var(--bg-surface)",
                        borderRadius: "8px",
                        borderLeft: `3px solid ${tensionColor}`,
                    }}>
                        <div style={{ fontSize: "11px", color: "var(--text-muted)" }}>TENSIÓN</div>
                        <div style={{ fontWeight: 600, color: tensionColor }}>{tensionLevel}</div>
                    </div>

                    {/* Shadow Mode Toggle */}
                    <button
                        onClick={onToggleShadowMode}
                        className={shadowMode ? "shadow-mode-badge" : ""}
                        style={{
                            padding: "8px 16px",
                            background: shadowMode ? "rgba(245, 158, 11, 0.15)" : "var(--twins-red)",
                            border: shadowMode ? "1px solid var(--warning-amber)" : "none",
                            borderRadius: "20px",
                            color: shadowMode ? "var(--warning-amber)" : "white",
                            fontSize: "12px",
                            fontWeight: 500,
                            cursor: "pointer",
                            display: "flex",
                            alignItems: "center",
                            gap: "6px",
                        }}
                    >
                        {shadowMode ? "🌑 SHADOW MODE" : "☀️ LIVE MODE"}
                    </button>
                </div>
            </div>
        </header>
    );
}
