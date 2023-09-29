package main

import (
	"net/http"
)

func (app *Config) routes() http.Handler {
	//! MUX is a common name for referring to "routing" in go
	mux := http.NewServeMux() //? Setting up the new router
	// mux := chi.NewRouter() //? Setting up the new router

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

	// more middleware setup
	mux.HandleFunc("/ping", app.Ping) //! This is a way to check at a later date if this service is still "alive" aka running

	//! ADDING ROUTES
	mux.HandleFunc("/log", app.WriteLog)

	handler := corsMiddleware(mux)

	return handler
}
