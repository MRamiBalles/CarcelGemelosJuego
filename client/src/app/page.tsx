"use client";

import { useState, useEffect } from "react";
import PrisonerDashboard from "@/components/PrisonerDashboard";
import TwinsControlPanel from "@/components/TwinsControlPanel";
import { AdminPanel } from "@/components/admin/AdminPanel";
import VARReplay from "@/components/VARReplay";
import Header from "@/components/Header";
import { useGameEngine } from "@/hooks/useGameEngine";

// Mock data for demonstration (would come from WebSocket in production)
const mockPrisoners = [
    { id: "P001", name: "Simón", archetype: "VETERAN", sanity: 72, loyalty: 45, online: true, dignity: 100, pot_contribution: 0 },
    { id: "P002", name: "Elena", archetype: "MYSTIC", sanity: 85, loyalty: 60, online: true, dignity: 100, pot_contribution: 0 },
    { id: "P003", name: "Marco", archetype: "SHOWMAN", sanity: 55, loyalty: 30, online: false, dignity: 100, pot_contribution: 0 },
    { id: "P004", name: "Lucía", archetype: "REDEEMED", sanity: 90, loyalty: 75, online: true, dignity: 100, pot_contribution: 0 },
    { id: "P005", name: "Diego", archetype: "VETERAN", sanity: 40, loyalty: 20, online: true, dignity: 100, pot_contribution: 0 },
    { id: "P006", name: "Carla", archetype: "SHOWMAN", sanity: 65, loyalty: 55, online: false, dignity: 100, pot_contribution: 0 },
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

    // Live WebSocket connection to Go Engine
    const { events, isConnected, triggerOracle, triggerTorture } = useGameEngine();

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
