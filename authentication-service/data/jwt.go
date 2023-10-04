package data

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims represents the custom claims you want to include in the JWT.
type MyClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type JWTPayload struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID string) (string, error) {
	// Create a new set of claims with the user ID and expiration time
	claims := MyClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	// Create the token using the HMAC SHA256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	secretKey := []byte(JWT_Secret) // Replace with your actual secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWTToken(tokenString string) (*JWTPayload, error) {
	secretKey := []byte(JWT_Secret)
	token, err := jwt.ParseWithClaims(tokenString, &JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	fmt.Println("Verify 1")
	if err != nil {
		return nil, err
	}
	fmt.Println("Verify 2")

	if claims, ok := token.Claims.(*JWTPayload); ok && token.Valid {
		// Token is valid, and you can access the claims
		return claims, nil
	} else {
		return nil, fmt.Errorf("token is not valid")
	}
}
