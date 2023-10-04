package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/broker-service/event"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/broker-service/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"` //? what do you want to do
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type RPCPayload struct {
	Name string
	Data string
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"Message"`
}

// because its a handler it needs a ResponseWriter and Request
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AT LEAST IT HIT")
	// receive payload
	payload := JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload) //! sending the header OUT
}

func (app *Config) HandelSubmission(w http.ResponseWriter, r *http.Request) {
	// Handle OPTIONS requests for CORS preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var requestPayload RequestPayload

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println("FIRST HIT")

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	case "log":
		app.LogEventViaRPC(w, requestPayload.Log)

	case "mail":
		app.sendMail(w, requestPayload.Mail)

	default:
		app.errorJSON(w, errors.New("unknown Action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some JSON we'll send to the microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	fmt.Printf("FIRST auth HIT %s", jsonData)

	// call the service (equivalent of a http request)
	//! "authentication-service" is what we called the auth micro-service in the docker-compose.yml file
	//? the url set in the authentication-service routes.go = "/authenticate"
	// use bytes.NewBuffer bc the jsonData needs to be passed in, in a specific format
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	fmt.Println("SECOND auth HIT")
	fmt.Println(response.StatusCode)

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusOK {
		fmt.Printf("print the resp %#v", response)
		fmt.Printf("print the resp %#v", response.StatusCode)
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.body into
	var jsonFromService JsonResponse

	// decode json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Register(w http.ResponseWriter, a any) {
	// create some JSON we'll send to the microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	// fmt.Printf("FIRST auth HIT %s", jsonData)

	request, err := http.NewRequest("POST", "http://authentication-service/addUser", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// create a variable we'll read response.body into
	var jsonFromService JsonResponse

	// decode json from the auth service
	_ = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted && response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {

	jsonData, _ := json.MarshalIndent(entry, "", "\t") //! MUST DO NEVER USE "MarshalIndent only use Marshall" (better for Dev env)
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	// conv to json obj
	jsonData, _ := json.MarshalIndent(msg, "", "\t")
	log.Println("sendMail in broker-service")

	// call the mail service
	mailServiceURL := "http://mailer-service/send" // correlates to the docker-compose.yml file Service Name

	// POST TO THE MAIL SERVICE
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	// send the final json response
	var payload JsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, payload LogPayload) {
	fmt.Println("LOGGED VIA RABBIT HIT")
	err := app.pushToQueue(payload.Name, payload.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var resp JsonResponse
	resp.Error = false
	resp.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) LogEventViaRPC(w http.ResponseWriter, l LogPayload) {
	fmt.Println("LOG EVENT VIA RPC HIT")
	client, err := rpc.Dial("tcp", "logger-Service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	fmt.Println("LOG EVENT VIA RPC HIT 2")
	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println("LOG EVENT VIA RPC HIT 3")
	payload := JsonResponse{
		Error:   false,
		Message: result,
	}

	fmt.Println("LOG EVENT VIA RPC HIT 4")
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	// Handle OPTIONS requests for CORS preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("LOG EVENT VIA GRPC HIT")
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println("LOG EVENT VIA GRPC HIT 2")
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	fmt.Println("LOG EVENT VIA GRPC HIT 3")
	client := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println("LOG EVENT VIA GRPC HIT 4")

	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println("WRITE LOG GRPC HIT (POST 2) - ", err)
	payload := JsonResponse{
		Error:   false,
		Message: "logged via GRPC",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// because its a handler it needs a ResponseWriter and Request
func (app *Config) Ping(w http.ResponseWriter, r *http.Request) {
	// receive payload
	payload := JsonResponse{
		Error:   false,
		Message: "PONG",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}
