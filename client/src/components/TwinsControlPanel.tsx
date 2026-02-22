"use client";

interface Decision {
    id: string;
    timestamp: string;
    action: string;
    target: string;
    approved: boolean;
    shadow: boolean;
    madViolation?: string;
}

interface TwinsControlPanelProps {
    decisions: Decision[];
    shadowMode: boolean;
    isLockdown?: boolean;
    day21Dilemmas?: { prisoner: string; decision: 'BETRAY' | 'COLLABORATE' | 'PENDING' }[];
    onToggleShadowMode: () => void;
    onTriggerOracle?: (target: string, message: string) => void;
    onTriggerTorture?: (soundId: string) => void;
}

export default function TwinsControlPanel({ decisions, shadowMode, isLockdown = false, day21Dilemmas = [], onToggleShadowMode, onTriggerOracle, onTriggerTorture }: TwinsControlPanelProps) {
    const actionEmoji: Record<string, string> = {
        NOISE_TORTURE: "🔊",
        REVEAL_SECRET: "👁️",
        RESOURCE_CUT: "🚰",
        REWARD: "🎁",
        OBSERVE: "🔍",
    };

    return (
        <div>
            {/* Twins Status */}
            <div className="card glow-red" style={{ marginBottom: "24px" }}>
                <div style={{ display: "flex", alignItems: "center", gap: "20px" }}>
                    <div className="twins-eye">
                        <span style={{ fontSize: "28px" }}>👁️</span>
                    </div>
                    <div style={{ flex: 1 }}>
                        <h2 style={{ fontSize: "1.25rem", fontWeight: 600, marginBottom: "8px" }}>
                            Los Gemelos
                        </h2>
                        <p style={{ color: "var(--text-secondary)", fontSize: "14px" }}>
                            Sistema cognitivo autónomo • MAD-BAD-SAD Framework
                        </p>
                    </div>
                    <div style={{ textAlign: "right" }}>
                        <div style={{ fontSize: "11px", color: "var(--text-muted)", marginBottom: "4px" }}>
                            MODO ACTUAL
                        </div>
                        <button
                            onClick={onToggleShadowMode}
                            style={{
                                padding: "12px 24px",
                                background: shadowMode ? "var(--warning-amber)" : "var(--twins-red)",
                                border: "none",
                                borderRadius: "8px",
                                color: "white",
                                fontWeight: 600,
                                cursor: "pointer",
                                fontSize: "14px",
                            }}
                        >
                            {shadowMode ? "🌑 SHADOW" : "☀️ EN VIVO"}
                        </button>
                    </div>
                </div>
            </div>

            {/* System Status Alerts */}
            <div style={{ display: "flex", gap: "16px", marginBottom: "24px", flexWrap: "wrap" }}>
                {/* Lockdown Indicator */}
                <div className="card" style={{ flex: 1, minWidth: "250px", display: "flex", alignItems: "center", gap: "12px", borderLeft: isLockdown ? "4px solid var(--twins-red)" : "4px solid var(--sanity-green)" }}>
                    <div style={{ fontSize: "24px" }}>{isLockdown ? "🔒" : "🔓"}</div>
                    <div>
                        <div style={{ fontWeight: 600, fontSize: "14px", color: isLockdown ? "var(--twins-red)" : "var(--sanity-green)" }}>
                            {isLockdown ? "LOCKDOWN ACTIVO" : "Celdas Abiertas"}
                        </div>
                        <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>
                            {isLockdown ? "00:00 - 08:00 (Aislamiento)" : "Horario de Convivencia"}
                        </div>
                    </div>
                </div>

                {/* Day 21 Dilemma */}
                {day21Dilemmas.length > 0 && (
                    <div className="card" style={{ flex: 2, minWidth: "300px", borderLeft: "4px solid var(--warning-amber)" }}>
                        <div style={{ fontWeight: 600, fontSize: "14px", color: "var(--warning-amber)", marginBottom: "8px" }}>
                            ⚖️ Dilema Final (Día 21)
                        </div>
                        <div style={{ display: "flex", gap: "16px" }}>
                            {day21Dilemmas.map((d, i) => (
                                <div key={i} style={{ flex: 1, background: "var(--bg-surface)", padding: "8px", borderRadius: "6px" }}>
                                    <div style={{ fontSize: "12px", color: "var(--text-muted)", marginBottom: "4px" }}>{d.prisoner}</div>
                                    <div style={{
                                        fontSize: "13px",
                                        fontWeight: 600,
                                        color: d.decision === 'BETRAY' ? "var(--twins-red)" :
                                            d.decision === 'COLLABORATE' ? "var(--sanity-green)" : "var(--text-secondary)"
                                    }}>
                                        {d.decision === 'PENDING' ? "⏳ ESPERANDO..." : d.decision === 'BETRAY' ? "🗡️ TRAICIÓN" : "🤝 COLABORACIÓN"}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </div>
            {/* Manual Overrides (Testing) */}
            <div className="card" style={{ marginBottom: "24px", borderLeft: "4px solid var(--text-primary)" }}>
                <div style={{ fontWeight: 600, fontSize: "14px", marginBottom: "16px" }}>
                    🛠️ Overrides Manuales (Pruebas de API)
                </div>
                <div style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}>
                    <button
                        onClick={() => onTriggerOracle?.("P001", "La audiencia te observa, Simón.")}
                        style={{
                            padding: "8px 16px",
                            background: "var(--bg-surface)",
                            border: "1px solid var(--border-subtle)",
                            borderRadius: "6px",
                            color: "var(--text-primary)",
                            cursor: "pointer",
                            fontSize: "13px"
                        }}
                    >
                        👁️ Enviar Oráculo (P001)
                    </button>

                    <button
                        onClick={() => onTriggerTorture?.("SCREAM_01")}
                        style={{
                            padding: "8px 16px",
                            background: "var(--bg-surface)",
                            border: "1px solid var(--twins-red)",
                            borderRadius: "6px",
                            color: "var(--twins-red)",
                            cursor: "pointer",
                            fontSize: "13px"
                        }}
                    >
                        🔊 Disparar Tortura (SCREAM_01)
                    </button>
                </div>
            </div>

            {/* Stats */}
            <div style={{
                display: "grid",
                gridTemplateColumns: "repeat(4, 1fr)",
                gap: "16px",
                marginBottom: "24px"
            }}>
                <StatCard label="Decisiones Hoy" value="12" />
                <StatCard label="Bloqueadas (MAD)" value="3" color="var(--twins-red)" />
                <StatCard label="Ejecuciones" value="7" color="var(--sanity-green)" />
                <StatCard label="En Sombra" value="2" color="var(--warning-amber)" />
            </div>

            {/* Decision History */}
            <h3 style={{
                fontSize: "1rem",
                fontWeight: 600,
                marginBottom: "16px",
                display: "flex",
                alignItems: "center",
                gap: "8px"
            }}>
                📋 Historial de Decisiones
            </h3>

            <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                {decisions.map((decision) => (
                    <DecisionCard
                        key={decision.id}
                        decision={decision}
                        emoji={actionEmoji[decision.action] || "❓"}
                    />
                ))}
            </div>
        </div>
    );
}

function StatCard({ label, value, color = "var(--text-primary)" }: {
    label: string;
    value: string;
    color?: string;
}) {
    return (
        <div className="card" style={{ textAlign: "center", padding: "16px" }}>
            <div style={{ fontSize: "28px", fontWeight: 700, color }}>{value}</div>
            <div style={{ fontSize: "12px", color: "var(--text-muted)", marginTop: "4px" }}>{label}</div>
        </div>
    );
}

function DecisionCard({ decision, emoji }: { decision: Decision; emoji: string }) {
    return (
        <div
            className="card"
            style={{
                padding: "16px",
                display: "flex",
                alignItems: "center",
                gap: "16px",
                borderLeft: decision.approved
                    ? "3px solid var(--sanity-green)"
                    : "3px solid var(--twins-red)",
                opacity: decision.shadow ? 0.7 : 1,
            }}
        >
            {/* Icon */}
            <div style={{
                width: "40px",
                height: "40px",
                borderRadius: "8px",
                background: "var(--bg-surface)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontSize: "20px",
            }}>
                {emoji}
            </div>

            {/* Details */}
            <div style={{ flex: 1 }}>
                <div style={{ fontWeight: 600, fontSize: "14px" }}>
                    {decision.action.replace("_", " ")}
                </div>
                <div style={{ fontSize: "12px", color: "var(--text-muted)" }}>
                    Objetivo: {decision.target}
                </div>
            </div>

            {/* Status Badges */}
            <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
                {decision.shadow && (
                    <span style={{
                        padding: "4px 8px",
                        background: "rgba(245, 158, 11, 0.15)",
                        border: "1px solid var(--warning-amber)",
                        borderRadius: "4px",
                        fontSize: "11px",
                        color: "var(--warning-amber)",
                    }}>
                        SHADOW
                    </span>
                )}
                {!decision.approved && (
                    <span style={{
                        padding: "4px 8px",
                        background: "rgba(220, 38, 38, 0.15)",
                        border: "1px solid var(--twins-red)",
                        borderRadius: "4px",
                        fontSize: "11px",
                        color: "var(--twins-red)",
                    }}>
                        MAD: {decision.madViolation || "BLOCKED"}
                    </span>
                )}
                {decision.approved && !decision.shadow && (
                    <span style={{
                        padding: "4px 8px",
                        background: "rgba(34, 197, 94, 0.15)",
                        border: "1px solid var(--sanity-green)",
                        borderRadius: "4px",
                        fontSize: "11px",
                        color: "var(--sanity-green)",
                    }}>
                        EJECUTADO
                    </span>
                )}
            </div>

            {/* Timestamp */}
            <div style={{
                fontSize: "12px",
                color: "var(--text-muted)",
                fontFamily: "var(--font-mono)",
            }}>
                {decision.timestamp}
            </div>
        </div>
    );
}
