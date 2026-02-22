"use client";

import { useState } from "react";
import { GameEvent } from "@/hooks/useGameEngine";

interface VARReplayProps {
    liveEvents?: GameEvent[];
}

export default function VARReplay({ liveEvents = [] }: VARReplayProps) {
    const [filter, setFilter] = useState<string>("ALL");
    const [showRevealed, setShowRevealed] = useState<boolean>(true);

    const filteredEvents = liveEvents.filter(e => {
        if (filter !== "ALL" && e.type !== filter) return false;
        if (!showRevealed && e.is_revealed) return false;
        return true;
    });

    return (
        <div>
            {/* Header */}
            <div style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                marginBottom: "24px"
            }}>
                <div>
                    <h2 style={{ fontSize: "1.25rem", fontWeight: 600 }}>
                        📼 VAR de la Traición
                    </h2>
                    <p style={{ color: "var(--text-muted)", fontSize: "13px", marginTop: "4px" }}>
                        Historial inmutable de todos los eventos • Append-Only Ledger
                    </p>
                </div>

                {/* Stats */}
                <div style={{ display: "flex", gap: "16px" }}>
                    <MiniStat label="Total Eventos" value={liveEvents.length.toString()} />
                    <MiniStat label="Traiciones" value={liveEvents.filter(e => e.type === "BETRAYAL").length.toString()} color="var(--twins-red)" />
                    <MiniStat label="Anomalías Ocultas" value={liveEvents.filter(e => !e.is_revealed).length.toString()} color="var(--warning-amber)" />
                </div>
            </div>

            {/* Filters */}
            <div style={{
                display: "flex",
                gap: "8px",
                marginBottom: "24px",
                flexWrap: "wrap"
            }}>
                <FilterButton active={filter === "ALL"} onClick={() => setFilter("ALL")}>
                    Todos
                </FilterButton>
                <FilterButton active={filter === "BETRAYAL"} onClick={() => setFilter("BETRAYAL")}>
                    🗡️ Traiciones
                </FilterButton>
                <FilterButton active={filter === "NOISE_EVENT"} onClick={() => setFilter("NOISE_EVENT")}>
                    🔊 Ruido
                </FilterButton>
                <FilterButton active={filter === "SANITY_CHANGE"} onClick={() => setFilter("SANITY_CHANGE")}>
                    🧠 Cordura
                </FilterButton>
                <FilterButton active={filter === "TWINS_DECISION"} onClick={() => setFilter("TWINS_DECISION")}>
                    🤖 Decisiones
                </FilterButton>

                <div style={{ marginLeft: "auto" }}>
                    <label style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        color: "var(--text-secondary)",
                        fontSize: "13px",
                        cursor: "pointer"
                    }}>
                        <input
                            type="checkbox"
                            checked={!showRevealed}
                            onChange={() => setShowRevealed(!showRevealed)}
                            style={{ accentColor: "var(--twins-red)" }}
                        />
                        Solo secretos ocultos
                    </label>
                </div>
            </div>

            {/* Timeline */}
            <div style={{ display: "flex", flexDirection: "column", gap: "4px" }}>
                {filteredEvents.map((event) => (
                    <EventCard key={event.id} event={event} />
                ))}

                {filteredEvents.length === 0 && (
                    <div style={{
                        textAlign: "center",
                        padding: "48px",
                        color: "var(--text-muted)"
                    }}>
                        No hay eventos que coincidan con los filtros
                    </div>
                )}
            </div>
        </div>
    );
}

function MiniStat({ label, value, color = "var(--text-primary)" }: {
    label: string;
    value: string;
    color?: string;
}) {
    return (
        <div style={{ textAlign: "center" }}>
            <div style={{ fontSize: "20px", fontWeight: 700, color }}>{value}</div>
            <div style={{ fontSize: "10px", color: "var(--text-muted)" }}>{label}</div>
        </div>
    );
}

function FilterButton({ active, onClick, children }: {
    active: boolean;
    onClick: () => void;
    children: React.ReactNode;
}) {
    return (
        <button
            onClick={onClick}
            style={{
                padding: "8px 16px",
                background: active ? "var(--twins-red)" : "var(--bg-surface)",
                border: "1px solid",
                borderColor: active ? "var(--twins-red)" : "var(--border-subtle)",
                borderRadius: "20px",
                color: active ? "white" : "var(--text-secondary)",
                fontSize: "13px",
                cursor: "pointer",
                transition: "all 0.2s ease",
            }}
        >
            {children}
        </button>
    );
}

function EventCard({ event }: { event: typeof mockEvents[0] }) {
    const impactColor = {
        POSITIVE: "var(--sanity-green)",
        NEGATIVE: "var(--twins-red)",
        NEUTRAL: "var(--text-muted)",
    }[event.impact];

    const typeEmoji: Record<string, string> = {
        NOISE_EVENT: "🔊",
        SANITY_CHANGE: "🧠",
        BETRAYAL: "🗡️",
        TWINS_DECISION: "🤖",
        AUDIENCE_TORTURE: "📺",
        RESOURCE_INTAKE: "🍞",
    };

    return (
        <div
            className={`timeline-event ${event.impact.toLowerCase()}`}
            style={{
                padding: "16px",
                paddingLeft: "32px",
                background: "var(--bg-cell)",
                borderRadius: "8px",
                marginLeft: "8px",
            }}
        >
            <div style={{ display: "flex", alignItems: "flex-start", gap: "16px" }}>
                {/* Icon */}
                <div style={{
                    width: "36px",
                    height: "36px",
                    borderRadius: "8px",
                    background: "var(--bg-surface)",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "18px",
                    flexShrink: 0,
                }}>
                    {typeEmoji[event.type] || "❓"}
                </div>

                {/* Content */}
                <div style={{ flex: 1 }}>
                    <div style={{
                        fontSize: "14px",
                        fontWeight: 500,
                        marginBottom: "4px",
                        color: event.revealed ? "inherit" : "var(--warning-amber)",
                        fontStyle: event.revealed ? "normal" : "italic"
                    }}>
                        {event.revealed ? event.summary : "⚠️ Anomalía Detectada: Procesando información del VAR..."}
                    </div>
                    <div style={{
                        fontSize: "12px",
                        color: "var(--text-muted)",
                        display: "flex",
                        gap: "16px"
                    }}>
                        <span>Actor: {event.revealed ? event.actor : "[DATOS_CORRUPTOS]"}</span>
                        <span>Día {event.day}</span>
                    </div>
                </div>

                {/* Badges */}
                <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
                    {!event.revealed && (
                        <span style={{
                            padding: "4px 8px",
                            background: "rgba(245, 158, 11, 0.15)",
                            border: "1px solid var(--warning-amber)",
                            borderRadius: "4px",
                            fontSize: "10px",
                            color: "var(--warning-amber)",
                        }}>
                            OCULTO
                        </span>
                    )}
                    <span style={{
                        padding: "4px 8px",
                        background: `${impactColor}20`,
                        border: `1px solid ${impactColor}`,
                        borderRadius: "4px",
                        fontSize: "10px",
                        color: impactColor,
                    }}>
                        {event.impact}
                    </span>
                </div>

                {/* Timestamp */}
                <div style={{
                    fontSize: "12px",
                    color: "var(--text-muted)",
                    fontFamily: "var(--font-mono)",
                    flexShrink: 0,
                }}>
                    {event.timestamp}
                </div>
            </div>
        </div>
    );
}
