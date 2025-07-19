const clientId = import.meta.env.VITE_SPOTIFY_CLIENT_ID;
const redirectUri = "http://127.0.0.1:5173/callback";
const scopes = [
    "user-read-recently-played",
    "user-read-private",
    "user-read-email",
    "streaming",
];

function getSpotifyLoginUrl() {
    const params = new URLSearchParams({
        client_id: clientId,
        response_type: "code",
        redirect_uri: redirectUri,
        scope: scopes.join(" "),
    });

    return `https://accounts.spotify.com/authorize?${params.toString()}`;
}

export default function LoginPage() {
    return (
        <div>
            <h1> 🎵 SpotiQuiz</h1>
            <button onClick={() => (window.location.href = getSpotifyLoginUrl())}>
                Zaloguj się przez Spotify
            </button>
        </div>
    );
}
