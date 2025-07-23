import { Routes, Route } from "react-router-dom";
import HomePage from "./components/HomePage";
import SpotifyCallback from "./components/SpotifyCallback";
import GamePage from "./components/GamePage";
import RoomLobby from "./components/RoomLobby";
import ScoreboardPage from "./components/ScoreboardPage";
export default function App() {
    return (
        <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/callback" element={<SpotifyCallback />} />
            <Route path="/room/:code" element={<GamePage />} />
            <Route path="/room/:code/lobby" element={<RoomLobby />} />
            <Route path="/scoreboard" element={<ScoreboardPage />} />
        </Routes>
    );
}
