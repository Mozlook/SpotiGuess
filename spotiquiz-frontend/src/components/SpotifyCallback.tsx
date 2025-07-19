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
                .then(async (res) => {
                    if (!res.ok) {
                        const text = await res.text(); // odczytaj jako tekst (bo może to nie być JSON)
                        throw new Error(`Auth failed: ${res.status} - ${text}`);
                    }
                    return res.json();
                })
                .then((data) => {
                    localStorage.setItem("access_token", data.access_token);
                    localStorage.setItem("spotify_id", data.spotify_id);
                    window.location.href = "/";
                })
                .catch((err) => {
                    console.error("Auth error:", err);
                    alert("Logowanie nie powiodło się");
                });
        }
    }, []);

    return <p>Logowanie przez Spotify...</p>;
}
