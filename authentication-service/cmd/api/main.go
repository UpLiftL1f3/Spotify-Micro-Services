package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/auth-service/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// ? this is going to be a WEB API SERVICE
// ! SO we need to listen on a specific PORT
// we can do this even tho Broker service is on port 80
// ! We can do this because Docker lets multiple containers listen on the same port and treat them as individual servers
const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	log.Println("Starting authentication service")

	//!  connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant connect to postgres")
	}

	//! Connect Config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

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
	dsn := os.Getenv("DSN")
	log.Println("looking for the DSN: ", dsn)

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
