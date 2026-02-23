import React, { useState } from 'react';
import { GameEvent } from '@/hooks/useGameEngine';

interface PollWidgetProps {
    events: GameEvent[];
    createPoll: (payload: any) => Promise<void>;
    votePoll: (pollId: string, option: string) => Promise<void>;
}

export default function PollWidget({ events, createPoll, votePoll }: PollWidgetProps) {
    const [question, setQuestion] = useState('');
    const [options, setOptions] = useState('P_001,P_002');
    const [rewardType, setRewardType] = useState('SUSHI');
    const [durationSec, setDurationSec] = useState(30);

    // Compute active polls
    const createdPolls = events.filter(e => e.type === 'POLL_CREATED');
    const resolvedPolls = events.filter(e => e.type === 'POLL_RESOLVED');
    const resolvedIds = new Set(resolvedPolls.map(e => e.payload.poll_id));

    // An active poll is one that hasn't been resolved yet
    const activePolls = createdPolls.filter(e => !resolvedIds.has(e.payload.poll_id));

    const handleCreate = async () => {
        const payload = {
            poll_id: `POLL_${Date.now()}`,
            question,
            options: options.split(',').map(o => o.trim()),
            reward_type: rewardType,
            duration_sec: durationSec
        };
        await createPoll(payload);
        setQuestion('');
    };

    return (
        <div className="bg-slate-800 p-4 rounded shadow-md mt-4">
            <h2 className="text-xl font-bold mb-4 text-white">Twitch Audience Polls</h2>

            <div className="flex flex-col mb-6 gap-2 bg-slate-900 p-4 rounded">
                <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-widest">Create New Poll</h3>
                <input
                    type="text"
                    placeholder="Question (e.g. Who gets sushi?)"
                    value={question}
                    onChange={e => setQuestion(e.target.value)}
                    className="p-2 rounded bg-slate-700 text-white"
                />
                <input
                    type="text"
                    placeholder="Options (comma separated IDs)"
                    value={options}
                    onChange={e => setOptions(e.target.value)}
                    className="p-2 rounded bg-slate-700 text-white"
                />
                <div className="flex gap-2 w-full">
                    <select
                        value={rewardType}
                        onChange={e => setRewardType(e.target.value)}
                        className="p-2 flex-1 rounded bg-slate-700 text-white"
                    >
                        <option value="SUSHI">Reward: Sushi</option>
                        <option value="TORTURE">Punishment: Torture</option>
                        <option value="ISOLATION">Punishment: Isolation</option>
                    </select>
                    <input
                        type="number"
                        value={durationSec}
                        onChange={e => setDurationSec(Number(e.target.value))}
                        className="p-2 w-24 rounded bg-slate-700 text-white"
                        title="Duration (Seconds)"
                    />
                </div>
                <button
                    onClick={handleCreate}
                    className="mt-2 bg-blue-600 hover:bg-blue-500 text-white p-2 rounded transition-colors font-bold"
                >
                    Start Broadcast Poll
                </button>
            </div>

            <div className="flex flex-col gap-4">
                <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-widest border-b border-slate-700 pb-2">Active Polls</h3>
                {activePolls.length === 0 ? (
                    <p className="text-sm text-gray-500 italic">No active polls running.</p>
                ) : (
                    activePolls.map(pollEvent => {
                        const payload = pollEvent.payload;
                        return (
                            <div key={pollEvent.id} className="bg-slate-700 p-3 rounded flex flex-col gap-2">
                                <div className="flex justify-between items-start">
                                    <h4 className="font-bold text-lg text-white">{payload.question}</h4>
                                    <span className="text-xs font-mono bg-blue-900 text-blue-300 px-2 py-1 rounded">
                                        ⏱️ {payload.duration_sec}s
                                    </span>
                                </div>
                                <p className="text-xs text-gray-400">Modifier: {payload.reward_type}</p>
                                <div className="flex gap-2 mt-2">
                                    {payload.options.map((opt: string) => (
                                        <button
                                            key={opt}
                                            onClick={() => votePoll(payload.poll_id, opt)}
                                            className="flex-1 bg-slate-600 hover:bg-slate-500 text-white p-2 text-sm rounded transition-colors"
                                        >
                                            Vote: {opt}
                                        </button>
                                    ))}
                                </div>
                            </div>
                        );
                    })
                )}
            </div>

            {/* Show last resolved poll briefly (Optional) */}
            <div className="mt-6 border-t border-slate-700 pt-4">
                <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-widest mb-2">Recent Results</h3>
                {resolvedPolls.slice(0, 3).map(p => (
                    <div key={p.id} className="text-sm text-white bg-slate-900 p-2 mb-2 rounded border border-green-800">
                        <span className="font-bold text-green-400">Winner:</span> {p.payload.winner_option} ({p.payload.reward_type})
                    </div>
                ))}
            </div>
        </div>
    );
}
