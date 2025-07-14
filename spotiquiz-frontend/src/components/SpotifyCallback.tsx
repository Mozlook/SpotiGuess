import { useEffect } from "react";

export default function SpotifyCallback() {
    useEffect(() => {
        const params = new URLSearchParams(window.location.search);
        const code = params.get("code");

        if (code) {
            fetch("http://localhost:8080/auth/callback", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ code }),
            })
                .then((res) => res.json())
                .then((data) => {
                    localStorage.setItem("access_token", data.access_token);
                    localStorage.setItem("spotify_id", data.spotify_id);
                    window.location.href = "/";
                });
        }
    }, []);

    return <p>Logowanie przez Spotify...</p>;
}
