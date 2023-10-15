package main

import (
	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
)

func (app *Config) isUserAuthorized(token string) (bool, error) {
	_, err := data.ValidateJWTToken(token)
	if err != nil {
		return false, err
	}

	return true, nil
}
