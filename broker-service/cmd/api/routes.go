package main

import (
	"net/http"
	"strings"
)

func (app *Config) routes() http.Handler {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// a way to make sure this service is still responding to service requests
	// Add your custom "Heartbeat" middleware for the "/ping" path
	mux.HandleFunc("/ping", app.Ping)

	// Create a CORS middleware instance
	corsMiddleware := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
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

	// Add your routes to the mux
	mux.HandleFunc("/", app.Broker)
	mux.HandleFunc("/log-grpc", app.LogViaGRPC)
	// mux.HandleFunc("/auth-grpc", app.AuthVerifyEmail)
	// mux.Handle("/log-grpc", runFunc(app.LogViaGRPC))
	mux.HandleFunc("/handle", app.HandelSubmission)

	return handler
}

// runFunc is a function that takes another function and applies MethodMiddleware.
func runFunc(handler http.HandlerFunc) http.Handler {
	return MethodMiddleware(http.MethodPost)(handler)
}

// MethodMiddleware is a middleware that restricts the allowed HTTP methods for a route.
func MethodMiddleware(allowedMethods ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, method := range allowedMethods {
		allowed[strings.ToUpper(method)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := allowed[r.Method]; !ok {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
