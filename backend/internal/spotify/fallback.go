package spotify

import (
	"backend/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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
