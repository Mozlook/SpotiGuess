import { useParams, useNavigate, useLocation } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import axios from "axios";

type GameMode = "players" | "playlist" | "artist";

const RoomLobby = () => {
    const { code } = useParams();
    const location = useLocation();
    const playerName = location.state;
    const isHost: boolean = localStorage.getItem("isHost") === "true";
    const playerID: string | null = localStorage.getItem("spotify_id");
    const token = localStorage.getItem("access_token");
    const url: string = import.meta.env.VITE_BACKEND_URL;
    const navigate = useNavigate();

    const socketRef = useRef<WebSocket | null>(null);

    const [playersList, setPlayersList] = useState<string[]>([]);
    const [gameMode, setGameMode] = useState<GameMode>("players");
    const [searchQuery, setSearchQuery] = useState("");
    const [searchResults, setSearchResults] = useState<
        { id: string; name: string; image: string }[]
    >([]);
    const [playlistUrl, setPlaylistUrl] = useState("");
    const [artistID, setArtistID] = useState("");

    // Start game
    const StartGame = async () => {
        try {
            const requestBody = {
                roomCode: code,
                hostId: playerID,
                gameMode: gameMode,
                tracksData: "",
            };

            if (gameMode === "playlist") requestBody.tracksData = playlistUrl;
            if (gameMode === "artist") requestBody.tracksData = artistID;

            const res = await axios.post(`${url}/start-game`, requestBody, {
                headers: {
                    ...(token && { Authorization: `Bearer ${token}` }),
                },
            });

            if (res.data.status) {
                navigate(`/room/${code}`);
            }
        } catch (err) {
            console.error(err);
        }
    };

    // WebSocket
    useEffect(() => {
        socketRef.current = new WebSocket(
            `ws://localhost:8080/ws/${code}/${playerName || playerID}`,
        );

        socketRef.current.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            if (msg.type === "game-started" && !isHost) {
                navigate(`/room/${code}`, { state: playerName });
            }
            if (msg.type === "new-player" && isHost) {
                setPlayersList((prev) => [...prev, msg.data]);
            }
        };

        return () => {
            socketRef.current?.close();
        };
    }, [code, playerID, navigate, isHost, playerName]);

    const handleSearch = async () => {
        if (!searchQuery || gameMode === "players") return;

        try {
            const res = await axios.get(`${url}/spotify/search`, {
                params: {
                    q: searchQuery,
                    type: gameMode,
                    userId: playerID,
                },
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            setSearchResults(res.data);
        } catch (err) {
            console.error("Search failed", err);
        }
    };

    return (
        <div className="min-h-screen bg-gradient-to-b from-emerald-300 via-gray-200 to-emerald-100 text-gray-800 flex flex-col items-center px-6 py-12 gap-8">
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
                        onClick={StartGame}
                        className="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-2 px-6 rounded shadow"
                    >
                        Start Game
                    </button>

                    <div className="flex flex-col gap-4 mt-6 w-full">
                        <label className="text-sm font-medium text-gray-600">
                            Select Game Mode:
                        </label>
                        <select
                            value={gameMode}
                            onChange={(e) => setGameMode(e.target.value as GameMode)}
                            className="p-2 rounded border bg-white text-gray-800"
                        >
                            <option value="players">Based on players</option>
                            <option value="playlist">From a playlist</option>
                            <option value="artist">From an artist</option>
                        </select>

                        {(gameMode === "playlist" || gameMode === "artist") && (
                            <>
                                <div className="flex gap-2">
                                    <input
                                        type="text"
                                        placeholder={`Search ${gameMode}`}
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                        className="px-4 py-2 rounded border bg-white text-gray-800 flex-1"
                                    />
                                    <button
                                        onClick={handleSearch}
                                        className="bg-indigo-500 hover:bg-indigo-600 text-white px-4 py-2 rounded"
                                    >
                                        Search
                                    </button>
                                </div>

                                {searchResults.length > 0 && (
                                    <div className="w-full bg-white rounded border p-2 max-h-64 overflow-y-auto shadow-sm">
                                        {searchResults.map((result) => (
                                            <div
                                                key={result.id}
                                                onClick={() => {
                                                    if (gameMode === "playlist") {
                                                        setPlaylistUrl(result.id);
                                                        setSearchQuery(result.name);
                                                    } else {
                                                        setArtistID(result.id);
                                                        setSearchQuery(result.name);
                                                    }
                                                    setSearchResults([]);
                                                }}
                                                className="flex items-center gap-3 p-2 hover:bg-gray-100 cursor-pointer rounded"
                                            >
                                                <img
                                                    src={result.image}
                                                    alt={result.name}
                                                    className="w-10 h-10 object-cover rounded"
                                                />
                                                <span className="text-sm">{result.name}</span>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </>
                        )}
                    </div>

                    {playersList.length > 0 && (
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
