package main

import (
	"backend/internal/room"
	"backend/internal/store"
	"log"
	"net/http"
)

func main() {

	store.InitRedis()

	r := http.NewServeMux()
	r.HandleFunc("/create-room", room.CreateRoomHandler)
	r.HandleFunc("/join-room", room.JoinRoomHandler)
	r.HandleFunc("/room/", room.GetRoomHandler)
	log.Println("Server on :8080")
	http.ListenAndServe(":8080", r)
}
