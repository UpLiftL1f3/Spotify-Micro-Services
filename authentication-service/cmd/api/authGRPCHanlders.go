package main

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/auths"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
	Models data.Models
}

func (a *authServiceServer) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	fmt.Println("verifyEmail 1")

	// -> make payload
	var verifyPayload data.VerifyEmailRequest
	formattedUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return generateEVResponse(false), err

	}
	verifyPayload.UserID = formattedUserID
	verifyPayload.Token = req.Token

	// -> get the Hashed token to compare
	fmt.Println("verifyEmail 2")
	evtBody, err := verifyPayload.FindEmailVerToken()
	if err != nil {
		return generateEVResponse(false), err
	}

	// -> compare the hashed token
	fmt.Println("verifyEmail 3")
	isValid, err := evtBody.CompareHashedToken(verifyPayload.Token)
	if err != nil || !isValid {
		return generateEVResponse(false), err
	}

	// -> declare the fields you're going to update
	fmt.Println("verifyEmail 4")
	updateFields := map[string]interface{}{
		"verified": true,
		// Add more fields as needed
	}

	// -> UPDATE the user
	err = a.Models.User.Update("spotifyClone_schema.users", updateFields, "id", verifyPayload.UserID.String())
	if err != nil {
		fmt.Println("updated error (pt2): ", err)
		return generateEVResponse(false), err
	}

	//! delete the EmailVerificationToken
	err = evtBody.DeleteByID()
	if err != nil {
		fmt.Println("updated error (p3): ", err)
		return generateEVResponse(false), err

	}

	fmt.Println("verifyEmail 5")
	respPayload := generateEVResponse(true)

	return respPayload, nil
}

// * Email Verification (EV)
func generateEVResponse(isVerified bool) *pb.VerifyEmailResponse {
	return &pb.VerifyEmailResponse{
		IsVerified: isVerified,
	}
}

func (a *authServiceServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	fmt.Println("GRPC Authenticate func hit 1")

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(req.Email)
	if err != nil {
		return generateSignInResponse(&data.User{}, []string{}, err)
	}

	fmt.Println("Authenticate func hit 3")
	valid, err := user.PasswordMatches(req.Password)
	fmt.Println("Authenticate func hit 3 valid: ", valid)
	if err != nil || !valid {
		newError := errors.New("invalid credentials")
		return generateSignInResponse(&data.User{}, []string{}, newError)
	}

	fmt.Println("Authenticate func hit 4")
	if len(user.Token) > 0 {
		isAuthorized, err := app.isUserAuthorized(user.Token[0])
		if err != nil || !isAuthorized {
			return generateSignInResponse(&data.User{}, []string{}, err)
		}
		//! USER IS AUTHORIZED
		return generateSignInResponse(user, []string{}, err)
	}

	jwt, err := data.GenerateJWT(user.ID.String())
	if err != nil {
		return generateSignInResponse(&data.User{}, []string{}, err)

	}

	tokenWithJWT := append(user.Token, jwt)

	// Convert []string to pq.StringArray
	pqTokenCopy := pq.StringArray(tokenWithJWT)

	fmt.Println("Authenticate func hit 6 tokenCopy")
	updateFields := map[string]interface{}{
		"token": pqTokenCopy,
	}
	if err := user.Update(data.UsersTableName, updateFields, "id", user.ID); err != nil {
		return generateSignInResponse(&data.User{}, tokenWithJWT, err)
	}

	return generateSignInResponse(user, tokenWithJWT, err)

}

func generateSignInResponse(user *data.User, token []string, err error) (*pb.SignInResponse, error) {
	if user == nil {
		return &pb.SignInResponse{}, err
	}

	if len(user.Token) > 0 {
		token = user.Token
	}

	response := &pb.SignInResponse{
		Active:    int64(user.Active),
		Id:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		Verified:  user.Verified,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Token:     token,
	}

	// Check if LastName is not nil before dereferencing
	if user.LastName != nil {
		response.LastName = *user.LastName
	}

	return response, nil
}
