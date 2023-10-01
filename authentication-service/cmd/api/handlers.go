package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
	types "github.com/UpLiftL1f3/Spotify-Micro-Services/shared/types"
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
	var userBody data.User

	err := app.readJSON(w, r, &userBody)
	fmt.Println("userBody:", userBody)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	fmt.Println("userBody:", body)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	userID, err := userBody.Insert()
	if err != nil {
		fmt.Println("Insert Error:", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Generate the Email Token
	token, err := data.GenerateToken(6, userID)

	// Send Email Verification
	mailErr, status := app.emailVerification(w, userBody.Email, token)
	fmt.Println("MAILER ERROR: ", mailErr)
	if mailErr != nil {
		app.errorJSON(w, err, status)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in as user %s", userBody.Email),
		Data:    userID,
	}
	fmt.Println("Insert HIT END:", err)

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

func (app *Config) emailVerification(w http.ResponseWriter, email string, token string) (error, int) {
	fmt.Println("Email 1")
	var mailPayload types.MailMessage
	mailPayload.To = email
	mailPayload.From = "spotifyClone@gmail.com"
	mailPayload.Subject = "Test Mail"
	mailPayload.Message = fmt.Sprintf("Your new generated token is: %s", token)

	jsonData, _ := json.MarshalIndent(mailPayload, "", "\t")
	mailServiceURL := "http://mailer-service/sender"

	fmt.Println("Mail 2")
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err, http.StatusBadRequest
	}

	fmt.Println("Mail 3")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err, http.StatusBadRequest
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		// Handle non-success status code
		return errors.New("Non-success status code: " + response.Status), response.StatusCode
	}

	// create a variable we'll read response.body into
	var jsonFromService JsonResponse

	// decode json from the auth service
	_ = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if jsonFromService.Error {
		return errors.New(jsonFromService.Message), http.StatusBadRequest
	}

	fmt.Println("Mail 4")
	return nil, http.StatusOK
}
