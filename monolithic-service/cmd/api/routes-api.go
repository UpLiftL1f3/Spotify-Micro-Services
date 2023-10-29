package main

import "net/http"

func (app *application) routes() http.Handler {
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

	// Add your routes to the mux
	mux.HandleFunc("/api/authenticate", app.CreateAuthToken)
	mux.HandleFunc("/api/insertUser", app.InsertUser)
	// mux.HandleFunc("/auth-grpc", app.AuthVerifyEmail)
	// mux.Handle("/log-grpc", runFunc(app.LogViaGRPC))
	// mux.HandleFunc("/handle", app.HandelSubmission)

	// Wrap the entire mux with CORS middleware
	handler := corsMiddleware(mux)

	return handler
}
