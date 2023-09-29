package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	log.Println("hit sendMail handler")

	var requestPayload mailMessage
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	log.Println("About to send back response from Mail handlers")

	app.writeJSON(w, http.StatusAccepted, payload)

}

// because its a handler it needs a ResponseWriter and Request
func (app *Config) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AT LEAST IT HIT")
	// receive payload
	payload := JsonResponse{
		Error:   false,
		Message: "PONG",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}
