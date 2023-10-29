##var for compiled version of the frontend app and broker app
FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp	
LOGGER_BINARY=loggerServiceApp	
AUTH_BINARY=authApp	
MAIL_BINARY=mailerApp
LISTENER_BINARY=listenerApp
MONOLITHIC_BINART=monolithicApp

## MONO APP

## up: starts all containers in the background without forcing build
up: 
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

up_monolithic: build_broker build_auth build_logger
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

# up_build: build_broker build_auth  
## build_broker: builds the broker binary as a linux executable
build_monolithic:
ifeq ($(wildcard ${MONOLITHIC_BINART}),)
	@echo "Building new ${MONOLITHIC_BINART}..."
	cd ../monolithic-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MONOLITHIC_BINART} ./cmd/api
else
	@echo "Deleting existing ${MONOLITHIC_BINART} file..."
	cd ../monolithic-service && rm -f ${MONOLITHIC_BINART}
	@echo "Building new ${MONOLITHIC_BINART}..."
	go build -o ${MONOLITHIC_BINART} ./cmd/api
	@echo "building monolithic binary..."
	cd ../monolithic-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MONOLITHIC_BINART} ./cmd/api
	@echo "Done!"
endif

## build_broker: builds the broker binary as a linux executable
build_broker:
ifeq ($(wildcard ${AUTH_BINARY}),)
	@echo "Building new ${BROKER_BINARY}..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
else
	@echo "Deleting existing ${BROKER_BINARY} file..."
	cd ../broker-service && rm -f ${BROKER_BINARY}
	@echo "Building new ${BROKER_BINARY}..."
	go build -o ${BROKER_BINARY} ./cmd/api
	@echo "building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"
endif
# build_broker:
# 	@echo "building broker binary..."
# 	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
# 	@echo "Done!"


## build_logger: builds the logger binary as a linux executable
build_logger:
ifeq ($(wildcard ${LOGGER_BINARY}),)
	@echo "Building new ${LOGGER_BINARY}..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
else
	@echo "Deleting existing ${LOGGER_BINARY} file..."
	cd ../logger-service && rm -f ${LOGGER_BINARY}
	@echo "Building new ${LOGGER_BINARY}..."
	go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "building logging binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"
endif
# build_logger:
# 	@echo "building logger binary..."
# 	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
# 	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
ifeq ($(wildcard ${AUTH_BINARY}),)
	@echo "Building new ${AUTH_BINARY}..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
else
	@echo "Deleting existing ${AUTH_BINARY} file..."
	cd ../broker-service && rm -f ${AUTH_BINARY}
	@echo "Building new ${AUTH_BINARY}..."
	go build -o ${AUTH_BINARY} ./cmd/api
	@echo "building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"
endif
# build_auth:
# @echo "building auth binary..."
# cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
# @echo "Done!"
## build_mailer: builds the mailer binary as a linux executable
build_mail:
	@echo "building mail binary"
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"
## build_listner: builds the listener binary as a linux executable
build_listener:
	@echo "building listener binary"
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

## build_front: builds the front end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd  ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop: 
	@echo "Stopping front end..."
	pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"