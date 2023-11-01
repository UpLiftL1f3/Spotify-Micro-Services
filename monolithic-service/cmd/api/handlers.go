package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

type SingInResponse struct {
	Id        string    `json:"id"`
	Active    int       `json:"active"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     string    `json:"last_four"`
}

// create sign in auth token
// -> SIGN IN FUNCTION
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

	//-> GET USER FROM DB
	user, err := app.DB.User.GetUserByEmail(userInput.Email)
	if err != nil {
		app.invalidCredentials(w)
		return
	}
	fmt.Println("HIT AUTH TOKEN 3", user)

	//-> Validate password
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

	fmt.Println("HIT AUTH TOKEN GENERATE")
	//-> GENERATE Token
	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	fmt.Println("HIT AUTH TOKEN GENERATE (TOKEN)", token)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	//-> SAVE TO DATABASE
	err = app.DB.Token.InsertToken(token, user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	fmt.Println("HIT AUTH TOKEN GENERATED TOKEN: ", token)
	//-> SEND RESPONSE
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		*SingInResponse
	}

	//-> Make Token suitable for the Frontend
	token.Hash = []byte(token.TokenHashToString())
	var data = &SingInResponse{
		Id:        user.ID.String(),
		Active:    user.Active,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     string(token.Hash),
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("a token for %s has been created", user.Email)
	payload.SingInResponse = data

	app.WriteJSON(w, http.StatusOK, payload)
}

// create a new user
func (app *application) InsertUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("insert user TOKEN")
	var createUser models.User

	if err := app.ReadJSON(w, r, &createUser); err != nil {
		app.ErrorJSON(w, err)
		return
	}
	fmt.Println("insert user TOKEN 2", createUser.String())

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
