package auth

import (
	"backend/internal/store"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Code == "" {
		log.Println("Invalid request body:", err)
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

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("Failed to build token request:", err)
		http.Error(w, "Token request error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send token request:", err)
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Token exchange failed (status %d): %s", resp.StatusCode, string(bodyBytes))
		http.Error(w, "Failed to exchange token", http.StatusBadRequest)
		return
	}

	var tokenRes struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"` // seconds
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		log.Println("Failed to decode token response:", err)
		http.Error(w, "Token decode error", http.StatusInternalServerError)
		return
	}

	// 🔍 Get Spotify profile
	userReq, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)

	userResp, err := client.Do(userReq)
	if err != nil {
		log.Println("Failed to get user profile:", err)
		http.Error(w, "User profile fetch failed", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(userResp.Body)
		log.Printf("Failed to get user profile (status %d): %s", userResp.StatusCode, string(bodyBytes))
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	var me struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&me); err != nil {
		log.Println("Failed to decode user profile:", err)
		http.Error(w, "User decode error", http.StatusInternalServerError)
		return
	}

	// 💾 Save token data to Redis
	tokenData, _ := json.Marshal(map[string]any{
		"access_token":  tokenRes.AccessToken,
		"refresh_token": tokenRes.RefreshToken,
		"expires_at":    time.Now().Add(time.Duration(tokenRes.ExpiresIn) * time.Second).Unix(),
	})
	userKey := "user:" + me.ID
	if err := store.Client.Set(store.Ctx, userKey, tokenData, 60*time.Minute).Err(); err != nil {
		log.Println(" Failed to save token data:", err)
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	log.Printf("Auth completed for user: %s", me.ID)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenRes.AccessToken,
		"spotify_id":   me.ID,
	})
}
