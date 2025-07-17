package spotify

import (
	"backend/internal/model"
	"encoding/json"
	"fmt"
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

	var apiResp recentlyPlayedResponse

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
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

// FetchRecommendations retrieves a list of recommended track titles from the Spotify API
// based on a given seed track ID.
//
// It performs a GET request to the Spotify endpoint:
//
//	https://api.spotify.com/v1/recommendations
//
// using the provided `trackID` as the seed parameter. The function requests
// 3 recommendations for the Polish market (market=PL).
//
// The response is parsed and the names of the recommended tracks are returned
// as a slice of strings, which can be used as fake quiz answers or distractors.
//
// Parameters:
//   - trackID: Spotify track ID to use as the recommendation seed
//   - token: a valid Spotify OAuth access token
//
// Returns:
//   - []string: list of track titles recommended by Spotify
//   - error: if the request or decoding fails
func FetchRecommendations(trackID string, token string) ([]string, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/recommendations?seed_tracks=%s&limit=3&market=PL", trackID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Spotify FetchRecommendations error:", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Spotify API returned %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("spotify error %d", resp.StatusCode)
	}

	var apiResp recommendationResponse

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}
	var fakeAnswers []string
	for _, track := range apiResp.Tracks {
		fakeAnswers = append(fakeAnswers, track.Name)
	}

	return fakeAnswers, nil
}
