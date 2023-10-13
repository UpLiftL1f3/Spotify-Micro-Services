package grpc

// import (
// 	"log"
// 	"net/http"

// 	// authPb "github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/auths"
// 	pb "github.com/UpLiftL1f3/Spotify-Micro-Services/broker-service/auths"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// type AuthServiceServer struct {
// 	pb.UnimplementedAuthServiceServer
// }

// func (s *AuthServiceServer) AuthVerifyEmail(w http.ResponseWriter, r *http.Request) {
// 	// func (s *AuthServiceServer) AuthVerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
// 	// Handle OPTIONS requests for CORS preflight
// 	if r.Method == http.MethodOptions {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	// fmt.Println("LOG EVENT VIA GRPC HIT")
// 	// var requestPayload pb.VerifyEmailRequest
// 	// err := app.readJSON(w, r, &requestPayload)
// 	// if err != nil {
// 	// 	app.errorJSON(w, err)
// 	// 	return
// 	// }

// 	conn, err := grpc.Dial("localhost:50041", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
// 	if err != nil {
// 		log.Fatalf("could not connect to authentication service: %v", err)
// 	}
// 	defer conn.Close()

// 	// authClient := pb.NewAuthServiceClient(conn)
// 	// // authResponse, err := authClient.AuthVerifyEmail(ctx, &pb.VerifyEmailRequest{
// 	// // 	UserID: req.UserID,
// 	// // 	Token:  req.Token,
// 	// // })
// 	if err != nil {
// 		log.Fatalf("authentication service call failed: %v", err)
// 	}

// 	// return &pb.VerifyEmailResponse{IsVerified: authResponse.IsVerified}, nil
// }
