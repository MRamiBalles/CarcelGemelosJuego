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
    onToggleShadowMode: () => void;
}

export default function TwinsControlPanel({ decisions, shadowMode, onToggleShadowMode }: TwinsControlPanelProps) {
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
