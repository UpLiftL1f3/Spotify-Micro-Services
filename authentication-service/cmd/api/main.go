package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"net/http"

	pb "github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/auths"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
)

// ? this is going to be a WEB API SERVICE
// ! SO we need to listen on a specific PORT
// we can do this even tho Broker service is on port 80
// ! We can do this because Docker lets multiple containers listen on the same port and treat them as individual servers
const (
	webPort  = "80"
	grpcPort = "50002"
)

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

var app Config

func main() {

	data.LoadEnvVariables()

	//!  connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant connect to postgres")
	}

	//! Connect Config
	app = Config{
		DB:     conn,
		Models: data.New(conn),
	}

	//! GRPC
	go GRPCListen()

	fmt.Printf("Starting auth service on port %s\n", webPort)
	//! set up a web server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

// func main() {
// 	data.LoadEnvVariables()

// 	// Connect to DB asynchronously
// 	conn := connectToDB()
// 	if conn == nil {
// 		log.Panic("Can't connect to postgres")
// 	}
// 	app = Config{
// 		DB:     conn,
// 		Models: data.New(conn),
// 	}

// 	// Wait a bit to allow the database connection to proceed
// 	time.Sleep(2 * time.Second)

// 	// Start gRPC server
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		if err := Listen(); err != nil {
// 			log.Panic(err)
// 		}
// 	}()

// 	fmt.Printf("Starting auth service on port %s\n", webPort)

// 	// set up a web server
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}

// 	// Run the HTTP server in a goroutine
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		if err := srv.ListenAndServe(); err != nil {
// 			log.Panic(err)
// 		}
// 	}()

// 	// Wait for both servers to finish starting
// 	wg.Wait()
// }

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	// Set the desired schema
	schema := "schema=spotifyClone_schema"

	// Construct the DSN with the schema
	dsn := fmt.Sprintf("host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable search_path=%s timezone=UTC", schema)

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		// if we have run the for loop more than 10 times THEN EXIT
		if counts > 10 {
			log.Println(err)
			log.Println("DSN below")
			log.Println(dsn)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func GRPCListen() {
	fmt.Println("STARTED GRPC LISTEN")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("STARTED GRPC LISTEN 2")
	grpcServer := grpc.NewServer()
	fmt.Println("STARTED GRPC LISTEN 3")
	pb.RegisterAuthServiceServer(grpcServer, &authServiceServer{Models: app.Models})
	fmt.Println("STARTED GRPC LISTEN 4")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
