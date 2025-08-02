import { useParams, useNavigate, useLocation } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import axios from "axios";
const RoomLobby = () => {
    const { code } = useParams();
    const location = useLocation();
    const playerName = location.state;
    const isHost: boolean = localStorage.getItem("isHost") === "true";
    const playerID: string | null = localStorage.getItem("spotify_id");
    const url: string = import.meta.env.VITE_BACKEND_URL;
    const navigate = useNavigate();
    const socketRef = useRef<WebSocket | null>(null);
    const [playersList, setPlayersList] = useState<string[]>([]);
    const StartGame = async () => {
        try {
            const token = localStorage.getItem("access_token");

            const res = await axios.post(
                `${url}/start-game`,
                {
                    roomCode: code,
                    hostId: playerID,
                },
                {
                    headers: {
                        ...(token && { Authorization: `Bearer ${token}` }),
                    },
                },
            );

            if (res.data.status) {
                navigate(`/room/${code}`);
            }
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        socketRef.current = new WebSocket(
            `ws://localhost:8080/ws/${code}/${playerName ? playerName : playerID}`,
        );

        socketRef.current.onopen = () => {
            console.log("WebSocket polaczony");
        };

        socketRef.current.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            console.log("WS widomosc:", msg);

            if (msg.type === "game-started" && !isHost) {
                navigate(`/room/${code}`, { state: playerName });
            }
            if (msg.type === "new-player" && isHost) {
                setPlayersList((playersList: string[]) => [...playersList, msg.data]);
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
    }, [code, playerID, navigate, isHost, playerName]);
    return (
        <div className="min-h-screen bg-gradient-to-b from-emerald-300 via-gray-200 to-emerald-100 text-gray-800 flex flex-col items-center justify-center px-6 py-12 gap-8">
            <div className="text-center">
                <h1 className="text-3xl font-bold mb-2">Room Code</h1>
                <p className="text-lg tracking-widest font-mono bg-gray-100 text-indigo-600 px-4 py-2 rounded shadow">
                    {code}
                </p>
            </div>

            {isHost ? (
                <div className="flex flex-col items-center gap-4 w-full max-w-md">
                    <p className="text-lg text-indigo-700 font-medium">
                        You are the host
                    </p>
                    <button
                        className="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-2 px-6 rounded shadow"
                        onClick={StartGame}
                    >
                        Start Game
                    </button>

                    {playersList && playersList.length > 0 && (
                        <div className="bg-white border border-gray-200 rounded-lg p-5 w-full shadow-md">
                            <h3 className="text-xl font-semibold text-center text-gray-700 mb-4">
                                Players in Room
                            </h3>
                            <ul className="space-y-2">
                                {playersList.map((player, index) => (
                                    <li
                                        key={index}
                                        className="flex items-center justify-between bg-gray-100 px-4 py-2 rounded-md text-sm font-medium text-gray-800"
                                    >
                                        <span className="truncate">{player}</span>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    )}
                </div>
            ) : (
                <div className="text-gray-600 text-lg font-medium italic">
                    Waiting for host to start the game...
                </div>
            )}
        </div>
    );
};
export default RoomLobby;
