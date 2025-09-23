import { useEffect, useMemo } from "react";
import axios from "axios";
import useSpotifyPlayer from "../hooks/useSpotifyPlayer";
import TimedProgress from "./TimedProgress";
import type { Question } from "../pages/GamePage";

type Props = {
    question: Question | null;
    scoreboard: Record<string, number> | null;
    view: string;
    accessToken: string | null;
    playerID: string;
};

const HostGame: React.FC<Props> = ({
    question,
    scoreboard,
    view,
    accessToken,
}) => {
    const { playerReady } = useSpotifyPlayer(accessToken);

    const displayQuestionNumber = useMemo(() => {
        if (!question?.id || question.id.length < 2) return "";
        return question.id.slice(1);
    }, [question?.id]);

    useEffect(() => {
        const device_id = localStorage.getItem("device_id");
        if (!device_id || !accessToken || !playerReady || !window.player) return;

        if (view === "question" && question?.trackId) {
            axios
                .put(
                    `https://api.spotify.com/v1/me/player/play?device_id=${device_id}`,
                    {
                        uris: [`spotify:track:${question.trackId}`],
                        position_ms: question.positionMs,
                    },
                    {
                        headers: {
                            Authorization: `Bearer ${accessToken}`,
                        },
                    },
                )
                .then(() => console.log("Track playing"))
                .catch((err) => console.error("Play error:", err));
        }

        if (view === "scoreboard") {
            axios
                .put(
                    `https://api.spotify.com/v1/me/player/pause?device_id=${device_id}`,
                    {},
                    {
                        headers: {
                            Authorization: `Bearer ${accessToken}`,
                        },
                    },
                )
                .then(() => console.log("Player paused"))
                .catch((err) => console.error("Pause error:", err));
        }
    }, [view, question, playerReady, accessToken]);

    return (
        <div className="w-full max-w-2xl bg-gray-100 text-gray-800 p-6 rounded-xl shadow-lg flex flex-col items-center border border-gray-200">
            <div className="text-sm text-indigo-500 font-medium mb-3 uppercase tracking-wide">
                Host Panel
            </div>

            {view === "question" && question && (
                <>
                    <div className="mb-2 text-gray-700">
                        <span className="text-sm">Question </span>
                        <span className="font-semibold text-indigo-700">
                            {displayQuestionNumber || "—"}/10
                        </span>
                    </div>

                    <div className="w-full mb-4">
                        <TimedProgress duration={15} />
                    </div>

                    <div className="text-center space-y-2 mb-4">
                        <p className="text-gray-600">Pick an answer on Your phone!</p>
                    </div>

                    <div className="w-full">
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                            {question.options.map((option) => (
                                <div
                                    key={option}
                                    className="w-full px-4 py-3 rounded-lg text-left text-sm font-medium border bg-white text-gray-800 border-gray-300 shadow-sm"
                                >
                                    {option}
                                </div>
                            ))}
                        </div>
                    </div>
                </>
            )}

            {view === "scoreboard" && scoreboard && (
                <>
                    <div className="w-full mb-4">
                        <TimedProgress duration={5} />
                    </div>
                    <div className="bg-indigo-200 w-full">
                        <div className="px-4 py-2 rounded-t-md rounded-b-lg text-indigo-700 font-medium text-center">
                            Player Rankings
                        </div>
                        <ul className="divide-y divide-gray-200 rounded-b-sm">
                            {Object.entries(scoreboard)
                                .sort(([, a], [, b]) => b - a)
                                .map(([playerId, score], idx) => (
                                    <li
                                        key={playerId}
                                        className="flex justify-between px-4 py-3 bg-indigo-50 transition rounded-b-sm"
                                    >
                                        <span className="font-medium">
                                            #{idx + 1} {playerId}
                                        </span>
                                        <span className="text-indigo-700 font-semibold">
                                            {score} pts
                                        </span>
                                    </li>
                                ))}
                        </ul>
                    </div>
                </>
            )}
        </div>
    );
};

export default HostGame;
