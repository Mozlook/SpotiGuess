package game

import (
	"backend/internal/model"
	"backend/internal/store"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func tracksFromPlayers(players []string, roomCode string) []model.Track {
	var allTracks []model.Track
	for _, playerID := range players {
		key := fmt.Sprintf("tracks:%s:%s", roomCode, playerID)
		raw, err := store.Client.Get(store.Ctx, key).Result()
		if err != nil {
			log.Println("no tracks for player:", playerID, err)
			continue
		}

		var tracks []model.Track
		err = json.Unmarshal([]byte(raw), &tracks)
		if err != nil {
			log.Println("invalid track data for player:", playerID, err)
			continue
		}

		allTracks = append(allTracks, tracks...)
	}
	return allTracks
}

func tracksFromPlaylist(playlistID string, token string) []model.Track {
	var allTracks []model.Track
	offset := 0
	limit := 100

	client := &http.Client{}

	for {
		reqURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?offset=%d&limit=%d", playlistID, offset, limit)
		req, _ := http.NewRequest("GET", reqURL, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			log.Println("error fetching playlist:", err)
			break
		}
		defer resp.Body.Close()

		var result struct {
			Items []struct {
				Track struct {
					ID         string `json:"id"`
					Name       string `json:"name"`
					DurationMs int    `json:"duration_ms"`
					Artists    []struct {
						Name string `json:"name"`
					} `json:"artists"`
				} `json:"track"`
			} `json:"items"`
			Next string `json:"next"`
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Println("error decoding playlist response:", err)
			break
		}

		for _, item := range result.Items {
			t := item.Track
			var artistNames []string
			for _, a := range t.Artists {
				artistNames = append(artistNames, a.Name)
			}
			allTracks = append(allTracks, model.Track{
				ID:       t.ID,
				Name:     t.Name,
				Artists:  artistNames,
				Duration: t.DurationMs,
			})
		}

		if result.Next == "" {
			break
		}
		offset += limit
	}

	rand.Shuffle(len(allTracks), func(i, j int) {
		allTracks[i], allTracks[j] = allTracks[j], allTracks[i]
	})

	if len(allTracks) > 25 {
		return allTracks[:25]
	}
	return allTracks
}

func tracksFromArtist(artistID string, token string) []model.Track {
	var allTracks []model.Track
	offset := 0
	limit := 50
	client := &http.Client{}

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums?include_groups=album,single&limit=%d&offset=%d", artistID, limit, offset)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error fetching albums:", err)
			break
		}
		defer resp.Body.Close()

		var albumResp struct {
			Items []struct {
				ID string `json:"id"`
			} `json:"items"`
			Next string `json:"next"`
		}
		json.NewDecoder(resp.Body).Decode(&albumResp)

		for _, album := range albumResp.Items {
			albumURL := fmt.Sprintf("https://api.spotify.com/v1/albums/%s/tracks", album.ID)
			albumReq, _ := http.NewRequest("GET", albumURL, nil)
			albumReq.Header.Set("Authorization", "Bearer "+token)
			albumResp, err := client.Do(albumReq)
			if err != nil {
				continue
			}
			defer albumResp.Body.Close()

			var trackResp struct {
				Items []struct {
					ID         string `json:"id"`
					Name       string `json:"name"`
					DurationMs int    `json:"duration_ms"`
					Artists    []struct {
						Name string `json:"name"`
					} `json:"artists"`
				} `json:"items"`
			}
			json.NewDecoder(albumResp.Body).Decode(&trackResp)

			for _, t := range trackResp.Items {
				var artistNames []string
				for _, a := range t.Artists {
					artistNames = append(artistNames, a.Name)
				}
				allTracks = append(allTracks, model.Track{
					ID:       t.ID,
					Name:     t.Name,
					Duration: t.DurationMs,
					Artists:  artistNames,
				})
			}
		}

		if albumResp.Next == "" {
			break
		}
		offset += limit
	}

	rand.Shuffle(len(allTracks), func(i, j int) {
		allTracks[i], allTracks[j] = allTracks[j], allTracks[i]
	})
	if len(allTracks) > 25 {
		return allTracks[:25]
	}
	return allTracks
}
