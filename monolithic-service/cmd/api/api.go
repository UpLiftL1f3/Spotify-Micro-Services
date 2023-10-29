package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/driver"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/env"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/models"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	// stripe struct {
	// 	secret string
	// 	key    string
	// }
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
	Users    models.UserModel
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf(fmt.Sprintf("starting back end server in %s mode on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe()
}

func main() {
	env.LoadEnvVariables()
	fmt.Println("POSTGRES DSN:", env.DSNEnv)

	var cfg config
	flag.IntVar(&cfg.port, "port", 80, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.db.dsn, "dsn", env.DSNEnv, "DSN")

	flag.Parse()

	// cfg.stripe.key = os.Getenv("STRIPE_KEY")
	// cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.ConnectToDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	modelsInstance := models.NewModels(conn)

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		DB:       modelsInstance.DB,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
