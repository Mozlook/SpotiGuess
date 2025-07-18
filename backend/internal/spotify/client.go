package spotify

import (
	"backend/internal/model"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type recentlyPlayedResponse struct {
	Items []struct {
		Track struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"track"`
	} `json:"items"`
}

type recommendationResponse struct {
	Tracks []struct {
		Name string `json:"name"`
		// (opcjonalnie)
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"tracks"`
}

// FetchRecentTracks retrieves the most recently played tracks for a Spotify user,
// using the provided OAuth access token.
//
// It performs a GET request to the Spotify Web API endpoint:
//
//	https://api.spotify.com/v1/me/player/recently-played?limit=25
//
// The function parses the response, extracts relevant fields from each track,
// and returns a slice of Track structs containing:
//   - ID: Spotify track ID
//   - Name: track title
//   - Artists: list of artist names
//
// Only basic metadata is extracted — the preview URL is ignored here.
// The function returns an error if the HTTP request or JSON decoding fails.
//
// Parameters:
//   - token: a valid Spotify OAuth access token
//
// Returns:
//   - []Track: a slice of tracks the user recently listened to
//   - error: if any step fails
func FetchRecentTracks(token string) ([]model.Track, error) {
	url := "https://api.spotify.com/v1/me/player/recently-played?limit=25"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("Spotify FetchRecentTracks status:", resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	log.Println("Spotify response body:", string(body))

	var apiResp recentlyPlayedResponse

	defer resp.Body.Close()
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}
	var tracks []model.Track
	for _, item := range apiResp.Items {
		trackData := item.Track
		var artistNames []string
		for _, artist := range trackData.Artists {
			artistNames = append(artistNames, artist.Name)

		}
		tracks = append(tracks, model.Track{
			ID:      trackData.ID,
			Name:    trackData.Name,
			Artists: artistNames,
		})
	}
	return tracks, nil
}
