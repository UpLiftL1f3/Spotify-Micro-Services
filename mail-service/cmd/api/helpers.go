package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // One Megabyte MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes)) // it is intended for limiting the size of incoming request bodies

	dec := json.NewDecoder(r.Body) // creates a new decoder type which can later be used to decode (or read) JSON data

	err := dec.Decode(data)
	if err != nil {
		return fmt.Errorf("Error decoding using readJSON: %w", err)
	}

	err = dec.Decode(&struct{}{}) // this is verifying if the dec is a single json obj
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
		for key, value := range headers[0] {
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

	// if the status parameter has a status code then set that as the new statusCode
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
