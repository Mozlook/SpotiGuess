package main

import (
	"backend/internal/room"
	"backend/internal/store"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {

	store.InitRedis()

	r := http.NewServeMux()
	r.HandleFunc("/", handler)
	r.HandleFunc("/create-room", room.CreateRoomHandler)
	r.HandleFunc("/join-room", room.JoinRoomHandler)
	log.Println("Server on :8080")
	http.ListenAndServe(":8080", r)
}
