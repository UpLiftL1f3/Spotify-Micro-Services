package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/logger-service/data"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/logger-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	fmt.Println("WRITE LOG GRPC HIT")
	input := req.GetLogEntry()

	// write log entry
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{Result: "logged!"}
	return res, nil

}

// start the listener

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen to grpc: %v", err)
	}

	srv := grpc.NewServer()

	logs.RegisterLogServiceServer(srv, &LogServer{Models: app.Models})
	log.Printf("grpc server started on port: %s", gRpcPort)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to listen to grpc: %v", err)
	}
}
