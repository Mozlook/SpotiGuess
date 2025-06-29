package auth

import (
	"backend/internal/store"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type AuthCallbackRequest struct {
	Code string `json:"code"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type SpotifyMeResponse struct {
	ID string `json:"id"`
}

// AuthCallbackHandler handles HTTP POST requests to /auth/callback.
//
// It expects a JSON payload containing an authorization code obtained from Spotify OAuth:
//
//	{
//	  "code": "abc123"
//	}
//
// The handler performs the following steps:
//
//  1. Validates the request body and parses the authorization code.
//
//  2. Sends a POST request to Spotify's token endpoint to exchange the code
//     for an access token and refresh token.
//
//  3. Sends a GET request to Spotify's /me endpoint to retrieve the authenticated user's Spotify ID.
//
//  4. Stores the access token, refresh token, and expiration time in Redis under the key:
//
//     "user:{spotify_id}" → {
//     "access_token": "...",
//     "refresh_token": "...",
//     "expires_at": 1719440000
//     }
//
//  5. Responds with a JSON object containing the access token and Spotify user ID:
//
//     Response:
//     {
//     "access_token": "BQD...",
//     "spotify_id": "user123"
//     }
//
// In case of an error (e.g. invalid code, Spotify API failure, or Redis write failure),
// responds with an appropriate HTTP error status.
func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var body AuthCallbackRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.Code == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := "http://127.0.0.1:5173/callback"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", body.Code)
	data.Set("redirect_uri", redirectURI)

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to exchange token", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	var tokenRes TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenRes)
	if err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	userReq, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)

	userResp, err := client.Do(userReq)
	if err != nil || userResp.StatusCode != 200 {
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	var me SpotifyMeResponse
	err = json.NewDecoder(userResp.Body).Decode(&me)
	if err != nil {
		http.Error(w, "Failed to parse user profile", http.StatusInternalServerError)
		return
	}

	userKey := "user:" + me.ID
	tokenData, _ := json.Marshal(map[string]interface{}{
		"access_token":  tokenRes.AccessToken,
		"refresh_token": tokenRes.RefreshToken,
		"expires_at":    time.Now().Add(time.Duration(tokenRes.ExpiresIn) * time.Second).Unix(),
	})

	err = store.Client.Set(store.Ctx, userKey, tokenData, 60*time.Minute).Err()
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenRes.AccessToken,
		"spotify_id":   me.ID,
	})
}

