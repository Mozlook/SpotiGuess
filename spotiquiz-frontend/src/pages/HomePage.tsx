import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import LoginPage from "../components/LoginPage";
import axios from "axios";
import CustomDialog from "../components/CustomDialog";
import CustomAlert from "../components/CustomAlert";
const HomePage = () => {
    const player_ID: string | null = localStorage.getItem("spotify_id");
    const apiUrl: string = import.meta.env.VITE_BACKEND_API_URL;
    const [roomCode, setRoomCode] = useState<string>("");
    const navigate = useNavigate();
    const [error, setError] = useState<string | null>(null);
    const [errorTitle, setErrorTitle] = useState<
        string | number | null | undefined
    >(null);
    localStorage.removeItem("isHost");

    useEffect(() => {
        const token = localStorage.getItem("access_token");
        if (!token && !player_ID) return;
        const ValidateToken = async () => {
            try {
                const res = await axios.post(`${apiUrl}/auth/validate-token`, {
                    clientId: player_ID,
                    token: token,
                });
                localStorage.setItem("access_token", res.data.access_token);
            } catch (err) {
                localStorage.removeItem("spotify_id");
                localStorage.removeItem("access_token");
                console.log(err);
                if (axios.isAxiosError(err)) {
                    setError(err.response?.data);
                }
            }
        };
        ValidateToken();
    }, []);

    useEffect(() => {
        if (!error) return;
        const timer = setTimeout(() => setError(null), 4000);
        return () => clearTimeout(timer);
    }, [error]);

    const CreateRoom = async () => {
        try {
            const token = localStorage.getItem("access_token");

            const res = await axios.post(
                `${apiUrl}/create-room`,
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
            localStorage.removeItem("isHost");
            if (axios.isAxiosError(err)) {
                setErrorTitle(err.response?.status);
                setError(err.response?.data);
            }
        }
    };

    const JoinRoom = async (name: string) => {
        const token = localStorage.getItem("access_token");
        if (!name.trim()) {
            setError("Enter Name");
            setErrorTitle("Error code: 400");
            return;
        }
        try {
            const res = await axios.post(
                `${apiUrl}/join-room`,
                {
                    roomCode: roomCode,
                    playerId: name,
                },
                {
                    headers: {
                        ...(token && { Authorization: `Bearer ${token}` }),
                    },
                },
            );
            localStorage.setItem("name", name);
            localStorage.setItem("isHost", "false");
            navigate(`/room/${res.data.roomCode}/lobby`, { state: name });
        } catch (err) {
            if (axios.isAxiosError(err) && err.response?.status === 400) {
                setError(err.response?.data);
                setErrorTitle("Error code: 400");
            } else if (axios.isAxiosError(err) && err.response?.status === 404) {
                setError(err.response?.data);
                setErrorTitle("Error code: 404");
            } else if (axios.isAxiosError(err) && err.response?.status === 500) {
                setError(err.response?.data);
                setErrorTitle("Error code: 500");
            } else if (axios.isAxiosError(err) && err.response?.status === 409) {
                setError(err.response?.data);
                setErrorTitle("Error code: 409");
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
        <div className="min-h-screen flex flex-col items-center justify-center gap-8 px-4 py-8 sm:px-6 sm:py-12 bg-gradient-to-b from-emerald-300 via-gray-200 to-emerald-100 text-gray-800">
            {error && <CustomAlert msg={error} title={errorTitle} />}

            <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-center">
                Welcome to SpotiGuess
            </h1>

            {player_ID && (
                <div className="flex flex-col items-center gap-3 text-center">
                    <span className="text-base sm:text-lg font-medium text-green-600">
                        Logged in as {player_ID}
                    </span>
                    <button
                        onClick={CreateRoom}
                        className="w-48 sm:w-56 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-2 px-4 rounded shadow"
                    >
                        Create room
                    </button>
                </div>
            )}

            <div className="flex flex-col items-center gap-3 w-full max-w-sm sm:max-w-md text-center">
                <label className="text-sm text-gray-600">Join a room</label>
                <div className="flex gap-2">
                    <input
                        type="text"
                        value={roomCode}
                        onChange={(e) => setRoomCode(e.target.value)}
                        placeholder="Room code"
                        className="w-full px-4 py-2 rounded bg-white text-gray-800 border border-gray-300 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                    <CustomDialog onConfirm={JoinRoom} roomCode={roomCode} />
                </div>
            </div>

            <div className="mt-6">
                {player_ID ? (
                    <button
                        onClick={handleLogout}
                        className="bg-rose-500 hover:bg-rose-600 text-white py-2 px-4 rounded shadow"
                    >
                        Logout
                    </button>
                ) : (
                    <LoginPage />
                )}
            </div>
        </div>
    );
};
export default HomePage;
