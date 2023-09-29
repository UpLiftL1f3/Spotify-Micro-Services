package main

import (
	"fmt"
	"net/http"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/logger-service/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("write Log")
	// read json into var
	var requestPayload JSONPayload
	fmt.Println("write Log PRE READ")
	_ = app.readJSON(w, r, &requestPayload)

	fmt.Println("write Log POST READ: ", requestPayload)
	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	fmt.Println("write Log POST READ 2", event.Name, event.Data)
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println("write Log POST READ 3")
	resp := JsonResponse{
		Error:   false,
		Message: "logged",
	}

	fmt.Println("write Log POST READ 4", resp)
	app.writeJSON(w, http.StatusAccepted, resp)
}

// because its a handler it needs a ResponseWriter and Request
func (app *Config) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AT LEAST IT HIT")
	// receive payload
	payload := JsonResponse{
		Error:   false,
		Message: "PONG",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}
