import { useParams, useNavigate, useLocation } from "react-router-dom";
import { useEffect, useRef, useState, useCallback } from "react";
import axios from "axios";

import SourcesColumn from "@/components/lobby/SourcesColumn";
import ShareStartColumn from "@/components/lobby/ShareStartColumn";
import PlayersColumn from "@/components/lobby/PlayersColumn";

import type { GameMode, SearchItem } from "@/types/lobby";

const RoomLobby = () => {
    const { code } = useParams();
    const location = useLocation();
    const playerName = location.state as string | null;

    const isHost: boolean = localStorage.getItem("isHost") === "true";
    const playerID: string | null = localStorage.getItem("spotify_id");
    const token = localStorage.getItem("access_token");

    const apiUrl: string = import.meta.env.VITE_BACKEND_API_URL;
    const wsUrl: string = import.meta.env.VITE_BACKEND_WS_URL;
    const navigate = useNavigate();

    const socketRef = useRef<WebSocket | null>(null);

    const [playersList, setPlayersList] = useState<string[]>([]);
    const [gameMode, setGameMode] = useState<GameMode>("players");
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [searchResults, setSearchResults] = useState<SearchItem[]>([]);
    const [selectedSources, setSelectedSources] = useState<SearchItem[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<Error | null>(null);

    const getLastGame = useCallback(() => {
        const lastGameRaw = localStorage.getItem("lastGame");
        if (!lastGameRaw) return;
        try {
            const lastGame = JSON.parse(lastGameRaw) as {
                gameMode: GameMode;
                selectedSources: SearchItem[];
            };
            setGameMode(lastGame?.gameMode ?? "players");
            setSelectedSources(lastGame?.selectedSources ?? []);
            setSearchQuery("");
            setSearchResults([]);
        } catch (err) {
            console.log(err);
        }
    }, []);

    const StartGame = useCallback(async () => {
        try {
            setLoading(true);
            setError(null);

            if (
                (gameMode === "artist" || gameMode === "playlist") &&
                selectedSources.length === 0
            ) {
                setError(new Error("Select Artist or Playlist"));
                return;
            }

            const tracksData = selectedSources.map((s) => s.id);

            const requestBody = {
                roomCode: code,
                hostId: playerID,
                gameMode: gameMode,
                tracksData: tracksData,
            };

            const res = await axios.post(`${apiUrl}/start-game`, requestBody, {
                headers: {
                    ...(token && { Authorization: `Bearer ${token}` }),
                },
            });

            if (res.data.status) {
                localStorage.setItem(
                    "lastGame",
                    JSON.stringify({
                        gameMode,
                        selectedSources,
                    }),
                );
                navigate(`/room/${code}`);
            }
        } catch (err) {
            console.error(err);
            if (axios.isAxiosError(err)) {
                setError(err);
            } else {
                setError(new Error("Failed to start the game"));
            }
        } finally {
            setLoading(false);
        }
    }, [apiUrl, code, gameMode, navigate, playerID, selectedSources, token]);

    useEffect(() => {
        if (!code) return;

        socketRef.current = new WebSocket(
            `${wsUrl}/ws/${code}/${playerName || playerID}`,
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
    }, [code, playerID, isHost, playerName, wsUrl, navigate]);

    const handleSearch = useCallback(async () => {
        if (!searchQuery || gameMode === "players") return;

        try {
            const res = await axios.get(`${apiUrl}/spotify/search`, {
                params: {
                    q: searchQuery,
                    type: gameMode,
                    userId: playerID,
                },
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            setSearchResults(Array.isArray(res.data) ? res.data : []);
        } catch (err) {
            console.error("Search failed", err);
        }
    }, [apiUrl, gameMode, playerID, searchQuery, token]);

    const handleChangeMode = useCallback((mode: GameMode) => {
        setGameMode(mode);
        setSelectedSources([]);
        setSearchResults([]);
        setSearchQuery("");
    }, []);

    const handleSelectSource = useCallback((item: SearchItem) => {
        setSelectedSources((prev) =>
            prev.some((s) => s.id === item.id) ? prev : [...prev, item],
        );
        setSearchQuery("");
        setSearchResults([]);
    }, []);

    const handleRemoveSource = useCallback((id: string) => {
        setSelectedSources((prev) => prev.filter((s) => s.id !== id));
    }, []);

    return (
        <div className="h-screen bg-gradient-to-b from-emerald-300 via-gray-200 to-emerald-100 text-gray-800 flex flex-col">
            {/* Header */}
            <div className="shrink-0 px-6 pt-6 pb-4 text-center">
                <h1 className="text-3xl font-bold mb-2">Room Code</h1>
                <p className="text-lg tracking-widest font-mono bg-gray-100 text-indigo-600 px-4 py-2 rounded shadow inline-block">
                    {code}
                </p>
            </div>

            {isHost ? (
                <div className="flex-1 px-6 pb-6">
                    <div className="h-full mx-auto w-full grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 overflow-hidden">
                        {/* LEFT */}
                        <SourcesColumn
                            gameMode={gameMode}
                            searchQuery={searchQuery}
                            searchResults={searchResults}
                            selectedSources={selectedSources}
                            onChangeGameMode={handleChangeMode}
                            onChangeSearchQuery={setSearchQuery}
                            onSearch={handleSearch}
                            onSelectSource={handleSelectSource}
                            onRemoveSource={handleRemoveSource}
                        />

                        {/* MIDDLE */}
                        <ShareStartColumn
                            roomCode={code!}
                            loading={loading}
                            errorMessage={error?.message}
                            onLoadPrevious={getLastGame}
                            onStartGame={StartGame}
                        />

                        {/* RIGHT */}
                        <PlayersColumn players={playersList} />
                    </div>
                </div>
            ) : (
                <div className="flex-1 px-6 pb-6">
                    <div className="h-full mx-auto max-w-2xl bg-white/70 backdrop-blur border border-gray-200 rounded-xl shadow-md p-8 text-center flex items-center justify-center">
                        <p className="text-gray-600 text-lg font-medium italic">
                            Waiting for host to start the game…
                        </p>
                    </div>
                </div>
            )}
        </div>
    );
};

export default RoomLobby;
