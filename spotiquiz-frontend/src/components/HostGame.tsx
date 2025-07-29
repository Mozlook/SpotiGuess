import { useEffect } from "react";
import axios from "axios";
import useSpotifyPlayer from "../hooks/useSpotifyPlayer";
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

    useEffect(() => {
        const device_id = localStorage.getItem("device_id");
        if (!device_id || !accessToken || !playerReady || !window.player) return;

        if (view === "question" && question?.trackId) {
            console.log(question.positionMs);
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
        <div className="w-full max-w-xl bg-gray-900 text-white p-6 rounded-xl shadow-md flex flex-col items-center">
            <div className="text-sm text-gray-400 mb-4">Host</div>
            <div className="text-lg font-semibold mb-6 capitalize">
                {view === "question" ? "Current Question" : "Scoreboard"}
            </div>

            {view === "question" && question && (
                <div className="text-center text-lg">
                    <p className="text-gray-300 mb-2">
                        Select correct answer on your device!
                    </p>
                    <p className="text-2xl font-bold">{question.trackName}</p>
                </div>
            )}

            {view === "scoreboard" && scoreboard && (
                <div className="w-full mt-4 space-y-2">
                    {Object.entries(scoreboard)
                        .sort(([, a], [, b]) => b - a)
                        .map(([playerId, score]) => (
                            <div
                                key={playerId}
                                className="flex justify-between px-4 py-2 bg-gray-800 rounded-lg"
                            >
                                <span className="text-gray-300">{playerId}</span>
                                <span className="font-bold text-white">{score} pkt</span>
                            </div>
                        ))}
                </div>
            )}
        </div>
    );
};

export default HostGame;
