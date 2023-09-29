package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

//!  This Broker service will ONLY EVER DO ONE THING
//?  Take Requests and forward them off to some microservice and send a response back

const webPort = "80" // docker can listen on port 80

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	fmt.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	// http.ListenAndServe(webPort, CORSMiddleware(app.routes()))
}

func connect() (*amqp.Connection, error) {
	// attempt to connect a fixed number of times
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Rabbitmq Not Yet Ready...")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ!!")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// exponentially increases the time by a power of 2 EX: 1, 4, 9, 16 => 1*1 2*2 3*3 4*4
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
