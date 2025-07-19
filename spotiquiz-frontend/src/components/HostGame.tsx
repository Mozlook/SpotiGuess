import { useEffect } from "react";
import axios from "axios";
import useSpotifyPlayer from "../hooks/useSpotifyPlayer";
import type { Question } from "./GamePage";
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
                .then(() => {
                    console.log("Track playing, will seek...");
                })
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
        <div>
            <div>Host</div>
            <div>{view}</div>
            {view === "question" && question && (
                <div>
                    <p>🎵 Select correct answer on your device!</p>
                </div>
            )}
            {view === "scoreboard" && scoreboard && (
                <div>
                    {Object.entries(scoreboard).map(([playerId, score]) => (
                        <div key={playerId}>
                            {playerId}: {score} pkt
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default HostGame;
