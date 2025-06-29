package main

import (
	"backend/internal/auth"
	"backend/internal/middleware"
	"backend/internal/room"
	"backend/internal/store"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	store.InitRedis()

	r := http.NewServeMux()
	r.HandleFunc("/create-room", room.CreateRoomHandler)
	r.HandleFunc("/join-room", room.JoinRoomHandler)
	r.HandleFunc("/room/", room.GetRoomHandler)
	r.HandleFunc("/auth/callback", auth.AuthCallbackHandler)
	handler := middleware.EnableCORS(r)
	log.Println("Server on :8080")
	http.ListenAndServe(":8080", handler)
}
