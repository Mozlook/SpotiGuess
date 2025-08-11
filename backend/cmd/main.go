package main

import (
	"backend/internal/auth"
	"backend/internal/game"
	"backend/internal/middleware"
	"backend/internal/room"
	"backend/internal/spotify"
	"backend/internal/store"
	"backend/internal/ws"
	"fmt"
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
	} else if len(parts) == 4 && parts[3] == "next-question" {
		game.GetNextQuestionHandler(w, r)
	} else {
		http.Error(w, "Invalid room route", http.StatusNotFound)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	go ws.GlobalHub.Run()

	store.InitRedis()

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GENERIC handler hit!", r.URL.Path)
		w.Write([]byte("Hello from root: " + r.URL.Path))
	})

	r.HandleFunc("/create-room", room.CreateRoomHandler)
	r.HandleFunc("/join-room", room.JoinRoomHandler)
	r.HandleFunc("/room/", roomRouterHandler)
	r.HandleFunc("/auth/callback", auth.AuthCallbackHandler)
	r.HandleFunc("/start-game", game.StartGameHandler)
	r.HandleFunc("/submit-answer", game.SubmitAnswerHandler)
	r.HandleFunc("/ws/", ws.WSHandler)
	r.HandleFunc("/auth/validate-token", auth.EnsureValidTokenHandler)
	r.HandleFunc("/spotify/search", spotify.SearchSpotifyHandler)
	handler := http.StripPrefix("/go", r)
	handler = middleware.EnableCORS(handler)
	log.Println("Server on :8081")
	http.ListenAndServe(":8081", handler)
}
