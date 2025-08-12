<div align="center">

# SpotiGuess

### A real-time multiplayer music quiz game powered by Spotify

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)](https://reactjs.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)
[![Spotify](https://img.shields.io/badge/Spotify-1ED760?style=for-the-badge&logo=spotify&logoColor=white)](https://developer.spotify.com/)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
[![GitHub last commit](https://img.shields.io/github/last-commit/Mozlook/SpotiGuess?style=flat-square)](https://github.com/Mozlook/SpotiGuess/commits/main)

---

_Challenge your friends and discover new music together with personalized quizzes generated from your Spotify listening history!_

</div>

## Features

<div align="center">

| **Feature**                   | **Description**                                                         |
| ----------------------------- | ----------------------------------------------------------------------- |
| **Personalized Quizzes**      | Questions generated from players' Spotify listening history             |
| **Real-time Multiplayer**     | Live gameplay with instant scoring using WebSockets                     |
| **Smart Question Generation** | Uses Last.fm API to find similar tracks for challenging questions       |
| **Easy Room System**          | Host creates a room, players join with a simple code                    |
| **Optional Authentication**   | Players can join anonymously or with Spotify for personalized questions |
| **Live Scoreboards**          | Real-time scoring with round-by-round and final results                 |

</div>

## Architecture

<table align="center">
<tr>
<td width="50%">

### Backend (`/backend`)

- **Language**: Go
- **APIs**: RESTful endpoints + WebSocket
- **Database**: Redis for session management
- **External APIs**:
  - Spotify Web API (auth & tracks)
  - Last.fm API (similar tracks)
- **Key Dependencies**:
  - Gorilla WebSocket
  - Redis client

</td>
<td width="50%">

### Frontend (`/spotiquiz-frontend`)

- **Framework**: React with Vite
- **Styling**: Tailwind CSS
- **Routing**: React Router
- **State**: React hooks (useState/useEffect)
- **Real-time**: WebSocket integration
- **Build**: Lightning-fast Vite bundler

</td>
</tr>
</table>

## Getting Started

<div align="center">

### Prerequisites

![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go&logoColor=white)
![Node.js](https://img.shields.io/badge/Node.js-16+-339933?style=flat&logo=node.js&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Latest-DC382D?style=flat&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Optional-2496ED?style=flat&logo=docker&logoColor=white)

</div>

---

- Go 1.19+ installed
- Node.js 16+ and npm
- Redis server running
- Docker (optional, for containerized deployment)
- Spotify Developer Account (for API keys)
- Last.fm API account

### Spotify App Registration

Before setting up the project, you need to register your application in the Spotify Developer Dashboard:

> **Quick Setup Guide**
>
> 1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
> 2. Log in with your Spotify account
> 3. Click **"Create App"**
> 4. Fill in the app details:
>    - **App Name**: `SpotiGuess` (or your preferred name)
>    - **App Description**: `Music quiz game using Spotify data`
>    - **Redirect URI**: `http://localhost:3000/callback`
> 5. Save your `Client ID` and `Client Secret` for environment setup

### Environment Setup

<table align="center">
<tr>
<td width="50%">

#### Backend Environment

Create a `.env` file in the **backend** directory:

```bash
# Spotify Configuration
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
SPOTIFY_REDIRECT_URI=your_spotify_redirect_uri

# Last.fm Configuration
LASTFM_API_KEY=your_lastfm_api_key

# Cors Configuration
ALLOWED_CORS=your_frontend_url

# Redis Configuration
REDIS_URL=localhost:6379
```

</td>
<td width="50%">

#### Frontend Environment

Create a `.env` file in the **spotiquiz-frontend** directory:

```bash
# API Configuration
VITE_BACKEND_API_URL=http://localhost:8080
VITE_BACKEND_WS_URL=ws://localhost:8080

# Spotify Configuration
VITE_SPOTIFY_CLIENT_ID=your_spotify_client_id
VITE_SPOTIFY_SPOTIFY_URI=your_spotify_callback_uri
```

</td>
</tr>
</table>

### Backend Setup

#### Standard Setup

```bash
cd backend
go mod tidy
go run ./cmd
```

#### Docker Setup

For containerized deployment, you can use the provided Dockerfile:

```bash
cd backend
docker build -t spotiguess-backend .
docker run -p 8080:8080 --env-file .env spotiguess-backend
```

The backend server will start and handle API requests and WebSocket connections.

### Frontend Setup

```bash
cd spotiquiz-frontend
npm install
npm run dev
```

The frontend will be available at `http://127.0.0.1:5173`

## How to Play

<table align="center">
<tr>
<td width="25%" align="center">
<h3>Host Setup</h3>
<p>Log in with Spotify<br/>Create game room<br/>Share join code</p>
</td>
<td width="25%" align="center">
<h3>Player Join</h3>
<p>Enter room code<br/>Optional Spotify auth<br/>Wait in lobby</p>
</td>
<td width="25%" align="center">
<h3>Gameplay</h3>
<p>Answer music questions<br/>Real-time competition<br/>Multiple rounds</p>
</td>
<td width="25%" align="center">
<h3>Results</h3>
<p>Live scoreboards<br/>Round summaries<br/>Final rankings</p>
</td>
</tr>
</table>

## API Endpoints

### Room Management

- `POST /create-room` - Create a new quiz room (requires Spotify token)
- `POST /join-room` - Join an existing room with code
- `GET /room/:code` - Fetch room information
- `GET /room/:code/questions` - Get quiz questions for room
- `GET /room/:code/scoreboard` - Retrieve current scores

### Game Flow

- `POST /start-game` - Generate quiz and launch game
- `POST /submit-answer` - Submit player answer and update score

## Project Structure

```
SpotiGuess/
├── backend/                 # Go server application
│   ├── cmd/                # Application entry point
│   ├── internal/           # Internal packages
│   │   ├── handlers/       # HTTP & WebSocket handlers
│   │   ├── game/          # Game logic and quiz generation
│   │   ├── spotify/       # Spotify API integration
│   │   └── lastfm/        # Last.fm API integration
│   ├── config/            # Configuration management
│   ├── Dockerfile         # Docker containerization
│   ├── .env               # Backend environment variables
│   └── go.mod             # Go dependencies
├── spotiquez-frontend/     # React frontend application
│   ├── src/               # Source code
│   │   ├── components/    # Reusable UI components
│   │   ├── pages/         # Route components
│   │   └── utils/         # Utility functions
│   ├── public/            # Static assets
│   ├── .env               # Frontend environment variables
│   └── package.json       # Node.js dependencies
└── README.md              # Project documentation
```

## Game Flow

1. **Room Creation**: Host authenticates with Spotify and creates room
2. **Player Joining**: Players enter room code and optionally authenticate
3. **Data Fetching**: Backend retrieves recently played tracks for authenticated users
4. **Quiz Generation**: System creates questions using real tracks and similar alternatives via Last.fm
5. **Real-time Gameplay**: Players answer questions with live score updates via WebSockets
6. **Scoring**: Points awarded for correct answers with time bonuses
7. **Results**: Round and final scoreboards displayed

## Key Technologies

- **Go**: High-performance backend with excellent concurrency
- **Docker**: Containerization for easy deployment and scalability
- **React**: Modern UI with component-based architecture
- **WebSockets**: Real-time bidirectional communication
- **Redis**: Fast in-memory storage for game state
- **Spotify API**: Music data and user authentication
- **Last.fm API**: Music discovery and recommendations
- **Tailwind CSS**: Utility-first styling framework

<div align="center">

---

## Contributing

We welcome contributions! Please feel free to submit issues and pull requests.

[![Contributors](https://img.shields.io/github/contributors/Mozlook/SpotiGuess?style=flat-square)](https://github.com/Mozlook/SpotiGuess/graphs/contributors)
[![Issues](https://img.shields.io/github/issues/Mozlook/SpotiGuess?style=flat-square)](https://github.com/Mozlook/SpotiGuess/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/Mozlook/SpotiGuess?style=flat-square)](https://github.com/Mozlook/SpotiGuess/pulls)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

### Enjoy the Music!

_SpotiGuess combines the joy of music discovery with competitive gaming. Challenge your friends, discover new tracks, and see who really knows their music best!_

<img src="https://img.shields.io/badge/Made_with-Love-red?style=for-the-badge"/>
<img src="https://img.shields.io/badge/Powered_by-Spotify-1ED760?style=for-the-badge&logo=spotify&logoColor=white"/>

</div>
