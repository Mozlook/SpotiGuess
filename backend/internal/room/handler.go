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

type CreateRoomRequest struct {
	HostID string `json:"hostId"`
}

type RoomResponse struct {
	RoomCode string `json:"roomCode"`
}

func generateRoomCode() string {
	code := ""
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for range 6 {
		code += string(characters[rand.Intn(len(characters))])
	}

	return code
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var t CreateRoomRequest
	room := new(Room)
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		panic(err)
	}

	Code := generateRoomCode()

	room.Code = Code
	room.HostId = t.HostID
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
