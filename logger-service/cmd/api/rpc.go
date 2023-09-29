package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/logger-service/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	fmt.Println("LOG INFO RPC HIT")

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo: ", err)
		return err
	}

	fmt.Println("LOG INFO RPC HIT 2")
	*resp = "Process payload payload via RPC: " + payload.Name
	return nil
}
