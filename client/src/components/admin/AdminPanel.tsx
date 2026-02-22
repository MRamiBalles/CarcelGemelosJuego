import React, { useState } from 'react';
import { Card, CardHeader, CardContent, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Brain, AlertOctagon, Terminal, Activity, Zap } from 'lucide-react';

interface Decision {
    action_type: string;
    target: string;
    intensity: number;
    reason?: string;
    justification: string;
    is_approved: boolean;
    metadata?: {
        reasoning?: string;
        model?: string;
    };
}

export function AdminPanel() {
    const [loading, setLoading] = useState(false);
    const [decision, setDecision] = useState<Decision | null>(null);
    const [error, setError] = useState<string | null>(null);

    const forceDecision = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await fetch('http://localhost:8080/api/twins/force-decision', {
                method: 'POST',
            });

            if (!response.ok) {
                throw new Error(`API Error: ${response.status} ${response.statusText}`);
            }

            const data = await response.json();
            if (data.status === 'ok' && data.decision) {
                setDecision(data.decision);
            } else {
                throw new Error('Invalid response format');
            }
        } catch (err: any) {
            setError(err.message || 'Failed to trigger decision');
        } finally {
            setLoading(false);
        }
    };

    const getIntensityColor = (level: number) => {
        switch (level) {
            case 1: return 'text-yellow-400 bg-yellow-400/10';
            case 2: return 'text-orange-500 bg-orange-500/10';
            case 3: return 'text-red-500 bg-red-500/10 hover:bg-red-500/20 shadow-[0_0_15px_rgba(239,68,68,0.5)]';
            default: return 'text-gray-400 bg-gray-400/10';
        }
    };

    const isBlocked = decision && !decision.is_approved;

    return (
        <Card className="w-full bg-slate-900 border-slate-700 shadow-xl overflow-hidden">
            <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-red-600 to-amber-500"></div>

            <CardHeader className="bg-slate-900/80 border-b border-slate-800 pb-4">
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                        <div className="p-2 bg-red-500/20 rounded-lg">
                            <Brain className="w-6 h-6 text-red-500" />
                        </div>
                        <div>
                            <CardTitle className="text-xl font-bold text-slate-100 uppercase tracking-wider">
                                PANÓPTICO ADMIN (GEMELOS)
                            </CardTitle>
                            <CardDescription className="text-slate-400">
                                Cognition Engine Override Console
                            </CardDescription>
                        </div>
                    </div>

                    <Button
                        onClick={forceDecision}
                        disabled={loading}
                        className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-6 rounded-md border border-red-500 shadow-[0_0_10px_rgba(220,38,38,0.4)] hover:shadow-[0_0_20px_rgba(220,38,38,0.6)] transition-all flex items-center gap-2 uppercase tracking-widest text-sm"
                    >
                        {loading ? (
                            <><Activity className="w-4 h-4 animate-spin" /> Procesando LLM...</>
                        ) : (
                            <><Zap className="w-4 h-4" /> Forzar Decisión Nv.1</>
                        )}
                    </Button>
                </div>
            </CardHeader>

            <CardContent className="p-0">
                {error && (
                    <div className="m-6 p-4 bg-red-900/30 border border-red-500/50 rounded-lg flex items-start gap-3 text-red-200">
                        <AlertOctagon className="w-5 h-5 text-red-500 mt-0.5" />
                        <div>
                            <h4 className="font-bold text-red-400">Error en el Oráculo</h4>
                            <p className="text-sm opacity-80">{error}</p>
                        </div>
                    </div>
                )}

                {!decision && !error && !loading && (
                    <div className="h-64 flex flex-col items-center justify-center text-slate-500 gap-4">
                        <Terminal className="w-12 h-12 opacity-20" />
                        <p className="tracking-widest uppercase text-sm font-medium">ESPERANDO INSTRUCCIÓN LLM</p>
                    </div>
                )}

                {decision && (
                    <div className="divide-y divide-slate-800">
                        {/* Action Header */}
                        <div className={`p-6 ${isBlocked ? 'bg-slate-800/50' : 'bg-red-950/20'}`}>
                            <div className="flex justify-between items-start mb-4">
                                <div>
                                    <div className="flex items-center gap-3 mb-1">
                                        <h3 className="text-2xl font-black text-slate-100 tracking-tight uppercase">
                                            {isBlocked ? 'BLOQUEADO' : decision.action_type}
                                        </h3>
                                        <span className={`px-2 py-0.5 rounded text-xs font-bold border ${isBlocked ? 'border-amber-500/50 text-amber-500 bg-amber-500/10' : 'border-green-500/50 text-green-500 bg-green-500/10'}`}>
                                            {isBlocked ? 'MAD RULES VIOLATION' : 'APPROVED'}
                                        </span>
                                    </div>
                                    <p className="text-slate-400 text-sm flex items-center gap-2">
                                        <span className="font-mono text-xs uppercase bg-slate-800 px-2 py-1 rounded">Target: {decision.target}</span>
                                    </p>
                                </div>

                                <div className={`flex flex-col items-end`}>
                                    <span className="text-xs text-slate-500 font-bold uppercase tracking-wider mb-1">Intensidad</span>
                                    <div className={`px-4 py-2 rounded-lg border font-black text-lg ${getIntensityColor(decision.intensity)}`}>
                                        LVL {decision.intensity}
                                    </div>
                                </div>
                            </div>

                            {/* Justification given to users */}
                            <div className="bg-slate-900 border border-slate-700 rounded-lg p-4 mt-4 relative overflow-hidden group">
                                <div className="absolute top-0 left-0 w-1 h-full bg-purple-500"></div>
                                <h4 className="text-purple-400 text-xs font-bold uppercase tracking-widest mb-2 font-mono flex items-center gap-2">
                                    <Terminal className="w-3 h-3" /> Dictamen Oficial (Audiencia)
                                </h4>
                                <p className="text-slate-200 leading-relaxed italic">
                                    "{decision.justification}"
                                </p>
                            </div>
                        </div>

                        {/* AI Chain of Thought (Metadata) */}
                        {decision.metadata?.reasoning && (
                            <div className="p-6 bg-slate-950/50 font-mono text-sm relative">
                                <div className="flex items-center justify-between mb-3 text-emerald-500 gap-2 uppercase tracking-widest text-xs font-bold">
                                    <span className="flex items-center gap-2">
                                        <Brain className="w-4 h-4" /> Rastro Cognitivo (Chain of Thought)
                                    </span>
                                    <span className="text-slate-600 bg-slate-900 px-2 py-1 rounded border border-slate-800">
                                        {decision.metadata?.model || 'Desconocido'}
                                    </span>
                                </div>
                                <div className="text-slate-400 whitespace-pre-wrap leading-relaxed border-l-2 border-emerald-900 pl-4 py-2 h-48 overflow-y-auto custom-scrollbar">
                                    {decision.metadata.reasoning}
                                </div>
                            </div>
                        )}
                    </div>
                )}
            </CardContent>
        </Card>
    );
}
