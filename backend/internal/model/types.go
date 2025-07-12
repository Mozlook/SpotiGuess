package model

import "time"

// Question represents a single quiz question.
type Question struct {
	ID            string   `json:"id"`
	TrackID       string   `json:"trackId"`
	TrackName     string   `json:"trackName"`
	AnswerOptions []string `json:"options"`
	CorrectAnswer string   `json:"correct"`
}

// Track represents a simplified track structure fetched from Spotify.
type Track struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []string `json:"artists"`
}

// Room holds the state of a quiz room.
type Room struct {
	Code        string         `json:"code"`
	HostId      string         `json:"hostId"`
	CreatedAt   time.Time      `json:"createdAt"`
	Players     []string       `json:"players"`
	GameState   string         `json:"gameState"`
	CurrentQIdx int            `json:"currentQIdx"`
	Scoreboard  map[string]int `json:"scoreboard"`
}

// CreateRoomRequest is the request body for /create-room.
type CreateRoomRequest struct {
	HostID string `json:"hostId"`
}

// RoomResponse is returned when creating a room.
type RoomResponse struct {
	RoomCode string `json:"roomCode"`
}

// JoinRoomRequest is the request body for /join-room.
type JoinRoomRequest struct {
	RoomCode string `json:"roomCode"`
	PlayerID string `json:"playerId"`
}

// StartGameRequest is the request body for /start-game.
type StartGameRequest struct {
	RoomCode string `json:"roomCode"`
	HostId   string `json:"hostId"`
}

// AnswerRequest is the request body for /submit-answer.
type AnswerRequest struct {
	RoomCode   string `json:"roomCode"`
	QuestionID string `json:"questionId"`
	Selected   string `json:"selected"`
	PlayerID   string `json:"playerId"`
}

// ScoreEntry can be used for sorting or returning top players.
type ScoreEntry struct {
	PlayerID string `json:"playerId"`
	Score    int    `json:"score"`
}
