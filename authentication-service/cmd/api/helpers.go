package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // omit it whenever its empty
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{}) // decoder will populate the provided empty struct with the decoded JSON data
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data) // converts data into a JSON obj
	if err != nil {
		return err
	}

	// if headers
	if len(headers) > 0 {
		// go through each header at headers[0]
		for key, value := range headers[0] {
			// make the obj with key value pairs
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest // by default if not specified

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()
	fmt.Println("status message is", payload.Message)

	return app.writeJSON(w, statusCode, payload)
}

func customAuthenticatedResp(user *data.User, token []string) interface{} {
	resp := map[string]any{
		"active":     user.Active,
		"id":         user.ID,
		"email":      user.Email,
		"firstName":  user.FirstName,
		"verified":   user.Verified,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"token":      token,
	}

	return resp
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if the Authorization header is present
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the header (assuming it's a Bearer token)
		token := strings.Split(authHeader, "Bearer ")[1]

		// Perform JWT validation here (use your own validation function)
		_, err := data.ValidateJWTToken(token)
		if err != nil {
			app.errorJSON(w, errors.New("Unauthorized"), http.StatusUnauthorized)
			return
		}

		// If the token is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}
