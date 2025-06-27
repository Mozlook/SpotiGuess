package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Track struct {
	ID      string
	Name    string
	Artists []string
}

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

func FetchRecentTracks(token string) ([]Track, error) {
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
	var tracks []Track
	for _, item := range apiResp.Items {
		trackData := item.Track
		var artistNames []string
		for _, artist := range trackData.Artists {
			artistNames = append(artistNames, artist.Name)

		}
		tracks = append(tracks, Track{
			ID:      trackData.ID,
			Name:    trackData.Name,
			Artists: artistNames,
		})
	}
	return tracks, nil
}

func FetchRecommendations(trackID string, token string) ([]string, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/recommendations?seed_tracks=%s&limit=3&market=PL", trackID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
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
