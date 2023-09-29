package main

import (
	"net/http"
)

func (app *Config) routes() http.Handler {
	// mux := chi.NewRouter() //-> New Router using CHI
	mux := http.NewServeMux() //-> New Router using Default library

	// a way to make sure this service is still responding to service requests
	// Add your custom "Heartbeat" middleware for the "/ping" path
	mux.HandleFunc("/ping", app.Ping)

	mux.HandleFunc("/authenticate", app.Authenticate)
	mux.HandleFunc("/addUser", app.InsertNewUser)

	// Create a CORS middleware instance
	corsMiddleware := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "https://*, http://*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Expose-Headers", "Link")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "300")

			// Call the wrapped handler
			handler.ServeHTTP(w, r)
		})
	}

	// Wrap the entire mux with CORS middleware
	handler := corsMiddleware(mux)

	return handler
}
