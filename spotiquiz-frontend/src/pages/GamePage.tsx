import { useEffect, useState, useRef } from "react";
import { useParams, useNavigate, useLocation } from "react-router-dom";
import HostGame from "../components/HostGame.tsx";
import PlayerGame from "../components/PlayerGame.tsx";
export type Question = {
    id: string;
    trackId: string;
    trackName: string;
    options: string[];
    correct: string;
    positionMs: number;
};
const GamePage = () => {
    const navigate = useNavigate();
    const isHost = localStorage.getItem("isHost") === "true";
    const token = localStorage.getItem("access_token");
    const { code } = useParams<string>();
    const playerID: string = getPlayerId();
    const location = useLocation();
    const playerName = location.state;
    const [question, setQuestion] = useState<Question | null>(null);
    const [scoreboard, setScoreboard] = useState<Record<string, number> | null>(
        null,
    );
    const [view, setView] = useState<string>("");
    const socketRef = useRef<WebSocket | null>(null);
    const [hasAnswered, setHasAnswered] = useState<boolean>(false);
    useEffect(() => {
        if (!code || !playerID) return;

        socketRef.current = new WebSocket(
            `ws://localhost:8080/ws/${code}/${playerName ? playerName : playerID}`,
        );

        socketRef.current.onopen = () => {
            console.log("WebSocket polaczony");
        };

        socketRef.current.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            console.log("WS widomosc:", msg);

            if (msg.type === "question" && msg.data) {
                setQuestion(msg.data);
                setView("question");
                setHasAnswered(false);
            }
            if (msg.type === "game-over") {
                console.log("Navigating with:", msg.data);
                navigate("/scoreboard", { state: msg.data });
            }
            if (msg.type === "scoreboard" && msg.data) {
                setScoreboard(msg.data);
                setView("scoreboard");
            }
        };

        socketRef.current.onclose = () => {
            console.log("WebSocket rozlaczony");
        };

        return () => {
            if (socketRef.current) {
                socketRef.current.close();
            }
        };
    }, [code, playerID, navigate, scoreboard, question, playerName]);

    function getPlayerId(): string {
        return (
            localStorage.getItem("spotify_id") ||
            localStorage.getItem("guest_id") ||
            (() => {
                const id = `guest:${crypto.randomUUID()}`;
                localStorage.setItem("guest_id", id);
                return id;
            })()
        );
    }

    return (
        <div className="min-h-screen bg-gradient-to-b from-emerald-300 via-gray-200 to-emerald-100 text-gray-800 flex flex-col items-center justify-start px-4 py-8">
            <h1 className="text-4xl font-bold text-center text-indigo-800 mb-10 drop-shadow-sm">
                SpotiGuess
            </h1>

            <div className="w-full max-w-3xl flex items-center justify-center">
                {isHost ? (
                    <HostGame
                        question={question}
                        scoreboard={scoreboard}
                        view={view}
                        playerID={playerID}
                        accessToken={token}
                    />
                ) : (
                    <PlayerGame
                        question={question}
                        scoreboard={scoreboard}
                        view={view}
                        hasAnswered={hasAnswered}
                        setHasAnswered={setHasAnswered}
                        playerName={playerName}
                        code={code}
                    />
                )}
            </div>
        </div>
    );
};

export default GamePage;
