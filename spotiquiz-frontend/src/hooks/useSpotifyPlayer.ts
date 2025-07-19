import { useEffect, useRef, useState } from "react";

declare global {
    interface Window {
        Spotify: typeof Spotify;
        onSpotifyWebPlaybackSDKReady: () => void;
        player: Spotify.Player;
    }
}

export default function useSpotifyPlayer(token: string | null) {
    const [playerReady, setPlayerReady] = useState(false);
    const playerRef = useRef<Spotify.Player | null>(null);

    useEffect(() => {
        if (!token) return;

        const script = document.createElement("script");
        script.src = "https://sdk.scdn.co/spotify-player.js";
        script.async = true;
        document.body.appendChild(script);

        window.onSpotifyWebPlaybackSDKReady = () => {
            const player = new window.Spotify.Player({
                name: "SpotiGuess Player",
                getOAuthToken: (cb) => {
                    cb(token);
                },
                volume: 0.5,
            });

            playerRef.current = player;

            window.player = player;

            player.addListener("ready", ({ device_id }: { device_id: string }) => {
                console.log("Ready with Device ID", device_id);
                localStorage.setItem("device_id", device_id);
                setPlayerReady(true);
            });

            player.addListener(
                "not_ready",
                ({ device_id }: { device_id: string }) => {
                    console.log("Device ID has gone offline", device_id);
                },
            );

            player.connect();
        };
    }, [token]);

    return { player: playerRef.current, playerReady };
}
