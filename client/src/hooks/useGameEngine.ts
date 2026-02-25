import { useState, useEffect, useRef } from 'react';

// Types mirroring the Go backend
export type EventType =
    | "RESOURCE_INTAKE" | "NOISE_EVENT" | "SOCIAL_ACTION" | "VOTE"
    | "BETRAYAL" | "PRIVACY_BREACH" | "TIME_TICK" | "SANITY_CHANGE"
    | "LOYALTY_CHANGE" | "TOILET_USE" | "DOOR_LOCK" | "DOOR_OPEN"
    | "AUDIO_TORTURE" | "AGGRESSIVE_EMOTE" | "LOCKDOWN_BANG"
    | "INSULT" | "STEAL" | "FINAL_DILEMMA_START" | "FINAL_DILEMMA_DECISION"
    | "POLL_CREATED" | "POLL_RESOLVED" | "ISOLATION_CHANGED" | "AUDIENCE_EXPULSION";

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
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const reconnectAttemptsRef = useRef(0);

    const connect = () => {
        // Prevent multiple connections
        if (wsRef.current && (wsRef.current.readyState === WebSocket.OPEN || wsRef.current.readyState === WebSocket.CONNECTING)) {
            return;
        }

        console.log(`Attempting to connect to ${url}...`);
        const ws = new WebSocket(url);

        ws.onopen = () => {
            console.log("Connected to Game Engine WebSocket");
            setIsConnected(true);
            reconnectAttemptsRef.current = 0; // Reset attempts on success
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

            // Exponential backoff reconnection
            const timeout = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current), 30000); // Max 30s
            console.log(`Reconnecting in ${timeout / 1000} seconds...`);

            reconnectTimeoutRef.current = setTimeout(() => {
                reconnectAttemptsRef.current++;
                connect();
            }, timeout);
        };

        ws.onerror = (error) => {
            console.error("WebSocket Error:", error);
            // OnError will usually be followed by OnClose, handled there.
        };

        wsRef.current = ws;
    };

    useEffect(() => {
        connect();

        return () => {
            if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);
            if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
                // Ensure we don't trigger the reconnect loop on deliberate unmount
                wsRef.current.onclose = null;
                wsRef.current.close();
            }
        };
    }, [url]);

    // Send a player action to the WebSocket directly (F5 implementation)
    const sendAction = (type: string, prisonerId: string, payload: any = {}) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            const message = JSON.stringify({
                type: type,
                prisoner_id: prisonerId,
                payload: payload
            });
            wsRef.current.send(message);
        } else {
            console.warn("Cannot send action, WebSocket is not open");
        }
    };

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

    const voteExpel = async (prisonerId: string) => {
        try {
            await fetch('http://localhost:8080/api/audience/vote_expel', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ prisoner_id: prisonerId })
            });
        } catch (err) {
            console.error("Failed to trigger Expulsion API", err);
        }
    };

    return { events, isConnected, sendAction, triggerOracle, triggerTorture, createPoll, votePoll, voteExpel };
}
