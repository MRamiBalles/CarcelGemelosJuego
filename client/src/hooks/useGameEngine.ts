import { useState, useEffect, useRef } from 'react';

// Types mirroring the Go backend
export type EventType =
    | "RESOURCE_INTAKE" | "NOISE_EVENT" | "SOCIAL_ACTION" | "VOTE"
    | "BETRAYAL" | "PRIVACY_BREACH" | "TIME_TICK" | "SANITY_CHANGE"
    | "LOYALTY_CHANGE" | "TOILET_USE" | "DOOR_LOCK" | "DOOR_OPEN"
    | "AUDIO_TORTURE" | "AGGRESSIVE_EMOTE" | "LOCKDOWN_BANG"
    | "INSULT" | "STEAL" | "FINAL_DILEMMA_START" | "FINAL_DILEMMA_DECISION"
    | "ORACLE_PAINFUL_TRUTH" | "POLL_CREATED" | "POLL_RESOLVED" | "ISOLATION_CHANGED";

export interface GameEvent {
    id: string;
    timestamp: string;
    type: EventType;
    actor_id: string;
    target_id: string;
    payload: any;
    game_day: number;
    is_revealed: boolean;
}

export function useGameEngine(url: string = 'ws://localhost:8080/ws') {
    const [events, setEvents] = useState<GameEvent[]>([]);
    const [isConnected, setIsConnected] = useState(false);
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        // Initialize WebSocket connection
        const ws = new WebSocket(url);

        ws.onopen = () => {
            console.log("Connected to Game Engine WebSocket");
            setIsConnected(true);
        };

        ws.onmessage = (event) => {
            try {
                const gameEvent: GameEvent = JSON.parse(event.data);

                // Add new event to the beginning of the list (Max 100 for memory)
                setEvents((prevEvents) => {
                    const newEvents = [gameEvent, ...prevEvents];
                    if (newEvents.length > 100) {
                        return newEvents.slice(0, 100);
                    }
                    return newEvents;
                });

            } catch (err) {
                console.error("Failed to parse incoming WS message:", err);
            }
        };

        ws.onclose = () => {
            console.log("Disconnected from Game Engine WebSocket");
            setIsConnected(false);
            // Could add reconnection logic here
        };

        ws.onerror = (error) => {
            console.error("WebSocket Error:", error);
            ws.close();
        };

        wsRef.current = ws;

        return () => {
            if (ws.readyState === WebSocket.OPEN) {
                ws.close();
            }
        };
    }, [url]);

    // Expose a way to manually trigger audience APIs
    const triggerOracle = async (target: string, message: string) => {
        try {
            await fetch('http://localhost:8080/api/trigger-oracle', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ target, message })
            });
        } catch (err) {
            console.error("Failed to trigger Oracle API", err);
        }
    };

    const triggerTorture = async (soundId: string) => {
        try {
            await fetch('http://localhost:8080/api/trigger-torture', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ soundId })
            });
        } catch (err) {
            console.error("Failed to trigger Torture API", err);
        }
    };

    const createPoll = async (payload: any) => {
        try {
            await fetch('http://localhost:8080/api/poll/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
        } catch (err) {
            console.error("Failed to trigger Poll API", err);
        }
    };

    const votePoll = async (pollId: string, option: string) => {
        try {
            await fetch('http://localhost:8080/api/poll/vote', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ poll_id: pollId, option })
            });
        } catch (err) {
            console.error("Failed to trigger Vote API", err);
        }
    };

    return { events, isConnected, triggerOracle, triggerTorture, createPoll, votePoll };
}
