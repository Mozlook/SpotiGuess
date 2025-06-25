package main

import (
	"backend/internal/room"
	"backend/internal/spotify"
	"backend/internal/store"
	"fmt"
	"log"
	"net/http"
)

func main() {

	store.InitRedis()

	r := http.NewServeMux()
	r.HandleFunc("/create-room", room.CreateRoomHandler)
	r.HandleFunc("/join-room", room.JoinRoomHandler)
	r.HandleFunc("/room/", room.GetRoomHandler)
	tracks, err := spotify.FetchRecentTracks("BQDuBQGWG4m9s1InPFtO8R50uzbhjJik5a-QFjRxbi6UUko4WFi_iqFv53w03hn-cPLL3hV-YixdP-6yoeR9qgZonIfXvNd4Ui5oCSrBVokoGlwSAToG9ozVM61QGWZFOP0TtVdgIvv1SmhZkF9-OdBKrGfx58fOcIz1I2NbP6-_BzoVWvUVeEKBeHJjmknTucHmRnwAdaQYefNB1sWTFzE7kSWLIfiUt18V")
	fmt.Println(tracks, err)
	log.Println("Server on :8080")
	http.ListenAndServe(":8080", r)
}
