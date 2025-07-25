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
// The request **must** include an Authorization header with a Spotify access token:
//
//	Authorization: Bearer <access_token>
//
// The handler performs the following steps:
//
//  1. Validates and decodes the JSON request body into a CreateRoomRequest struct.
//
//  2. Extracts the Spotify access token from the Authorization header.
//     - If the header is missing or invalid, responds with HTTP 401 Unauthorized.
//
//  3. Stores the access token in Redis under the key "player:{hostId}".
//
//  4. Generates a 6-character room code.
//
//  5. Constructs a new Room object with the given hostId and state set to "waiting".
//
//  6. Stores the Room in Redis under the key "room:{roomCode}" with a 60-minute TTL.
//
//  7. Responds with a JSON object containing the generated room code:
//
//     Response:
//     {
//     "RoomCode": "ABC123"
//     }
//
// On JSON parsing failure or Redis write failure, responds with an appropriate HTTP 400/500 status.
func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var request model.CreateRoomRequest
	room := new(model.Room)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	room.HostId = request.HostID

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Spotify token required", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	userKey := "player:" + room.HostId
	tokenData, _ := json.Marshal(map[string]string{
		"access_token": token,
	})
	err = store.Client.Set(store.Ctx, userKey, tokenData, 60*time.Minute).Err()
	if err != nil {
		log.Println("Failed to save token during room creation:", err)
	}

	room.Code = generateRoomCode()
	room.CreatedAt = time.Now()
	room.GameState = "waiting"

	data, _ := json.Marshal(room)
	roomKey := "room:" + room.Code
	err = store.Client.Set(store.Ctx, roomKey, data, 60*time.Minute).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{
		"RoomCode": room.Code,
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
//  1. Decodes the request body into a JoinRoomRequest struct.
//
//  2. Retrieves the Room object from Redis under "room:{roomCode}".
//     - If not found, responds with HTTP 404.
//
//  3. Check if player with the same name already exists
//     - If yes, respond with HTTP 409.
//
//  4. Appends the joining playerId to the room's Players slice.
//
//  5. Updates the room in Redis with a 60-minute TTL.
//
//  6. If the request contains a valid Authorization header:
//     - Extracts the Spotify access token.
//     - Fetches the player's 25 most recently played tracks via Spotify API.
//     - Stores the tracks in Redis under "tracks:{roomCode}:{playerId}".
//     - Also stores the access token in Redis under "player:{playerId}".
//
//  7. Responds with a JSON object confirming the join:
//
//     Response:
//     {
//     "status": "joined",
//     "roomCode": "ABC123",
//     "playerId": "spotify-user-456"
//     }
//
// Notes:
// - The Authorization header is optional, but required to store player track data.
// - If the Spotify token is expired or invalid, track saving will silently fail.
//
// On JSON parsing failure or Redis error, responds with appropriate HTTP 400/500.
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
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	var room model.Room

	err = json.Unmarshal([]byte(data), &room)
	if err != nil {
		http.Error(w, "internal error (unmarshal)", http.StatusInternalServerError)
		return
	}

	normalized := strings.ToLower(strings.TrimSpace(request.PlayerID))
	for _, player := range room.Players {
		if strings.TrimSpace(strings.ToLower(player)) == normalized {
			http.Error(w, "Player already exists", http.StatusConflict)
			return
		}
	}

	room.Players = append(room.Players, request.PlayerID)

	jsonData, _ := json.Marshal(room)
	err = store.Client.Set(store.Ctx, roomKey, string(jsonData), 60*time.Minute).Err()
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
		} else {
			log.Println("saved tracks:", len(tracks))
		}
		userKey := "player:" + request.PlayerID
		tokenData, _ := json.Marshal(map[string]string{
			"access_token": token,
		})

		err = store.Client.Set(store.Ctx, userKey, tokenData, 60*time.Minute).Err()
		if err != nil {
			log.Println("error saving user token", err)
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
