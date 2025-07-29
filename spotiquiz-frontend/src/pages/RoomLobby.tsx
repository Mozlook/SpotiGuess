import { useParams, useNavigate, useLocation } from "react-router-dom";
import { useEffect, useRef } from "react";
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
        if (isHost !== true) {
            socketRef.current = new WebSocket(
                `ws://localhost:8080/ws/${code}/${playerName ? playerName : playerID}`,
            );

            socketRef.current.onopen = () => {
                console.log("WebSocket polaczony");
            };

            socketRef.current.onmessage = (event) => {
                const msg = JSON.parse(event.data);
                console.log("WS widomosc:", msg);

                if (msg.type === "game-started") {
                    navigate(`/room/${code}`, { state: playerName });
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
        }
    }, [code, playerID, navigate, isHost, playerName]);
    return (
        <div className="text-white bg-gray-900 min-h-screen flex flex-col items-center justify-center gap-4">
            <span className="text-lg">
                Room code: <strong>{code}</strong>
            </span>
            {isHost === true ? (
                <div className="flex flex-col items-center gap-2">
                    <div>You are the host</div>
                    <button
                        className="bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded"
                        onClick={StartGame}
                    >
                        Start Game
                    </button>
                </div>
            ) : (
                <div>Waiting for host to start the game...</div>
            )}
        </div>
    );
};
export default RoomLobby;
