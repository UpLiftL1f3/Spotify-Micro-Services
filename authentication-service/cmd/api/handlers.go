package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
)

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

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authenticate func hit")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	fmt.Println("Authenticate func hit 2")
	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	fmt.Println("Authenticate func hit 2 user: ", user)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusInternalServerError)
		return
	}

	fmt.Println("Authenticate func hit 3")
	valid, err := user.PasswordMatches(requestPayload.Password)
	fmt.Println("Authenticate func hit 3 valid: ", valid)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusInternalServerError)
		return
	}

	fmt.Println("Authenticate func hit 4")
	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) InsertNewUser(w http.ResponseWriter, r *http.Request) {
	var userParams data.User

	err := app.readJSON(w, r, &userParams)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	userID, err := app.Models.User.Insert(userParams)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		// app.errorJSON(w, errors.New("invalid credentials"), http.StatusInternalServerError)
		return
	}

	// log authentication
	// err = app.logRequest("authentication", fmt.Sprintf("%v logged in", userParams.Email))
	// if err != nil {
	// 	app.errorJSON(w, err)
	// 	return
	// }

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", userParams.Email),
		Data:    userID,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) logRequest(name, data string) error {
	fmt.Println("LOG 1")
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	fmt.Println("LOG 2")
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	fmt.Println("LOG 3")
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	fmt.Println("LOG 4")
	return nil
}
