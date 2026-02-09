"use client";

import { useState, useEffect } from "react";
import PrisonerDashboard from "@/components/PrisonerDashboard";
import TwinsControlPanel from "@/components/TwinsControlPanel";
import VARReplay from "@/components/VARReplay";
import Header from "@/components/Header";

// Mock data for demonstration (would come from WebSocket in production)
const mockPrisoners = [
    { id: "P001", name: "Simón", archetype: "VETERAN", sanity: 72, loyalty: 45, online: true },
    { id: "P002", name: "Elena", archetype: "MYSTIC", sanity: 85, loyalty: 60, online: true },
    { id: "P003", name: "Marco", archetype: "SHOWMAN", sanity: 55, loyalty: 30, online: false },
    { id: "P004", name: "Lucía", archetype: "REDEEMED", sanity: 90, loyalty: 75, online: true },
    { id: "P005", name: "Diego", archetype: "VETERAN", sanity: 40, loyalty: 20, online: true },
    { id: "P006", name: "Carla", archetype: "SHOWMAN", sanity: 65, loyalty: 55, online: false },
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

            {/* Content */}
            <div className="container" style={{ padding: "24px" }}>
                {activeTab === "dashboard" && (
                    <PrisonerDashboard prisoners={mockPrisoners} />
                )}
                {activeTab === "twins" && (
                    <TwinsControlPanel
                        decisions={mockTwinsDecisions}
                        shadowMode={shadowMode}
                        onToggleShadowMode={() => setShadowMode(!shadowMode)}
                    />
                )}
                {activeTab === "var" && (
                    <VARReplay />
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
