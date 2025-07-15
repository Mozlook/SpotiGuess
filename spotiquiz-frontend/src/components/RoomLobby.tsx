import { useParams, useNavigate } from "react-router-dom";
import { useEffect, useRef } from "react";
import axios from "axios";
const RoomLobby = () => {
    const { code } = useParams();
    const isHost = localStorage.getItem("isHost");
    const playerID = localStorage.getItem("spotify_id");
    const url = import.meta.env.VITE_BACKEND_URL;
    const navigate = useNavigate();
    const socketRef = useRef<WebSocket | null>(null);
    const StartGame = async () => {
        try {
            const res = await axios.post(`${url}/start-game`, {
                roomCode: code,
                hostId: playerID,
            });
            if (res.data.status) {
                navigate(`/room/${code}`);
            }
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => {
        if (isHost !== "true") {
            socketRef.current = new WebSocket(
                `ws://localhost:8080/ws/${code}/${playerID}`,
            );

            socketRef.current.onopen = () => {
                console.log("WebSocket polaczony");
            };

            socketRef.current.onmessage = (event) => {
                const msg = JSON.parse(event.data);
                console.log("WS widomosc:", msg);

                if (msg.type === "game-started") {
                    navigate(`/room/${code}`);
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
    }, [code, playerID, navigate, isHost]);
    return (
        <>
            <span>Room code: {code}</span>
            {isHost === "true" ? (
                <div>
                    <div>You are a host</div>
                    <button onClick={StartGame}>Start Game</button>
                </div>
            ) : (
                <div>Waiting for host</div>
            )}
        </>
    );
};
export default RoomLobby;
