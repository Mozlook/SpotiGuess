import { Routes, Route } from "react-router-dom";
import LoginPage from "./components/LoginPage";
import SpotifyCallback from "./components/SpotifyCallback";

export default function App() {
    return (
        <Routes>
            <Route path="/" element={<LoginPage />} />
            <Route path="/callback" element={<SpotifyCallback />} />
        </Routes>
    );
}
