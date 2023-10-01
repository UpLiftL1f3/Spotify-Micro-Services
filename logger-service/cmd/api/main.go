package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
	rpcPort = "5001"
	// mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	fmt.Println("logger api main func ran")

	data.LoadEnvVariables()

	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Listen to RPC connection
	rpcErr := rpc.Register(new(RPCServer))
	if rpcErr != nil {
		log.Panic(err)
	}
	go app.rpcListen()

	// Listen to GRPC connection
	fmt.Println("logger GRPC listen hit")
	go app.gRPCListen()
	fmt.Println("logger GRPC listen hit 2")

	// start web server
	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC Server on port: ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		log.Println("Listening for RPC requests to accept")

		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	fmt.Println("connect to mongo hit")
	clientOptions := options.Client().ApplyURI(data.ConnectionString)

	// Create a new MongoDB client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		fmt.Println("MongoDB Client Error:")
		log.Fatal(err)
		return nil, err
	}

	// Connect to MongoDB
	err = client.Connect(context.Background())
	if err != nil {
		fmt.Println("MongoDB connect Client Error")
		log.Fatal(err)
		return nil, err
	}

	// Ping the MongoDB server
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("MongoDB connect Client ping:", err)
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}

// func connectToMongo() (*mongo.Client, error) {
// 	// Use the SetServerAPIOptions() method to set the Stable API version to 1
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
// 	opts := options.Client().ApplyURI(data.ConnectionString).SetServerAPIOptions(serverAPI)
// 	// opts.SetAuth(options.Credential{
// 	// 	Username: "admin",
// 	// 	Password: "password",
// 	// })

// 	client, err := mongo.Connect(context.TODO(), opts)
// 	if err != nil {
// 		log.Println("Error connection: ", err)
// 		return nil, err
// 	}

// 	return client, nil
// }
