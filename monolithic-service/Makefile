MONOLITHIC_BINARY=monolithicApp
STRIPE_SECRET=asfasgsdpfogime_fasd
STRIPE_KEY=ajisgrk134_fasdfw4t57

## MONO APP

## up: starts all containers in the background without forcing build
up: 
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started"


up_monolithic: build_monolithic
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
ifeq ($(wildcard ${MONOLITHIC_BINARY}),)
	@echo "Building new ${MONOLITHIC_BINARY}..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${MONOLITHIC_BINARY} ./cmd/api
else
	@echo "Deleting existing ${MONOLITHIC_BINARY} file..."
	 rm -f ${MONOLITHIC_BINARY}
	@echo "Building new ${MONOLITHIC_BINARY}..."
	go build -o ${MONOLITHIC_BINARY} ./cmd/api
	@echo "building monolithic binary..."
	cd ../monolithic-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MONOLITHIC_BINARY} ./cmd/api
	@echo "Done!"
endif

#build_mail:
#	@echo "building mail binary"
#	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
#	@echo "Done!"
## build_listner: builds the listener binary as a linux executable
#build_listener:
#	@echo "building listener binary"
#	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
#	@echo "Done!"

## build_front: builds the front end binary
#build_front:
#	@echo "Building front end binary..."
#	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
#	@echo "Done!"

## start: starts the front end
#start: build_front
#	@echo "Starting front end"
#	cd  ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop: 
	@echo "Stopping front end..."
	pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"