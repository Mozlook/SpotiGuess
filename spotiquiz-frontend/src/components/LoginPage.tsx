const clientId = import.meta.env.VITE_SPOTIFY_CLIENT_ID;
const redirectUri = "http://127.0.0.1:5173/callback";
const scopes = ["user-read-recently-played"];

const loginUrl = `https://accounts.spotify.com/authorize?client_id=${clientId}&response_type=code&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${scopes.join("%20")}`;

export default function LoginPage() {
    return (
        <div>
            <h1> 🎵 SpotiQuiz</h1>
            <a href={loginUrl}>
                <button>Zaloguj się przez Spotify</button>
            </a>
        </div>
    );
}
