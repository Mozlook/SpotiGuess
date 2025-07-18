package game

import (
	"backend/internal/model"
	"backend/internal/store"
	"backend/internal/ws"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

// RunQuizLoop runs the automated quiz loop for a given room.
//
// This function is executed asynchronously (as a goroutine) after /start-game is called.
// It handles the full lifecycle of the quiz by broadcasting questions and scoreboards
// via WebSocket to all clients in the room, with delays between rounds.
//
// It performs the following steps:
//
//  1. Waits 2 seconds before starting (to ensure clients are ready).
//
//  2. Retrieves the Room object from Redis ("room:{roomCode}").
//
//  3. Retrieves the generated []Question from Redis ("questions:{roomCode}").
//
//  4. Iterates over each question:
//     a. Broadcasts a WebSocket message:
//
//     {
//     "type": "question",
//     "data": { ...question }
//     }
//
//     b. Increments and updates the room's CurrentQIdx in Redis.
//     c. Waits x seconds for players to answer.
//     d. Gathers scores for each player from Redis ("score:{roomCode}:{playerId}").
//     e. Broadcasts the scoreboard:
//
//     {
//     "type": "scoreboard",
//     "data": { "player1": 2000, "guest:xyz": 1000 }
//     }
//
//     f. Waits another x seconds before continuing.
//
//  5. After all questions, broadcasts a final message:
//
//     {
//     "type": "game-over"
//     }
//
//  6. Cleans up by deleting all room-related keys from Redis:
//     - "room:{roomCode}"
//     - "questions:{roomCode}"
//     - "tracks:{roomCode}:{playerId}"
//     - "score:{roomCode}:{playerId}"
//
// This function assumes all track, question, and score data exists and is valid.
// It logs any Redis or decode failures and continues where possible.
func RunQuizLoop(roomCode string) {
	time.Sleep(2 * time.Second)
	log.Println("Starting quiz loop for room:", roomCode)

	roomKey := "room:" + roomCode
	data, err := store.Client.Get(store.Ctx, roomKey).Result()
	if err != nil {
		panic(err)
	}

	var room model.Room
	err = json.Unmarshal([]byte(data), &room)
	if err != nil {
		panic(err)
	}

	questionsKey := "questions:" + roomCode
	data, err = store.Client.Get(store.Ctx, questionsKey).Result()
	if err != nil {
		panic(err)
	}

	var questions []model.Question
	err = json.Unmarshal([]byte(data), &questions)
	if err != nil {
		panic(err)
	}

	log.Printf("Room has %d players, %d questions", len(room.Players), len(questions))

	for i, question := range questions {
		log.Printf("Broadcasting question %d", i+1)

		message := map[string]any{
			"type": "question",
			"data": question,
		}

		payload, _ := json.Marshal(message)
		ws.GlobalHub.Broadcast <- ws.BroadcastMessage{
			RoomCode: roomCode,
			Data:     payload,
		}
		room.CurrentQIdx = i + 1
		roomData, _ := json.Marshal(room)
		store.Client.Set(store.Ctx, roomKey, roomData, 60*time.Minute)

		time.Sleep(5 * time.Second)

		scoreboard := make(map[string]int)
		for _, player := range room.Players {
			scoreKey := "score:" + roomCode + ":" + player
			data, err = store.Client.Get(store.Ctx, scoreKey).Result()

			var score int
			score, err = strconv.Atoi(data)
			if err != nil {
				log.Printf("invalid score for player %s: %v", player, err)
				score = 0
			}

			scoreboard[player] = score
		}

		message = map[string]any{
			"type": "scoreboard",
			"data": scoreboard,
		}

		payload, _ = json.Marshal(message)
		ws.GlobalHub.Broadcast <- ws.BroadcastMessage{RoomCode: roomCode, Data: payload}

		time.Sleep(5 * time.Second)
	}

	message := map[string]any{
		"type": "game-over",
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

}
