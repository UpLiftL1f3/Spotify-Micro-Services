package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	// authPb "github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/auths"
	pb "github.com/UpLiftL1f3/Spotify-Micro-Services/broker-service/auths"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//	type AuthServiceServer struct {
//		pb.UnimplementedAuthServiceServer
//	}
const (
	authTarget = "authentication-service:50002"
)

func (app *Config) AuthGRPCRouter(w http.ResponseWriter, r RequestPayload) {
	desiredRoute := strings.Split(r.Action, "/")
	switch desiredRoute[1] {
	case "verifyEmail":
		app.AuthVerifyEmail(w, r.Auth)
	case "signIn":
		fmt.Printf("sign in Received payload: %+v\n", r)
		app.SignIn(w, r.Auth)

	default:
		app.errorJSON(w, errors.New("unknown Action"))
	}
}

func (app *Config) AuthVerifyEmail(w http.ResponseWriter, r AuthPayload) {
	fmt.Println("LOG EVENT VIA GRPC HIT")

	conn, err := grpc.Dial(authTarget, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println("gRPC call failed:", err)
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	fmt.Println("LOG EVENT VIA GRPC HIT pt 2")
	authClient := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println("LOG EVENT VIA GRPC HIT pt 3")
	_, err = authClient.VerifyEmail(ctx, &pb.VerifyEmailRequest{
		UserID: r.UserID,
		Token:  r.Token,
	})
	fmt.Println("LOG EVENT VIA GRPC HIT pt 4")
	if err != nil {
		log.Println("gRPC call failed pt2:", err)
		app.errorJSON(w, err)
		return
	}
	fmt.Println("LOG EVENT VIA GRPC HIT pt 5")

	payload := JsonResponse{
		Error:   false,
		Message: "verified email via GRPC",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) SignIn(w http.ResponseWriter, r AuthPayload) {
	fmt.Println("LOG EVENT VIA GRPC HIT")

	conn, err := grpc.Dial(authTarget, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println("gRPC call failed:", err)
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	fmt.Println("LOG EVENT VIA GRPC HIT pt 2")
	authClient := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println("LOG EVENT VIA GRPC HIT pt 3")
	resp, err := authClient.SignIn(ctx, &pb.SignInRequest{
		Email:    r.Email,
		Password: r.Password,
	})
	fmt.Println("LOG EVENT VIA GRPC HIT pt 4")
	if err != nil {
		log.Println("gRPC call failed pt2:", err)
		app.errorJSON(w, err)
		return
	}
	fmt.Println("LOG EVENT VIA GRPC HIT pt 5")

	payload := JsonResponse{
		Error:   false,
		Message: "user successfully Signed In",
		Data:    resp,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
