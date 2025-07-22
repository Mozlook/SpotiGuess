import { useEffect, useState, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import HostGame from "./HostGame.tsx";
import PlayerGame from "./PlayerGame.tsx";
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
            `ws://localhost:8080/ws/${code}/${playerID}`,
        );

        socketRef.current.onopen = () => {
            console.log("WebSocket polaczony");
        };

        socketRef.current.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            console.log("WS widomosc:", msg);

            if (msg.type === "question" && msg.data) {
                console.log(msg.data);
                setQuestion(msg.data);
                console.log(question);
                setView("question");
                setHasAnswered(false);
            }
            if (msg.type === "game-over") {
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
    }, [code, playerID, navigate, scoreboard, question]);

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

    function handleAnswer(selected: string) {
        if (!question) return;

        const playerId =
            localStorage.getItem("spotify_id") || localStorage.getItem("guest_id");
        const payload = {
            type: "answer",
            data: {
                questionId: question.id,
                selected: selected,
                playerId: playerId,
            },
        };
        if (socketRef.current) {
            socketRef.current.send(JSON.stringify(payload));
            setHasAnswered(true);
        }
    }

    return (
        <div className="min-h-screen bg-gray-900 text-white flex flex-col items-center justify-start px-4 py-8">
            <h1 className="text-4xl font-bold mb-8">SpotiGuess - Gra</h1>
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
                    playerID={playerID}
                    hasAnswered={hasAnswered}
                    handleAnswer={handleAnswer}
                />
            )}
        </div>
    );
};

export default GamePage;
