package spotify

import (
	"backend/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type recentlyPlayedResponse struct {
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
	body, _ := io.ReadAll(resp.Body)

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
			ID:       trackData.ID,
			Name:     trackData.Name,
			Artists:  artistNames,
			Duration: trackData.DurationMs,
		})
	}
	return tracks, nil
}

func SearchSpotifyHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")
	userID := r.URL.Query().Get("userId")

	if query == "" || (searchType != "playlist" && searchType != "artist") || userID == "" {
		http.Error(w, "Missing or invalid query parameters", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	results, err := SearchSpotify(query, searchType, token)
	if err != nil {
		http.Error(w, "Spotify search failed", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func SearchSpotify(query string, searchType string, token string) ([]map[string]string, error) {
	base := "https://api.spotify.com/v1/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", searchType)
	params.Set("limit", "10")

	req, err := http.NewRequest("GET", base+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var results []map[string]string

	if searchType == "playlist" {
		items := data["playlists"].(map[string]any)["items"].([]any)
		for _, raw := range items {
			p := raw.(map[string]any)

			imageURL := ""
			images := p["images"].([]any)
			if len(images) > 0 {
				image := images[0].(map[string]any)
				imageURL = image["url"].(string)
			}

			results = append(results, map[string]string{
				"id":    p["id"].(string),
				"name":  p["name"].(string),
				"owner": p["owner"].(map[string]any)["display_name"].(string),
				"image": imageURL,
			})
		}
	}

	if searchType == "artist" {
		items := data["artists"].(map[string]any)["items"].([]any)
		for _, raw := range items {
			a := raw.(map[string]any)

			imageURL := ""
			images := a["images"].([]any)
			if len(images) > 0 {
				image := images[0].(map[string]any)
				imageURL = image["url"].(string)
			}

			results = append(results, map[string]string{
				"id":    a["id"].(string),
				"name":  a["name"].(string),
				"image": imageURL,
			})
		}
	}

	return results, nil
}
