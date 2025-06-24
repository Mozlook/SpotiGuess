package room

import (
	"backend/internal/store"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Room struct {
	Code        string
	HostId      string
	CreatedAt   time.Time
	Players     []string
	GameState   string
	CurrentQIdx int
	Scoreboard  map[string]int
}

type Player struct {
	PlayerId string
	Score    int
}

type CreateRoomRequest struct {
	HostID string `json:"hostId"`
}

type RoomResponse struct {
	RoomCode string `json:"roomCode"`
}

type JoinRoomRequest struct {
	RoomCode string `json:"roomCode"`
	PlayerID string `json:"playerId"`
}

func generateRoomCode() string {
	code := ""
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for range 6 {
		code += string(characters[rand.Intn(len(characters))])
	}

	return code
}

// CreateRoomHandler handles HTTP POST requests to /create-room.
//
// It expects a JSON payload in the following format:
//
//	{
//	  "hostId": "spotify-user-id"
//	}
//
// The handler performs the following steps:
//
//  1. Decodes the JSON request body into a CreateRoomRequest struct.
//
//  2. Generates a random 6-character room code.
//
//  3. Constructs a new Room object with initial state ("waiting").
//
//  4. Serializes the Room object to JSON.
//
//  5. Stores it in Redis with a TTL of 60 minutes.
//
//  6. Responds with a JSON object containing the generated room code:
//
//     Response:
//     {
//     "roomCode": "ABC123"
//     }
//
// In case of an error (e.g. Redis failure), responds with HTTP 500.
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateRoomRequest
	room := new(Room)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	Code := generateRoomCode()

	room.Code = Code
	room.HostId = request.HostID
	room.CreatedAt = time.Now()
	room.GameState = "waiting"

	data, _ := json.Marshal(room)
	err = store.Client.Set(store.Ctx, Code, data, 60*time.Minute).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := new(RoomResponse)
	response.RoomCode = Code
	err = json.NewEncoder(w).Encode(response)

}

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinRoomRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	data, err := store.Client.Get(store.Ctx, request.RoomCode).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var room Room

	err = json.Unmarshal([]byte(data), &room)
	if err != nil {
		http.Error(w, "internal error (unmarshal)", http.StatusInternalServerError)
		return
	}
	room.Players = append(room.Players, request.PlayerID)

	jsonData, _ := json.Marshal(room)
	err = store.Client.Set(store.Ctx, request.RoomCode, string(jsonData), 60*time.Minute).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "joined",
		"roomCode": request.RoomCode,
		"playerId": request.PlayerID,
	})
}
