package spotify

import (
	"backend/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SimiliarFallback attempts to generate fallback track recommendations
// by using Spotify's Search API instead of relying on Last.fm.
//
// It takes the name and artists of the original track and performs a
// search against Spotify's `/search` endpoint.
//
// It filters out unwanted variants (e.g. remixes, live versions) and ensures
// that the original track does not appear in the results.
//
// Parameters:
//   - track: model.Track object containing ID, Name, Artists (used for search query)
//   - token: A valid Spotify access token with `user-read-private` scope
//
// Returns:
//   - []string: A slice of 3 recommended (but filtered) track names
//   - error: If the API request fails or the response cannot be parsed
//
// Example response:
//
//	["Good Vibes", "Let It Go", "Feel the Beat"]
//
// Filtering logic:
//   - Reject tracks whose names contain keywords like "remix", "acoustic", "live", etc.
//   - Reject tracks that match the original `track.Name` (case-insensitive)
//   - Only return the first 3 valid alternative tracks
//
// Spotify endpoint used:
//
//	GET https://api.spotify.com/v1/search?q=<track+name+artist>&type=track&limit=10&market=PL
func SimiliarFallback(track model.Track, token string) ([]string, error) {
	query := fmt.Sprintf("%s %s", track.Name, strings.Join(track.Artists, " "))
	searchURL := "https://api.spotify.com/v1/search?" + url.Values{
		"q":      {query},
		"type":   {"track"},
		"limit":  {"10"},
		"market": {"PL"},
	}.Encode()

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Tracks struct {
			Items []struct {
				Name string `json:"name"`
			} `json:"items"`
		} `json:"tracks"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var tracks []string
	bannedKeywords := []string{"acoustic", "remix", "live", "instrumental", "karaoke", "cover"}
	for _, item := range result.Tracks.Items {
		skip := false
		for _, keyword := range bannedKeywords {
			if strings.Contains(strings.ToLower(item.Name), keyword) {
				skip = true
				break
			}
		}
		if strings.EqualFold(track.Name, item.Name) || skip {
			continue
		}

		tracks = append(tracks, item.Name)
		if len(tracks) == 3 {
			break
		}
	}

	return tracks, nil
}
