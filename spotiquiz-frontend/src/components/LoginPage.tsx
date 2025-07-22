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
        <div className="flex flex-col items-center justify-center p-8 text-white bg-gray-900 rounded-md shadow-md">
            <h1 className="text-3xl font-bold mb-6">🎵 SpotiQuiz</h1>
            <button
                onClick={() => (window.location.href = getSpotifyLoginUrl())}
                className="bg-green-600 hover:bg-green-700 text-white py-2 px-6 rounded shadow"
            >
                Zaloguj się przez Spotify
            </button>
        </div>
    );
}
