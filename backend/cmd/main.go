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

	err := store.Client.Set(store.Ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := store.Client.Get(store.Ctx, "foo").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("foo", val)

	r := http.NewServeMux()
	r.HandleFunc("/", handler)
	r.HandleFunc("/create-room", room.CreateRoomHandler)
	log.Println("Server on :8080")
	http.ListenAndServe(":8080", r)
}
