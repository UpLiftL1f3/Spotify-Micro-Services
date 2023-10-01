package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/shared/functions"
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
	var userBody data.CreateUserRequest
	if err := functions.ReadAndCustomValidate(w, r, &userBody, userBody.Validate); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	userID, err := userBody.Insert(data.UsersTableName)
	if err != nil {
		fmt.Println("Insert Error:", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Generate the Email Token
	token, err := data.GenerateTokenAndCreateEVDocument(userID)
	if err != nil {
		fmt.Println("Insert Error:", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Send Email Verification
	mailErr, status := app.sendEmailWithString(w, userBody.Email, token)
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

func (app *Config) sendEmailWithString(w http.ResponseWriter, email string, str string) (error, int) {
	fmt.Println("Email 1")
	var mailPayload types.MailMessage
	mailPayload.To = email
	mailPayload.From = "spotifyClone@gmail.com"
	mailPayload.Subject = "Test Mail"
	mailPayload.Message = fmt.Sprintf("Your new generated token is: %s", str)

	jsonData, _ := json.MarshalIndent(mailPayload, "", "\t")
	mailServiceURL := "http://mailer-service/send"

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

func (app *Config) verifyEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("verifyEmail 1")

	// -> make payload
	var verifyPayload data.VerifyEmailRequest
	if err := functions.ReadAndCustomValidate(w, r, &verifyPayload, verifyPayload.Validate); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// -> get the Hashed token to compare
	fmt.Println("verifyEmail 2")
	evtBody, err := verifyPayload.FindEmailVerToken()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// -> compare the hashed token
	fmt.Println("verifyEmail 3")
	isValid, err := evtBody.CompareHashedToken(verifyPayload.Token)
	if err != nil || !isValid {
		app.errorJSON(w, fmt.Errorf("issue with token"))
		return
	}

	// -> declare the fields you're going to update
	fmt.Println("verifyEmail 4")
	updateFields := map[string]interface{}{
		"verified": true,
		// Add more fields as needed
	}

	// -> UPDATE the user
	err = app.Models.User.Update("spotifyClone_schema.users", updateFields, "id", verifyPayload.UserID.String())
	if err != nil {
		fmt.Println("updated error (pt2): ", err)

		app.errorJSON(w, err)
		return
	}

	//! delete the EmailVerificationToken
	err = evtBody.DeleteByID()
	if err != nil {
		fmt.Println("updated error (p3): ", err)
		app.errorJSON(w, err)
		return
	}

	fmt.Println("verifyEmail 5")
	payload := JsonResponse{
		Error:   false,
		Message: "Your email is verified",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendReverificationEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reverification Email 1")

	// -> make payload
	var verifyPayload data.VerifyEmailRequest
	if err := functions.ReadAndCustomValidate(w, r, &verifyPayload, verifyPayload.ValidateID); err != nil {
		app.errorJSON(w, err)
		return
	}

	//* Get User
	var user *data.User
	user, err := user.GetByID(verifyPayload.UserID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//! delete the EmailVerificationToken
	err = verifyPayload.DeleteByUserID()
	if err != nil {
		fmt.Println("updated error (p3): ", err)
		app.errorJSON(w, err)
		return
	}

	// -> Create new Token and EmailVerificationDocument
	// Generate the Email Token
	token, err := data.GenerateTokenAndCreateEVDocument(verifyPayload.UserID)
	if err != nil {
		fmt.Println("Insert Error:", err)
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Send Email Verification
	mailErr, status := app.sendEmailWithString(w, user.Email, token)
	fmt.Println("MAILER ERROR: ", mailErr)
	if mailErr != nil {
		app.errorJSON(w, err, status)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("verification email sent to %s", user.Email),
		Data:    verifyPayload.UserID,
	}
	fmt.Println("reverification HIT END:", err)

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) generateResetPasswordViaEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reverification Email 1")

	// -> make payload
	var resetPassPayload data.ResetPasswordRequest
	if err := functions.ReadAndCustomValidate(w, r, &resetPassPayload, resetPassPayload.ValidateEmail); err != nil {
		app.errorJSON(w, err)
		return
	}

	//* Get User
	var user *data.User
	user, err := user.GetByEmail(resetPassPayload.Email)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if err := user.Validate(); err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println("USER ID: ", user.ID)

	//* FIND ONE IF EXISTS AND DELETE
	if err := data.FindAndDeleteByID(user.ID); err != nil {
		app.errorJSON(w, err)
		return
	}

	token, err := data.GenerateTokenAndCreateRPDocument(user.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resetPassLink := user.GenerateResetLink(token)

	// Send Email Verification
	mailErr, status := app.sendEmailWithString(w, user.Email, resetPassLink)
	fmt.Println("MAILER ERROR: ", mailErr)
	if mailErr != nil {
		app.errorJSON(w, err, status)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("reset password link sent to %s", user.Email),
	}
	fmt.Println("Reset Password HIT END:", err)

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) validateResetPassToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("validate reset pass 1")

	// -> make payload
	var resetPassPayload data.ResetPasswordRequest
	if err := functions.ReadAndCustomValidate(w, r, &resetPassPayload, resetPassPayload.ValidateWithoutPassword); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// -> get the Hashed token to compare
	fmt.Println("verifyEmail 2")
	resetPassDoc, err := resetPassPayload.FindResetPassToken()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// -> compare the hashed token
	fmt.Println("verifyEmail 3")
	isValid, err := resetPassDoc.CompareHashedToken(resetPassPayload.Token)
	if err != nil || !isValid {
		app.errorJSON(w, fmt.Errorf("issue with token"))
		return
	}

	//! delete the ResetPassTokenDocument
	err = resetPassDoc.DeleteByID()
	if err != nil {
		fmt.Println("updated error (p3): ", err)
		app.errorJSON(w, err)
		return
	}

	fmt.Println("verifyEmail 5")
	payload := JsonResponse{
		Error:   false,
		Message: "You can reset your password",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	// -> make payload
	var resetPassPayload data.ResetPasswordRequest
	if err := functions.ReadAndCustomValidate(w, r, &resetPassPayload, resetPassPayload.ValidateWithoutToken); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	// -> Find USER
	fmt.Println("verifyEmail 4")

	user, err := app.Models.User.GetByID(resetPassPayload.UserID)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}

	// -> RESET PASSWORD
	user.ResetPassword(resetPassPayload.Password)

	fmt.Println("verifyEmail 5")
	payload := JsonResponse{
		Error:   false,
		Message: "password reset!",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
