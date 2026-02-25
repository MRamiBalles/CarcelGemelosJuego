"use client";

import { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface Prisoner {
    id: string;
    name: string;
    archetype: string;
    inventory: any[];
}

interface PlayerActionPanelProps {
    prisoners: Prisoner[];
    sendAction: (type: string, prisonerId: string, payload?: any) => void;
}

export default function PlayerActionPanel({ prisoners, sendAction }: PlayerActionPanelProps) {
    const [selectedActor, setSelectedActor] = useState<string>("P001");
    const [selectedTarget, setSelectedTarget] = useState<string>("P002");
    const [selectedItem, setSelectedItem] = useState<string>("");

    const handleAction = (type: string, payload: any = {}) => {
        sendAction(type, selectedActor, payload);
    };

    // Derived lists based on selection
    const actor = prisoners.find(p => p.id === selectedActor);
    const validTargets = prisoners.filter(p => p.id !== selectedActor);

    return (
        <Card className="bg-void border-border-subtle mt-6">
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    🕹️ Simulador de Jugador (F5)
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Actor Settings */}
                    <div className="space-y-4">
                        <div>
                            <label className="text-sm font-medium text-text-secondary block mb-1">Actor (Quién ejecuta la acción)</label>
                            <select
                                className="w-full bg-bg-surface border border-border-light rounded-md px-3 py-2 text-text-primary"
                                value={selectedActor}
                                onChange={(e) => setSelectedActor(e.target.value)}
                            >
                                {prisoners.map(p => (
                                    <option key={p.id} value={p.id}>{p.id} - {p.name}</option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className="text-sm font-medium text-text-secondary block mb-1">Objetivo Relacional (Robo / Snitch)</label>
                            <select
                                className="w-full bg-bg-surface border border-border-light rounded-md px-3 py-2 text-text-primary"
                                value={selectedTarget}
                                onChange={(e) => setSelectedTarget(e.target.value)}
                            >
                                {validTargets.map(p => (
                                    <option key={p.id} value={p.id}>{p.id} - {p.name}</option>
                                ))}
                            </select>
                        </div>

                        <div>
                            <label className="text-sm font-medium text-text-secondary block mb-1">Ítem a Consumir / Usar</label>
                            <select
                                className="w-full bg-bg-surface border border-border-light rounded-md px-3 py-2 text-text-primary"
                                value={selectedItem}
                                onChange={(e) => setSelectedItem(e.target.value)}
                                disabled={!actor || !actor.inventory || actor.inventory.length === 0}
                            >
                                <option value="">-- Selecciona Item --</option>
                                {actor?.inventory?.map(i => (
                                    <option key={i.type} value={i.type}>{i.type} ({i.quantity})</option>
                                ))}
                            </select>
                        </div>
                    </div>

                    {/* Action Buttons */}
                    <div className="grid grid-cols-2 gap-3">
                        <Button
                            variant="destructive"
                            onClick={() => handleAction("STEAL", { targetCell: "CELL_A" })} // Simplified for simulation
                            title="Robar a P"
                        >
                            🥷 STEAL (Cell)
                        </Button>

                        <Button
                            variant="outline"
                            className="text-twins-red border-twins-red hover:bg-twins-red/10"
                            onClick={() => handleAction("SNITCH", { target_id: selectedTarget, action: "SUSPICIOUS" })}
                        >
                            🐀 SNITCH Target
                        </Button>

                        <Button
                            variant="outline"
                            className="bg-sanity-green/10 text-sanity-green hover:bg-sanity-green/20"
                            onClick={() => handleAction("EAT", { item_type: selectedItem })}
                            disabled={!selectedItem}
                        >
                            🥪 COMER Ítem
                        </Button>

                        <Button
                            variant="secondary"
                            onClick={() => handleAction("TOILET")}
                        >
                            🚽 USAR BAÑO
                        </Button>

                        {actor?.archetype === "MYSTIC" && (
                            <>
                                <Button
                                    className="bg-loyalty-blue hover:bg-loyalty-blue/80 text-white"
                                    onClick={() => handleAction("MEDITATE", { target_cells: [] })}
                                >
                                    🧘 MEDITAR (Congelar)
                                </Button>
                                <Button
                                    className="bg-purple-600 hover:bg-purple-500 text-white"
                                    onClick={() => handleAction("USE_ORACLE", { target_id: selectedTarget })}
                                >
                                    🔮 ORÁCULO Target
                                </Button>
                            </>
                        )}

                        <Button
                            className="bg-amber-600 hover:bg-amber-500 text-white col-span-2"
                            onClick={() => handleAction("USE_RED_PHONE")}
                        >
                            ☎️ CONTESTAR TELÉFONO ROJO
                        </Button>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}
