package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/models"
)

type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	ExpiryMonth   int    `json:"exp_month"`
	ExpiryYear    int    `json:"exp_year"`
	LastFour      string `json:"last_four"`
	Plan          string `json:"plan"`
	ProductID     string `json:"product_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

func (app *application) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HIT AUTH TOKEN")
	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := app.ReadJSON(w, r, &userInput); err != nil {
		app.ErrorJSON(w, err)
		return
	}
	fmt.Println("HIT AUTH TOKEN 2", userInput)

	//* GET USER FROM DB
	user, err := app.DB.User.GetUserByEmail(userInput.Email)
	if err != nil {
		app.invalidCredentials(w)
		return
	}
	fmt.Println("HIT AUTH TOKEN 3", user)

	//* Validate password
	validPassword, err := user.PasswordMatches(userInput.Password)
	if err != nil {
		app.invalidCredentials(w)
		return
	}
	fmt.Println("HIT AUTH TOKEN Valid Password", validPassword)

	if !validPassword {
		app.invalidCredentials(w)
		return
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = false
	payload.Message = "success!"

	app.WriteJSON(w, http.StatusOK, payload)
}

func (app *application) InsertUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("insert user TOKEN")
	var createUser models.User

	if err := app.ReadJSON(w, r, &createUser); err != nil {
		app.ErrorJSON(w, err)
		return
	}
	fmt.Println("insert user TOKEN 2", createUser)

	//* GET USER FROM DB
	userID, err := app.DB.User.InsertUser(createUser)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	payload.Error = false
	payload.Message = "success adding User!"
	payload.Data = userID

	app.WriteJSON(w, http.StatusOK, payload)
}

// because its a handler it needs a ResponseWriter and Request
func (app *application) Ping(w http.ResponseWriter, r *http.Request) {
	// receive payload
	payload := JSONResponse{
		Error:   false,
		Message: "PONG",
	}
	_ = app.WriteJSON(w, http.StatusOK, payload)
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.errorLog.Println(err)
		return
	}
}
