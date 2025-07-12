package room

import (
	"backend/internal/model"
	"backend/internal/spotify"
	"backend/internal/store"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

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
	var request model.CreateRoomRequest
	room := new(model.Room)
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
	roomKey := "room:" + Code
	err = store.Client.Set(store.Ctx, roomKey, data, 60*time.Minute).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{
		"RoomCode": Code,
	})

}

// JoinRoomHandler handles HTTP POST requests to /join-room.
//
// It expects a JSON payload in the following format:
//
//	{
//	  "roomCode": "ABC123",
//	  "playerId": "spotify-user-456"
//	}
//
// The handler performs the following steps:
//
//  1. Decodes the JSON request body into a JoinRoomRequest struct.
//
//  2. Retrieves the Room object from Redis using the provided room code.
//
//  3. Adds the player ID to the room's Players list.
//
//  4. Updates the room object in Redis with a new 60-minute TTL.
//
//  5. Responds with a JSON object confirming the join:
//
//     Response:
//     {
//     "status": "joined",
//     "roomCode": "ABC123",
//     "playerId": "spotify-user-456"
//     }
//
// In case of an error (e.g. room not found, JSON parsing error), responds with an appropriate HTTP status code.
func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var request model.JoinRoomRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	roomKey := "room:" + request.RoomCode
	data, err := store.Client.Get(store.Ctx, roomKey).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var room model.Room

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

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)
		tracks, err := spotify.FetchRecentTracks(token)
		if err != nil {
			panic(err)
		}
		jsonData, _ := json.Marshal(tracks)

		err = store.Client.Set(store.Ctx, "tracks:"+request.RoomCode+":"+request.PlayerID, jsonData, 60*time.Minute).Err()

		if err != nil {
			log.Println("error saving tracks:", err)
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":   "joined",
		"roomCode": request.RoomCode,
		"playerId": request.PlayerID,
	})
}

// GetRoomHandler handles HTTP GET requests to /room/{code}.
//
// It extracts the room code from the URL path, for example:
//
//	GET /room/ABC123
//
// The handler performs the following steps:
//
//  1. Parses the room code from the URL path.
//
//  2. Retrieves the Room object from Redis.
//
//  3. Deserializes the JSON data into a Room struct.
//
//  4. Responds with the full Room data in JSON format:
//
//     Response:
//     {
//     "code": "ABC123",
//     "hostId": "host123",
//     "createdAt": "...",
//     "players": ["player1", "player2"],
//     "gameState": "waiting",
//     "currentQIdx": 0,
//     "scoreboard": { ... }
//     }
//
// If the room does not exist or the URL is malformed, responds with a 404 or 500 status code.
func GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	roomKey := "room:" + parts[2]
	data, err := store.Client.Get(store.Ctx, roomKey).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var room model.Room
	err = json.Unmarshal([]byte(data), &room)
	if err != nil {
		http.Error(w, "failed to parse room data", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(room)
}
