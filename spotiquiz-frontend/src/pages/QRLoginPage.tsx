import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const QRLoginPage: React.FC = () => {
    const { code } = useParams();
    const navigate = useNavigate();

    const apiUrl: string = import.meta.env.VITE_BACKEND_API_URL;
    const [playerName, setPlayerName] = useState<string>("");
    const [error, setError] = useState<string | null>(null);

    if (!code) {
        navigate("/");
    }
    useEffect(() => {
        const stored = localStorage.getItem("name");
        if (stored) setPlayerName(stored);
    }, []);

    const JoinRoom = async () => {
        const token = localStorage.getItem("access_token");
        if (!playerName.trim()) {
            setError("Enter name");
            return;
        }
        try {
            const res = await axios.post(
                `${apiUrl}/join-room`,
                {
                    roomCode: code,
                    playerId: playerName,
                },
                {
                    headers: {
                        ...(token && { Authorization: `Bearer ${token}` }),
                    },
                },
            );
            localStorage.setItem("name", playerName);
            localStorage.setItem("isHost", "false");
            navigate(`/room/${res.data.roomCode}/lobby`, { state: playerName });
        } catch (err) {
            if (axios.isAxiosError(err) && err.response) {
                const handledStatuses = [400, 404, 409, 500];
                const status = err.response.status;

                if (handledStatuses.includes(status)) {
                    setError(err.response.data);
                } else {
                    console.error("Unhandled error:", err);
                }
            } else {
                console.error("Unknown error:", err);
            }

            localStorage.removeItem("roomCode");
            localStorage.removeItem("isHost");
        }
    };

    return (
        <div className="min-h-screen flex flex-col items-center justify-center gap-8 px-4 py-8 sm:px-6 sm:py-12 bg-gradient-to-b from-emerald-300 via-gray-200 to-emerald-100 text-gray-800">
            <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-center">
                Welcome to SpotiGuess
            </h1>

            <div className="flex flex-col items-center gap-4 w-full max-w-sm">
                <input
                    type="text"
                    placeholder="Enter your name"
                    value={playerName}
                    onChange={(e) => setPlayerName(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                />
                <button
                    onClick={JoinRoom}
                    className="w-full bg-emerald-600 hover:bg-emerald-700 text-white font-semibold py-2 px-4 rounded-lg shadow transition-colors"
                >
                    Join
                </button>

                {error && (
                    <p className="text-red-600 font-medium text-center">{error}</p>
                )}
            </div>
        </div>
    );
};

export default QRLoginPage;
