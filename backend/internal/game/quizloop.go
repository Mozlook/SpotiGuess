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
