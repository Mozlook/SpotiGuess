import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import LoginPage from "./LoginPage";
import axios from "axios";
const HomePage = () => {
    const player_ID: string | null = localStorage.getItem("spotify_id");
    const [playerName, setPlayerName] = useState<string>("");
    const url: string = import.meta.env.VITE_BACKEND_URL;
    const [roomCode, setRoomCode] = useState<string>("");
    const navigate = useNavigate();

    localStorage.removeItem("isHost");

    useEffect(() => {
        const token = localStorage.getItem("access_token");
        if (!token && !player_ID) return;
        const ValidateToken = async () => {
            try {
                const res = await axios.post(`${url}/auth/validate-token`, {
                    clientId: player_ID,
                    token: token,
                });
                localStorage.setItem("access_token", res.data);
            } catch (err) {
                localStorage.removeItem("spotify_id");
                localStorage.removeItem("access_token");
                console.log(err);
            }
        };
        ValidateToken();
    });
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
        if (!playerName.trim()) {
            alert("Enter Name");
            return;
        }
        try {
            const res = await axios.post(
                `${url}/join-room`,
                {
                    roomCode: roomCode,
                    playerId: playerName,
                },
                {
                    headers: {
                        ...(token && { Authorization: `Bearer ${token}` }),
                    },
                },
            );
            localStorage.setItem("isHost", "false");
            navigate(`/room/${res.data.roomCode}/lobby`, { state: playerName });
        } catch (err) {
            if (axios.isAxiosError(err) && err.response?.status === 400) {
                alert(err.response?.data);
            } else if (axios.isAxiosError(err) && err.response?.status === 404) {
                alert(err.response?.data);
            } else if (axios.isAxiosError(err) && err.response?.status === 500) {
                alert(err.response?.data);
            } else if (axios.isAxiosError(err) && err.response?.status === 409) {
                alert(err.response?.data);
            }
            console.error(err);
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
                    <span className="text-lg font-medium">Logged in</span>
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
                <input
                    type="text"
                    value={playerName}
                    onChange={(e) => setPlayerName(e.target.value)}
                    placeholder="Display name"
                    className="px-4 py-2 w-full rounded bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-600"
                />

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
