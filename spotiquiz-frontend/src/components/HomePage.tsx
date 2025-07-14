import { useNavigate } from "react-router-dom";
import { useState } from "react";
import LoginPage from "./LoginPage";
import axios from "axios";
const HomePage = () => {
    const player_ID = localStorage.getItem("spotify_id");
    const url = import.meta.env.VITE_BACKEND_URL;
    const [roomCode, setRoomCode] = useState("");
    const navigate = useNavigate();

    const CreateRoom = async () => {
        const res = await axios.post(`${url}/create-room`, {
            hostId: player_ID,
        });
        localStorage.setItem("isHost", "true");
        navigate(`/room/${res.data.roomCode}/lobby`);
    };

    const JoinRoom = async () => {
        const token = localStorage.getItem("access_token");
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
    };

    return (
        <div className="flex flex-col justify-center items-center gap-8">
            {player_ID ? <span>Logged in</span> : <LoginPage />}
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
