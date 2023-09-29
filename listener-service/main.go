package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go" // advanced messaging que protocol
)

func main() {
	// try to connect to rabbit mq
	//! need a driver (3rd party package to connect)
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	//create consumer (consumes messages from the que)
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the que and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	// attempt to connect a fixed number of times
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		// c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
