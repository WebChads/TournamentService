package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check header format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Prepare request to authorization service
		authReq := struct {
			Token string `json:"token"`
		}{Token: tokenString}

		reqBody, err := json.Marshal(authReq)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create HTTP client with timeout
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		// Make POST request to authorization service
		resp, err := client.Post(
			"http://auth-service.com/validate-token",
			"application/json",
			bytes.NewBuffer(reqBody),
		)
		if err != nil {
			http.Error(w, "Auth Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		type AuthResponse struct {
			Valid   bool   `json:"valid"`
			UserID  string `json:"user_id"`
		}
		var authResp AuthResponse
		if err := json.Unmarshal(body, &authResp); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !authResp.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Define own type to avoid collisions
		type userId string
		const id userId = "user_id"

		// Add user information to context
		ctx := context.WithValue(r.Context(), id, authResp.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
