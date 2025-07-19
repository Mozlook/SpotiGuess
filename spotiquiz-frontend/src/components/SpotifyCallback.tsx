import { useEffect } from "react";
import axios from "axios";

export default function SpotifyCallback() {
    useEffect(() => {
        const code = new URLSearchParams(window.location.search).get("code");

        if (code) {
            (async () => {
                try {
                    const res = await axios.post("http://localhost:8080/auth/callback", {
                        code,
                    });

                    const { access_token, spotify_id } = res.data;
                    localStorage.setItem("access_token", access_token);
                    localStorage.setItem("spotify_id", spotify_id);
                    window.location.href = "/";
                } catch (err) {
                    console.error("Auth error:", err);
                    alert("Login failed");
                }
            })();
        }
    }, []);

    return <p>Logowanie przez Spotify...</p>;
}
