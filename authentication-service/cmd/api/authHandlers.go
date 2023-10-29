package main

import (
	"fmt"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
)

func (app *Config) isUserAuthorized(token string) (bool, error) {
	_, err := data.ValidateJWTToken(token)
	if err != nil {
		return false, err
	}

	return true, nil
}
func (app *Config) isAuthorizedAndClaims(tokenString string) (data.MyClaims, error) {
	claims, err := data.ValidateAndExtractClaims(tokenString, true)
	if err != nil {
		return data.MyClaims{}, err
	}

	fmt.Println("claims in helper: ", claims)

	myClaims := data.ConvertClaims(claims)

	fmt.Println("myClaims in helper: ", myClaims)

	return myClaims, nil
}
