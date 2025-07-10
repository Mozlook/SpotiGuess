package main

import (
	"backend/internal/auth"
	"backend/internal/game"
	"backend/internal/middleware"
	"backend/internal/room"
	"backend/internal/store"
	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

func roomRouterHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) == 3 {
		room.GetRoomHandler(w, r)
	} else if len(parts) == 4 && parts[3] == "questions" {

		game.GetQuestionsHandler(w, r)
	} else if len(parts) == 4 && parts[3] == "scoreboard" {
		game.GetScoreboardHandler(w, r)
	} else {
		http.Error(w, "Invalid room route", http.StatusNotFound)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	store.InitRedis()

	r := http.NewServeMux()
	r.HandleFunc("/create-room", room.CreateRoomHandler)
	r.HandleFunc("/join-room", room.JoinRoomHandler)
	r.HandleFunc("/room/", roomRouterHandler)
	r.HandleFunc("/auth/callback", auth.AuthCallbackHandler)
	r.HandleFunc("/start-game", game.StartGameHandler)
	r.HandleFunc("/submit-answer", game.SubmitAnswerHandler)
	handler := middleware.EnableCORS(r)
	log.Println("Server on :8080")
	http.ListenAndServe(":8080", handler)
}
