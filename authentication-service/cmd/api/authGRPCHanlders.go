package main

import (
	"context"
	"fmt"

	pb "github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/auths"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
	"github.com/google/uuid"
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
