import { useNavigate } from "react-router-dom";
import { useState } from "react";
import LoginPage from "./LoginPage";
import axios from "axios";
const HomePage = () => {
    const player_ID: string | null = localStorage.getItem("spotify_id");
    const url: string = import.meta.env.VITE_BACKEND_URL;
    const [roomCode, setRoomCode] = useState<string>("");
    const navigate = useNavigate();
    localStorage.removeItem("isHost");
    const CreateRoom = async () => {
        try {
            const token = localStorage.getItem("access_token");

            const res = await axios.post(
                `${url}/create-room`,
                {
                    hostId: player_ID,
                },
                {
                    headers: {
                        ...(token && { Authorization: `Bearer ${token}` }),
                    },
                },
            );
            localStorage.setItem("isHost", "true");
            navigate(`/room/${res.data.RoomCode}/lobby`);
        } catch (err) {
            console.error(err);
            alert("Couldn't create room");
            localStorage.removeItem("isHost");
        }
    };

    const JoinRoom = async () => {
        const token = localStorage.getItem("access_token");
        try {
            const res = await axios.post(
                `${url}/join-room`,
                {
                    roomCode: roomCode,
                    playerId: player_ID,
                },
                {
                    headers: {
                        ...(token && { Authorization: `Bearer ${token}` }),
                    },
                },
            );
            localStorage.setItem("isHost", "false");
            navigate(`/room/${res.data.roomCode}/lobby`);
        } catch (err) {
            console.error(err);
            alert("Couldn't join room");
            localStorage.removeItem("roomCode");
            localStorage.removeItem("isHost");
        }
    };
    const handleLogout = () => {
        localStorage.removeItem("access_token");
        localStorage.removeItem("spotify_id");
        window.location.reload();
    };

    return (
        <div className="min-h-screen flex flex-col justify-center items-center gap-8 bg-gray-900 text-white p-8">
            {player_ID ? (
                <div className="flex flex-col items-center gap-2">
                    <span className="text-lg font-medium">✅ Logged in</span>
                    <button
                        onClick={handleLogout}
                        className="bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded"
                    >
                        Logout
                    </button>
                </div>
            ) : (
                <LoginPage />
            )}

            <button
                onClick={CreateRoom}
                className="bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-4 rounded shadow"
            >
                Create room
            </button>

            <div className="flex flex-col items-center gap-2">
                <label className="text-sm font-light text-gray-300">Join room</label>
                <div className="flex gap-2">
                    <input
                        type="text"
                        value={roomCode}
                        onChange={(e) => setRoomCode(e.target.value)}
                        placeholder="ABC123"
                        className="px-4 py-2 rounded bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-600"
                    />
                    <button
                        onClick={JoinRoom}
                        className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded"
                    >
                        Join
                    </button>
                </div>
            </div>
        </div>
    );
};
export default HomePage;
