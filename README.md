# SpotiGuess

SpotiGuess is a real-time multiplayer music quiz game that uses players' Spotify listening history to generate personalized questions.
It is built with a Go backend and a React + Tailwind frontend.

## 🧠 How it works

- The **host** logs in via Spotify, creates a game room, and shares the join code.
- **Players** join using the code and optionally authenticate with Spotify.
- The backend fetches each user's recently played tracks.
- A quiz is generated dynamically with real or similar songs using the Last.fm API.
- Players answer questions in real-time using WebSockets.
- Scoreboards are shown after each round and at the end.

## 📦 Backend (Go)

- **Location**: `/backend`
- **Main package**: `cmd/`
- Uses:
  - Gorilla WebSocket
  - Redis for state (rooms, scores, tokens)
  - Spotify API for login & tracks
  - Last.fm API for similar tracks

### Key endpoints

- `POST /create-room` – create a new quiz room (requires Spotify token)
- `POST /join-room` – join an existing room
- `POST /start-game` – generate and launch quiz
- `POST /submit-answer` – submit answer and update score
- `GET /room/:code` – fetch room data
- `GET /room/:code/questions` – fetch questions
- `GET /room/:code/scoreboard` – get scores

## 💻 Frontend (React)

- **Location**: `/spotiquiz-frontend`
- **Stack**: Vite + React + TailwindCSS
- **Routing**: React Router
- **State**: useState/useEffect only (no Redux)

### Pages

- `/` – Home with Spotify login / Create / Join room
- `/room/:code/lobby` – Waiting room
- `/room/:code` – Game in progress
- `/scoreboard` – Final results

## 🧪 Local Development

### Backend

```bash
cd backend
go run ./cmd
```

Requires a running Redis instance and `.env` with:
- `SPOTIFY_CLIENT_ID`
- `SPOTIFY_CLIENT_SECRET`
- `LASTFM_API_KEY`

### Frontend

```bash
cd spotiquiz-frontend
npm install
npm run dev
```

App runs on [http://127.0.0.1:5173](http://127.0.0.1:5173)

## 🔗 API Integrations

- **Spotify Web API**
  - Login with OAuth
  - Fetch recently played tracks
- **Last.fm API**
  - Get similar songs for quiz generation

## 📂 Project structure

```
SpotiGuess/
├── backend/                 # Go server with REST + WebSocket API
│   ├── internal/            # Handlers, game logic, Spotify & LastFM integration
│   ├── cmd/                 # Entry point
│   └── config/, go.mod, ...
├── spotiquiz-frontend/      # React frontend with Tailwind and React Router
│   └── src/                 # Components & pages
```

## 📝 License

MIT
