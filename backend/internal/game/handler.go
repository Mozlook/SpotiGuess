package game

import (
	"backend/internal/model"
	"backend/internal/store"
	"backend/internal/ws"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// StartGameHandler handles HTTP POST requests to /start-game.
//
// It expects a JSON payload in the following format:
//
//	{
//	  "roomCode": "ABC123",
//	  "hostId": "spotify-user-456"
//	}
//
// The handler performs the following steps:
//
//  1. Decodes the JSON request body into a StartGameRequest struct.
//  2. Retrieves the Room object from Redis using key "room:{roomCode}".
//  3. Verifies that the requesting user (hostId) matches the room's HostId.
//  4. Iterates over all players in the room and attempts to fetch their saved tracks
//     from Redis under the key "tracks:{roomCode}:{playerId}".
//     - Invalid or missing track data is logged and skipped.
//  5. Combines all retrieved tracks, shuffles them, and selects the first 10 (or fewer).
//  6. Calls GenerateQuestions with the selected tracks to create quiz questions.
//  7. Stores the generated []Question in Redis under key "questions:{roomCode}" with a TTL of 60 minutes.
//  8. Launches the quiz loop asynchronously via RunQuizLoop(roomCode).
//  9. Broadcasts a "game-started" message via WebSocket to all clients in the room.
//
// 10. Responds with a JSON object containing:
//
//	Response:
//	{
//	  "status": "started",
//	  "questionsCount": 10
//	}
//
// If the host is invalid, Redis access fails, or question generation fails,
// the handler responds with an appropriate HTTP error (e.g. 400, 403, 500).
func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	var request model.StartGameRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	roomKey := "room:" + request.RoomCode
	data, err := store.Client.Get(store.Ctx, roomKey).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var room model.Room
	err = json.Unmarshal([]byte(data), &room)
	if err != nil {
		http.Error(w, "Invalid room object", http.StatusInternalServerError)
		return
	}
	if request.HostId != room.HostId {
		http.Error(w, "Invalid HostId", http.StatusForbidden)
		return
	}
	var allTracks []model.Track
	for _, playerID := range room.Players {
		key := fmt.Sprintf("tracks:%s:%s", request.RoomCode, playerID)
		raw, err := store.Client.Get(store.Ctx, key).Result()
		if err != nil {
			log.Println("no tracks for player:", playerID, err)
			continue
		}

		var tracks []model.Track
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
	var selectedTracks []model.Track
	if len(allTracks) < 10 {
		selectedTracks = allTracks
	} else {
		selectedTracks = allTracks[:10]
	}

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

	go RunQuizLoop(request.RoomCode)
	message := map[string]any{
		"type": "game-started",
	}

	payload, _ := json.Marshal(message)
	ws.GlobalHub.Broadcast <- ws.BroadcastMessage{
		RoomCode: request.RoomCode,
		Data:     payload,
	}
	json.NewEncoder(w).Encode(map[string]any{
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
	var questions []model.Question
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
	var request model.AnswerRequest
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

	var questions []model.Question
	err = json.Unmarshal([]byte(data), &questions)
	if err != nil {
		http.Error(w, "Invalid questions data", http.StatusInternalServerError)
		return
	}
	var question model.Question
	found := false
	for _, q := range questions {
		if q.ID == request.QuestionID {
			question = q
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}
	scoreKey := fmt.Sprintf("score:%s:%s", request.RoomCode, request.PlayerID)
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
	json.NewEncoder(w).Encode(map[string]any{
		"correct": request.Selected == question.CorrectAnswer,
		"score":   currentScore,
	})
}

// GetScoreboardHandler handles HTTP GET requests to /room/{code}/scoreboard.
//
// It expects the room code to be embedded in the URL path, e.g.:
//
//	GET /room/ABC123/scoreboard
//
// The handler performs the following steps:
//
//  1. Parses the room code from the URL.
//
//  2. Retrieves the corresponding Room object from Redis (key: "room:{roomCode}").
//
//  3. Iterates through all players in the room.
//
//  4. For each player, attempts to retrieve their current score from Redis
//     under the key "score:{roomCode}:{playerId}". If no score is found or parsing fails,
//     the player is assumed to have a score of 0.
//
//  5. Builds a scoreboard as a map of player IDs to scores.
//
//  6. Responds with the full scoreboard as a JSON object:
//
//     Response:
//     {
//     "scoreboard": {
//     "spotify-user-1": 2000,
//     "spotify-user-2": 1000,
//     "guest123": 0
//     }
//     }
//
// In case of an error (e.g. room not found or Redis failure),
// responds with the appropriate HTTP error status.
func GetScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	var room model.Room
	roomCode := parts[2]
	data, err := store.Client.Get(store.Ctx, "room:"+roomCode).Result()
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	err = json.Unmarshal([]byte(data), &room)
	if err != nil {
		http.Error(w, "Failed to parse room", http.StatusInternalServerError)
		return
	}

	scoreboard := make(map[string]int)
	for _, player := range room.Players {

		scoreKey := fmt.Sprintf("score:%s:%s", room.Code, player)
		scoreData, err := store.Client.Get(store.Ctx, scoreKey).Result()
		if err != nil {
			log.Printf("Failed to fetch or parse score for player %s: %v", player, err)

		}

		score := 0
		if err == nil {

			score, err = strconv.Atoi(scoreData)
			if err != nil {

				log.Println("Failed to update score:", err)
			}
		}

		scoreboard[player] = score
	}
	json.NewEncoder(w).Encode(map[string]any{
		"scoreboard": scoreboard,
	})

}

// GetNextQuestionHandler handles HTTP GET requests to /room/{code}/next-question.
//
// It expects the room code to be embedded in the URL path, for example:
//
//	GET /room/ABC123/next-question
//
// The handler performs the following steps:
//
//  1. Parses the room code from the URL.
//
//  2. Retrieves the Room object from Redis (key: "room:{roomCode}").
//
//  3. Checks the CurrentQIdx field of the room to determine which question is next.
//
//  4. Retrieves the full list of questions from Redis (key: "questions:{roomCode}").
//
//  5. If there are no more questions (i.e. CurrentQIdx >= len(questions)),
//     responds with HTTP 204 No Content to indicate the end of the quiz.
//
//  6. Otherwise:
//     - Increments the CurrentQIdx by 1,
//     - Updates the Room in Redis,
//     - Returns the next question and its index.
//
//     Example response:
//
//     {
//     "question": { ... },
//     "index": 2,
//     "total": 10,
//     }
//
// In case of errors (e.g. invalid room code, Redis error, parse failure),
// responds with the appropriate HTTP error status.
func GetNextQuestionHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	roomCode := parts[2]

	var room model.Room
	data, err := store.Client.Get(store.Ctx, "room:"+roomCode).Result()
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	err = json.Unmarshal([]byte(data), &room)
	if err != nil {
		http.Error(w, "Failed to parse room", http.StatusInternalServerError)
		return
	}

	currentQuestionIdx := room.CurrentQIdx
	room.CurrentQIdx++
	updatedRoomData, _ := json.Marshal(room)
	store.Client.Set(store.Ctx, "room:"+roomCode, updatedRoomData, 60*time.Minute)

	var questions []model.Question
	data, err = store.Client.Get(store.Ctx, "questions:"+roomCode).Result()
	if err != nil {
		http.Error(w, "Questions not found", http.StatusNotFound)
		return
	}

	err = json.Unmarshal([]byte(data), &questions)
	if err != nil {
		http.Error(w, "Failed to parse questions", http.StatusInternalServerError)
		return
	}
	if currentQuestionIdx >= len(questions) {
		scoreboard := make(map[string]int)
		for _, player := range room.Players {
			data, err := store.Client.Get(store.Ctx, "score:"+roomCode+":"+player).Result()
			if err != nil {
				log.Printf("Failed to fetch or parse score for player %s: %v", player, err)
			}

			score := 0
			if err == nil {

				score, err = strconv.Atoi(data)
				if err != nil {

					log.Println("Failed to update score:", err)
				}
			}

			scoreboard[player] = score
		}

		message := map[string]any{
			"type": "game-over",
			"data": scoreboard,
		}
		payload, _ := json.Marshal(message)
		ws.GlobalHub.Broadcast <- ws.BroadcastMessage{
			RoomCode: roomCode,
			Data:     payload,
		}

		store.Client.Del(store.Ctx, "room:"+roomCode)
		store.Client.Del(store.Ctx, "questions:"+roomCode)
		for _, player := range room.Players {
			store.Client.Del(store.Ctx, "score:"+roomCode+":"+player)
			store.Client.Del(store.Ctx, "tracks:"+roomCode+":"+player)

		}
		return

	}

	currentQuestion := questions[currentQuestionIdx]
	message := map[string]any{
		"type": "question",
		"data": currentQuestion,
	}
	payload, _ := json.Marshal(message)

	ws.GlobalHub.Broadcast <- ws.BroadcastMessage{
		RoomCode: roomCode,
		Data:     payload,
	}

	json.NewEncoder(w).Encode(map[string]any{
		"question": currentQuestion,
		"index":    currentQuestionIdx,
		"total":    len(questions),
	})

}
