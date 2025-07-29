import { Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import SpotifyCallback from "./components/SpotifyCallback";
import GamePage from "./pages/GamePage";
import RoomLobby from "./pages/RoomLobby";
import ScoreboardPage from "./pages/ScoreboardPage";
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
