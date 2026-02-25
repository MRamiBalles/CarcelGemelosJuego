"use client";

import { useState, useEffect } from "react";
import PrisonerDashboard from "@/components/PrisonerDashboard";
import TwinsControlPanel from "@/components/TwinsControlPanel";
import { AdminPanel } from "@/components/admin/AdminPanel";
import PollWidget from "@/components/admin/PollWidget";
import VARReplay from "@/components/VARReplay";
import Header from "@/components/Header";
import { useGameEngine, GameEvent } from "@/hooks/useGameEngine";

// Mock data for demonstration (would come from WebSocket in production)
const mockPrisoners = [
    { id: "P001", name: "Simón", archetype: "VETERAN", sanity: 72, loyalty: 45, online: true, dignity: 100, pot_contribution: 0, stamina: 80, hunger: 20, thirst: 10, states: { Asleep: true }, inventory: [{ type: "WATER", quantity: 2, is_contraband: false }] },
    { id: "P002", name: "Elena", archetype: "MYSTIC", sanity: 85, loyalty: 60, online: true, dignity: 100, pot_contribution: 0, stamina: 100, hunger: 0, thirst: 0, states: { Meditating: true }, inventory: [{ type: "ELIXIR", quantity: 1, is_contraband: true }] },
    { id: "P003", name: "Marco", archetype: "SHOWMAN", sanity: 55, loyalty: 30, online: false, dignity: 100, pot_contribution: 0, stamina: 10, hunger: 50, thirst: 40, states: { Exhausted: true }, inventory: [] },
    { id: "P004", name: "Lucía", archetype: "REDEEMED", sanity: 90, loyalty: 75, online: true, dignity: 100, pot_contribution: 0, stamina: 60, hunger: 10, thirst: 20, states: {}, inventory: [{ type: "RICE", quantity: 1, is_contraband: false }, { type: "SUSHI", quantity: 1, is_contraband: true }] },
    { id: "P005", name: "Diego", archetype: "VETERAN", sanity: 40, loyalty: 20, online: true, dignity: 100, pot_contribution: 0, stamina: 50, hunger: 60, thirst: 70, states: { Isolated: true }, inventory: [] },
    { id: "P006", name: "Carla", archetype: "SHOWMAN", sanity: 65, loyalty: 55, online: false, dignity: 100, pot_contribution: 0, stamina: 0, hunger: 100, thirst: 100, states: { Dead: true }, inventory: [] },
];

const mockTwinsDecisions = [
    { id: "D001", timestamp: "17:30", action: "NOISE_TORTURE", target: "BLOCK_A", approved: true, shadow: false },
    { id: "D002", timestamp: "16:45", action: "REVEAL_SECRET", target: "P003", approved: true, shadow: true },
    { id: "D003", timestamp: "15:20", action: "RESOURCE_CUT", target: "ALL", approved: false, shadow: true, madViolation: "NO_DAY_ONE_CRUELTY" },
    { id: "D004", timestamp: "14:00", action: "OBSERVE", target: "NONE", approved: true, shadow: false },
];

