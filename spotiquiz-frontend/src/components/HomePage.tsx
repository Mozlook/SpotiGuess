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
        <div className="flex flex-col justify-center items-center gap-8">
            {player_ID ? (
                <div>
                    <span>Logged in</span>
                    <button className="border" onClick={handleLogout}>
                        Logout
                    </button>
                </div>
            ) : (
                <LoginPage />
            )}
            <button className="border" onClick={CreateRoom}>
                Create room
            </button>
            <label className="border">Join room</label>
            <div className="flex gap-4">
                <input
                    type="text"
                    value={roomCode}
                    onChange={(e) => setRoomCode(e.target.value)}
                    placeholder="ABC123"
                    className="border"
                />
                <button className="border" onClick={JoinRoom}>
                    Join
                </button>
            </div>
        </div>
    );
};
export default HomePage;
