package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/helperFunc"
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

type SignInResponse struct {
	ID        string        `json:"id"`
	Active    int           `json:"active"`
	Email     string        `json:"email"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Avatar    models.Avatar `json:"avatar"`
	Verified  bool          `json:"verified"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Token     string        `json:"token,omitempty"`
}

// -> does the user have a valid token
func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	//-> Validate the Token and get associated user
	user, err := app.authenticateToken(r)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	fmt.Println("auth check 10")
	var data = &SignInResponse{
		ID:        user.ID.String(),
		Active:    user.Active,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    helperFunc.UserAvatarDereference(user),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	fmt.Println("auth check 11")
	//-> Valid User - SEND RESPONSE
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    *SignInResponse
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("authenticated user %s", user.Email)
	payload.Data = data

	app.WriteJSON(w, http.StatusOK, payload)
}

func (app *application) authenticateToken(r *http.Request) (*models.User, error) {
	fmt.Println("Auth check 1")
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	fmt.Println("Auth check 2")
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	fmt.Println("Auth check 3", headerParts)
	token := strings.TrimSpace(headerParts[1])
	// if len(token) != 26 {
	// 	return nil, errors.New("no authorization header received")
	// }

	fmt.Println("Auth check 4")
	//! get the user from the Tokens table
	user, err := app.DB.Token.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	fmt.Println("Auth check 5")
	return user, nil
}

// create sign in auth token
// -> SIGN IN FUNCTION
func (app *application) SignIn(w http.ResponseWriter, r *http.Request) {
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

	//-> Make Token suitable for the Frontend

	var data = &SignInResponse{
		ID:        user.ID.String(),
		Active:    user.Active,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    helperFunc.UserAvatarDereference(user),
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token:     token.PlainText,
	}

	//-> SEND RESPONSE
	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		*SignInResponse
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("a token for %s has been created", user.Email)
	payload.SignInResponse = data

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
