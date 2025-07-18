package lastfm

import (
	"backend/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type similarResponse struct {
	SimilarTracks struct {
		Track []struct {
			Name string `json:"name"`
		} `json:"track"`
	} `json:"similartracks"`
}

func FetchSimilar(track model.Track) ([]string, error) {

	api_key := os.Getenv("LASTFM_API_KEY")
	endpoint := fmt.Sprintf(
		"http://ws.audioscrobbler.com/2.0/?method=track.getsimilar&track=%s&artist=%s&api_key=%s&format=json&limit=3",
		url.QueryEscape(track.Name),
		url.QueryEscape(track.Artists[0]),
		api_key,
	)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result similarResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var titles []string
	for _, t := range result.SimilarTracks.Track {
		titles = append(titles, t.Name)
	}
	return titles, nil

}
