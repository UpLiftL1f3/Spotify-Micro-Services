package functions

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// type Validator interface {
// 	Validate() error
// }

// Generate a NEW email Verification Token
func GenerateToken(length int) (string, error) {
	verificationToken := ""

	for i := 0; i < length; i++ {
		randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNumber := randomGenerator.Intn(length) + 1

		verificationToken += fmt.Sprint(randomNumber)
	}

	return verificationToken, nil
}

// Generate Hex token
func GenerateHexToken(length int) (string, error) {
	tokenBytes := make([]byte, length)
	_, err := cryptoRand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(tokenBytes)
	return token, nil
}

func HashString(str string) (string, error) {
	// Generate a salt for the hash
	salt, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Hash the string (str) with the generated salt
	hashedString, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Concatenate the salt and hashed token (you might want to store them separately in your database)
	finalHash := append(salt, hashedString...)

	// Encode the final hash as a string (if you need to store it as a string)
	return string(finalHash), nil
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{}) // decoder will populate the provided empty struct with the decoded JSON data
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// ReadAndValidateJSON reads JSON from the request and validates it.
// func ReadAndValidateJSON(w http.ResponseWriter, r *http.Request, target interface{}) error {
// 	if err := readJSON(w, r, target); err != nil {
// 		return err
// 	}

// 	if validator, ok := target.(Validator); ok {
// 		if err := validator.Validate(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func ReadAndCustomValidate(w http.ResponseWriter, r *http.Request, payload any, validationFunc func() error) error {
	if err := readJSON(w, r, payload); err != nil {
		return err
	}

	if err := validationFunc(); err != nil {
		return err
	}

	return nil
}
