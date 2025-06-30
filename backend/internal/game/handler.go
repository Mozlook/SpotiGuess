package game

import (
	"backend/internal/room"
	"backend/internal/spotify"
	"backend/internal/store"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type StartGameRequest struct {
	RoomCode string `json:"roomCode"`
	HostId   string `json:"hostId"`
}

// StartGameHandler handles HTTP POST requests to /start-game.
//
// It expects a JSON payload in the following format:
//
//	{
//	  "roomCode": "ABC123",
//	  "hostId": "spotify-user-id"
//	}
//
// The handler performs the following steps:
//
//  1. Validates the request and ensures that the provided host ID matches the host of the room.
//
//  2. Retrieves the room from Redis using the given room code.
//
//  3. For each player in the room, attempts to fetch their previously stored track history
//     from Redis under the key "tracks:{roomCode}:{playerId}". Invalid or missing data is skipped.
//
//  4. Combines all retrieved tracks, shuffles the list, and selects up to 10 unique tracks.
//
//  5. Retrieves the host's Spotify access token from Redis (under "user:{hostId}").
//
//  6. Calls GenerateQuestions with the selected tracks and token to generate quiz questions.
//
//  7. Stores the generated []Question into Redis under the key "questions:{roomCode}"
//     with a 60-minute TTL.
//
//  8. Responds with a JSON object indicating the game has started:
//
//     Response:
//     {
//     "status": "started",
//     "questionsCount": 10
//     }
//
// In case of any failure (invalid input, Redis error, token retrieval failure, question generation failure),
// an appropriate HTTP error is returned.
func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	var request StartGameRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	data, err := store.Client.Get(store.Ctx, request.RoomCode).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var room room.Room
	err = json.Unmarshal([]byte(data), &room)
	if err != nil {
		http.Error(w, "Invalid room object", http.StatusInternalServerError)
		return
	}
	if request.HostId != room.HostId {
		http.Error(w, "Invalid HostId", http.StatusForbidden)
		return
	}
	var allTracks []spotify.Track
	for _, playerID := range room.Players {
		key := fmt.Sprintf("tracks:%s:%s", request.RoomCode, playerID)
		raw, err := store.Client.Get(store.Ctx, key).Result()
		if err != nil {
			log.Println("no tracks for player:", playerID, err)
			continue
		}

		var tracks []spotify.Track
		err = json.Unmarshal([]byte(raw), &tracks)
		if err != nil {
			log.Println("invalid track data for player:", playerID, err)
			continue
		}

		allTracks = append(allTracks, tracks...)
	}

	rand.Shuffle(len(allTracks), func(i, j int) {
		allTracks[i], allTracks[j] = allTracks[j], allTracks[i]
	})
	var selectedTracks []spotify.Track
	if len(allTracks) < 10 {
		selectedTracks = allTracks
	} else {
		selectedTracks = allTracks[:10]
	}

	key := "user:" + request.HostId
	raw, err := store.Client.Get(store.Ctx, key).Result()
	var tokenData struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal([]byte(raw), &tokenData)
	if err != nil {
		http.Error(w, "Failed to parse token", http.StatusInternalServerError)
		return
	}
	token := tokenData.AccessToken

	questions, err := GenerateQuestions(selectedTracks, token)
	if err != nil {
		http.Error(w, "Failed to generate questions", http.StatusInternalServerError)
		return
	}

	jsonData, _ := json.Marshal(questions)
	err = store.Client.Set(store.Ctx, "questions:"+request.RoomCode, jsonData, 60*time.Minute).Err()
	if err != nil {
		http.Error(w, "Failed to save questions", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "started",
		"questionsCount": len(questions),
	})

}