export default function Home() {
    const [activeTab, setActiveTab] = useState<"dashboard" | "twins" | "var">("dashboard");
    const [shadowMode, setShadowMode] = useState(true);
    const [gameDay, setGameDay] = useState(7);
    const [tensionLevel, setTensionLevel] = useState("HIGH");
    const [toasts, setToasts] = useState<{ id: string, message: string, type: string }[]>([]);

    // Live WebSocket connection to Go Engine
    const { events, isConnected, triggerOracle, triggerTorture, createPoll, votePoll } = useGameEngine();

    // Effect to catch new events and create toasts for significant ones
    useEffect(() => {
        if (events.length > 0) {
            const latestEvent = events[0];
            let toastMessage = "";
            let toastType = "info";

            switch (latestEvent.type) {
                case "STEAL":
                    toastMessage = `🚨 ROBBERY: ${latestEvent.actor_id} intentó robar.`;
                    toastType = "warning";
                    break;
                case "SNITCH":
                    toastMessage = `🐀 SNITCH: ${latestEvent.actor_id} chivó sobre ${latestEvent.target_id}!`;
                    toastType = "error";
                    break;
                case "AUDIO_TORTURE":
                    toastMessage = `🔊 TORTURA AUDITIVA: ${latestEvent.payload.soundName || "Ruido"}`;
                    toastType = "error";
                    break;
                case "RED_PHONE_MESSAGE":
                    toastMessage = `☎️ EL TELÉFONO ROJO ESTÁ SONANDO...`;
                    toastType = "warning";
                    break;
                case "TOILET_USE":
                    toastMessage = `🚽 ${latestEvent.actor_id} está usando el baño.`;
                    break;
                case "AUDIENCE_EXPULSION":
                    toastMessage = `💀 LA AUDIENCIA HA EXPULSADO A ${latestEvent.target_id}`;
                    toastType = "error";
                    break;
                case "MEDICAL_EVACUATION":
                    toastMessage = `🚑 EMERGENCIA: ${latestEvent.target_id} requiere evacuación.`;
                    toastType = "error";
                    break;
                case "ORACLE_USE":
                    toastMessage = `🔮 ORÁCULO: Tartaria ha revelado un secreto mortal de ${latestEvent.target_id}`;
                    toastType = "warning";
                    break;
                case "MEDITATE":
                    toastMessage = `🧘 Tartaria está meditando (Drenando celdas contiguas)`;
                    toastType = "info";
                    break;
            }

            if (toastMessage) {
                const id = latestEvent.id;
                setToasts(prev => {
                    // Prevent duplicate toasts if React re-renders quickly
                    if (prev.find(t => t.id === id)) return prev;
                    return [...prev, { id, message: toastMessage, type: toastType }].slice(-5);
                });

                // Auto-remove toast after 5 seconds
                setTimeout(() => {
                    setToasts(prev => prev.filter(t => t.id !== id));
                }, 5000);
            }
        }
    }, [events]);

    return (
        <main className="min-h-screen" style={{ background: "var(--bg-void)" }}>
            <Header
                gameDay={gameDay}
                tensionLevel={tensionLevel}
                shadowMode={shadowMode}
                onToggleShadowMode={() => setShadowMode(!shadowMode)}
            />

            {/* Navigation Tabs */}
            <nav style={{
                background: "var(--bg-cell)",
                borderBottom: "1px solid var(--border-subtle)",
                padding: "0 24px"
            }}>
                <div className="container" style={{ display: "flex", gap: "4px" }}>
                    <TabButton
                        active={activeTab === "dashboard"}
                        onClick={() => setActiveTab("dashboard")}
                    >
                        👁️ Prisioneros
                    </TabButton>
                    <TabButton
                        active={activeTab === "twins"}
                        onClick={() => setActiveTab("twins")}
                    >
                        🤖 Los Gemelos
                    </TabButton>
                    <TabButton
                        active={activeTab === "var"}
                        onClick={() => setActiveTab("var")}
                    >
                        📼 VAR Replay
                    </TabButton>
                </div>
            </nav>

            {/* Toasts Container */}
            <div style={{
                position: "fixed",
                bottom: "24px",
                right: "24px",
                display: "flex",
                flexDirection: "column",
                gap: "8px",
                zIndex: 1000
            }}>
                {toasts.map(toast => (
                    <div key={toast.id} style={{
                        background: toast.type === "error" ? "var(--twins-red)" : toast.type === "warning" ? "var(--warning-amber)" : "var(--bg-surface)",
                        color: toast.type === "error" ? "white" : toast.type === "warning" ? "black" : "var(--text-primary)",
                        padding: "12px 24px",
                        borderRadius: "8px",
                        boxShadow: "0 10px 25px -5px rgba(0, 0, 0, 0.5)",
                        border: "1px solid rgba(255,255,255,0.1)",
                        animation: "slideIn 0.3s ease-out forwards",
                        fontWeight: 500,
                        fontSize: "14px",
                        maxWidth: "400px",
                    }}>
                        {toast.message}
                    </div>
                ))}
            </div>

            {/* Content  - Powered by WebSocket Events */}
            <div className="container" style={{ padding: "24px" }}>
                {!isConnected && (
                    <div style={{ color: "var(--warning-amber)", marginBottom: "16px", padding: "12px", background: "rgba(255,193,7,0.1)", border: "1px solid var(--warning-amber)", borderRadius: "6px" }}>
                        ⚠️ Conexión con el Servidor (VAR) interrumpida. Mostrando datos mockeados o no disponibles.
                    </div>
                )}

                {activeTab === "dashboard" && (
                    <PrisonerDashboard prisoners={mockPrisoners} />
                )}
                {activeTab === "twins" && (
                    <div className="flex flex-col gap-6">
                        <AdminPanel />
                        <TwinsControlPanel
                            decisions={mockTwinsDecisions}
                            shadowMode={shadowMode}
                            onToggleShadowMode={() => setShadowMode(!shadowMode)}
                            // Adding trigger callbacks mapped to the API
                            onTriggerOracle={(target, message) => triggerOracle(target, message)}
                            onTriggerTorture={(soundId) => triggerTorture(soundId)}
                        />
                        <PollWidget events={events} createPoll={createPoll} votePoll={votePoll} />
                    </div>
                )}
                {activeTab === "var" && (
                    <VARReplay liveEvents={events} />
                )}
            </div>
        </main>
    );
}

function TabButton({ active, onClick, children }: {
    active: boolean;
    onClick: () => void;
    children: React.ReactNode;
}) {
    return (
        <button
            onClick={onClick}
            style={{
                padding: "16px 24px",
                background: "transparent",
                border: "none",
                borderBottom: active ? "2px solid var(--twins-red)" : "2px solid transparent",
                color: active ? "var(--text-primary)" : "var(--text-secondary)",
                fontSize: "14px",
                fontWeight: 500,
                cursor: "pointer",
                transition: "all 0.2s ease",
            }}
        >
            {children}
        </button>
    );
}
