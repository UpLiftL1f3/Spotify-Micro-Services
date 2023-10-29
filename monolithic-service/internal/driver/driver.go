package driver

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var counts int64

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

func ConnectToDB() (*sql.DB, error) {
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
			return connection, nil
		}

		// if we have run the for loop more than 10 times THEN EXIT
		if counts > 10 {
			log.Println(err)
			log.Println("DSN below")
			log.Println(dsn)
			return nil, err
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
