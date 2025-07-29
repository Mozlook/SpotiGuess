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
        <div className="flex flex-col items-center justify-center p-6 text-white">
            <button
                onClick={() => (window.location.href = getSpotifyLoginUrl())}
                className="w-full max-w-sm flex items-center justify-center gap-3 px-4 py-2 bg-green-700 hover:bg-green-600 text-white font-medium rounded shadow-md transition duration-200"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="w-5 h-5 fill-white"
                    viewBox="0 0 168 168"
                >
                    <path d="M84,0A84,84,0,1,0,168,84,84.09,84.09,0,0,0,84,0Zm38.11,121.29a5.25,5.25,0,0,1-7.21,1.68c-19.73-12.07-44.64-14.8-73.91-8.12a5.25,5.25,0,0,1-2.32-10.26c32.24-7.29,60.2-4.1,83,9.45A5.25,5.25,0,0,1,122.11,121.29Zm10.21-20.89a6.56,6.56,0,0,1-9,2.11c-22.61-14-57.2-18.06-83.94-9.91a6.56,6.56,0,0,1-3.82-12.57c30.64-9.34,70.84-4.88,97.33,11.41A6.56,6.56,0,0,1,132.32,100.4Zm1.45-21.58c-27.59-16.38-73.29-17.86-99-9.87a7.86,7.86,0,1,1-4.78-15c29.94-9.54,81.05-7.83,113,11.11a7.86,7.86,0,0,1-8.93,13.76Z" />
                </svg>
                Login with Spotify to host game
            </button>
        </div>
    );
}
