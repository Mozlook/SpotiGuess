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
	"strconv"
	"strings"
	"time"
)

type StartGameRequest struct {
	RoomCode string `json:"roomCode"`
	HostId   string `json:"hostId"`
}

type AnswerRequest struct {
	RoomCode   string `json:"roomCode"`
	QuestionId string `json:"questionId"`
	Selected   string `json:"selected"`
	PlayerId   string `json:"playerId"`
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
	roomKey := "room:" + request.RoomCode
	data, err := store.Client.Get(store.Ctx, roomKey).Result()
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

// GetQuestionsHandler handles HTTP GET requests to /room/{code}/questions.
//
// It expects the room code to be embedded in the URL path as the third segment,
// followed by the "questions" keyword, e.g.:
//
//	GET /room/ABC123/questions
//
// The handler performs the following steps:
//
//  1. Parses the room code from the URL.
//
//  2. Retrieves the list of quiz questions for that room from Redis,
//     stored under the key "questions:{roomCode}".
//
//  3. Deserializes the stored JSON into a slice of Question structs.
//
//  4. Responds with the full list of questions as a JSON array:
//
//     Response:
//     [
//     { "id": "q1", "trackId": "...", "trackName": "...", "options": [...], "correct": "..." },
//     ...
//     ]
//
// If the room or questions cannot be found, responds with HTTP 500.
func GetQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	key := "questions:" + parts[2]
	raw, err := store.Client.Get(store.Ctx, key).Result()
	if err != nil {
		http.Error(w, "Failed to get questions", http.StatusInternalServerError)
		return
	}
	var questions []Question
	err = json.Unmarshal([]byte(raw), &questions)
	if err != nil {
		http.Error(w, "Failed to parse questions", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(questions)
}

// SubmitAnswerHandler handles HTTP POST requests to /submit-answer.
//
// It expects a JSON payload in the following format:
//
//	{
//	  "roomCode": "ABC123",
//	  "questionId": "q1",
//	  "selected": "Photograph",
//	  "playerId": "spotify-user-456"
//	}
//
// The handler performs the following steps:
//
//  1. Retrieves the question list for the specified room from Redis ("questions:{roomCode}").
//
//  2. Locates the question matching the given questionId.
//
//  3. Compares the player's selected answer to the correct answer.
//
//  4. Retrieves the player's current score from Redis (under "score:{roomCode}:{playerId}"),
//     or initializes it to 0 if not found.
//
//  5. If the answer is correct, adds 1000 points to the player's score and updates the Redis entry.
//
//  6. Responds with a JSON object indicating whether the answer was correct and the player's updated score:
//
//     Response:
//     {
//     "correct": true,
//     "score": 2000
//     }
//
// In case of an error (e.g. invalid request, question not found, Redis error), responds with the appropriate HTTP status code.
func SubmitAnswerHandler(w http.ResponseWriter, r *http.Request) {
	var request AnswerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	key := "questions:" + request.RoomCode
	data, err := store.Client.Get(store.Ctx, key).Result()
	if err != nil {
		http.Error(w, "Failed to get questions", http.StatusInternalServerError)
		return
	}

	var questions []Question
	err = json.Unmarshal([]byte(data), &questions)
	if err != nil {
		http.Error(w, "Invalid questions data", http.StatusInternalServerError)
		return
	}
	var question Question
	found := false
	for _, q := range questions {
		if q.ID == request.QuestionId {
			question = q
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}
	scoreKey := fmt.Sprintf("score:%s:%s", request.RoomCode, request.PlayerId)
	currentScore := 0

	rawScore, err := store.Client.Get(store.Ctx, scoreKey).Result()
	if err == nil {
		currentScore, err = strconv.Atoi(rawScore)
		if err != nil {
			log.Println("Invalid score in Redis, resetting to 0")
			currentScore = 0
		}
	}

	if request.Selected == question.CorrectAnswer {
		currentScore += 1000
		store.Client.Set(store.Ctx, scoreKey, currentScore, 60*time.Minute)
		if err != nil {
			log.Println("Failed to update score:", err)
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"correct": request.Selected == question.CorrectAnswer,
		"score":   currentScore,
	})
}
