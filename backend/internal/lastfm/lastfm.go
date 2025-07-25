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

// FetchSimilar queries the Last.fm API for similar tracks based on the given track.
//
// It uses the Last.fm `track.getsimilar` endpoint to fetch up to 3 related track names,
// which can be used as "fake answers" for a quiz question.
//
// The function performs the following steps:
//
//  1. Constructs a GET request to:
//     http://ws.audioscrobbler.com/2.0/
//     with parameters:
//     - method=track.getsimilar
//     - track: track.Name
//     - artist: track.Artists[0]
//     - api_key: from environment variable LASTFM_API_KEY
//     - format=json
//     - limit=3
//
//  2. Sends the request with Accept: application/json.
//
//  3. Parses the JSON response into a similarResponse struct.
//
//  4. Extracts and returns a slice of similar track names (`[]string`).
//
// Returns:
//   - []string with up to 3 similar track titles
//   - error if the request fails or the response cannot be parsed
//
// Example output:
//
//	["Starboy", "I Feel It Coming", "Can't Feel My Face"]
//
// Note:
// - The function assumes that `track.Artists` is non-empty.
// - If the API key is missing or Last.fm returns an error, the function will fail.
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
