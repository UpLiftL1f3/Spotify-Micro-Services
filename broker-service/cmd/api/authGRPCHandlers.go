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

// type AuthServiceServer struct {
// 	pb.UnimplementedAuthServiceServer
// }

func (app *Config) AuthGRPCRouter(w http.ResponseWriter, r RequestPayload) {
	desiredRoute := strings.Split(r.Action, "/")
	switch desiredRoute[1] {
	case "verifyEmail":
		fmt.Printf("Received payload: %+v\n", r)
		app.AuthVerifyEmail(w, r.Auth)

	default:
		app.errorJSON(w, errors.New("unknown Action"))
	}
}

func (app *Config) AuthVerifyEmail(w http.ResponseWriter, r AuthPayload) {
	// func (s *AuthServiceServer) AuthVerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	// Handle OPTIONS requests for CORS preflight
	// if r.Method == http.MethodOptions {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }

	fmt.Println("LOG EVENT VIA GRPC HIT")
	// var requestPayload pb.VerifyEmailRequest
	// err := app.readJSON(w, r, &requestPayload)
	// if err != nil {
	// 	app.errorJSON(w, err)
	// 	return
	// }

	fmt.Println("LOG EVENT VIA GRPC HIT FUCKKKK")
	conn, err := grpc.Dial("authentication-service:50002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
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
