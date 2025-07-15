import { useEffect, useState, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";

type Question = {
    id: string;
    trackId: string;
    trackName: string;
    options: string[];
    correct: string;
};
const GamePage = () => {
    const navigate = useNavigate();
    const { code } = useParams();
    const playerID = getPlayerId();
    const [question, setQuestion] = useState<Question | null>(null);
    const socketRef = useRef<WebSocket | null>(null);
    const [hasAnswered, setHasAnswered] = useState(false);
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
                setQuestion(msg.data);
                setHasAnswered(false);
            }
            if (msg.type === "game-over" && msg.data) {
                navigate("/scoreboard", { state: msg.data });
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
    }, [code, playerID, navigate]);

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
            type: "question",
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
        <div>
            <h1>SpotiGuess - Gra</h1>
            {question && (
                <div>
                    <h2>{question.trackName}</h2>
                    {hasAnswered && <p>Czekam na kolejne pytanie...</p>}
                    <ul>
                        {question.options.map((opt) => (
                            <li key={opt}>
                                <button
                                    onClick={() => handleAnswer(opt)}
                                    disabled={hasAnswered}
                                >
                                    {opt}
                                </button>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
            {!question && <p>Czekam na rozpoczęcie gry...</p>}
        </div>
    );
};

export default GamePage;
