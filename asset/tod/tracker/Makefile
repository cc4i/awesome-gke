#

all:	manager

test:	fmt	vet
	REDIS_SERVER_ADDRESS="127.0.0.1:6379" \
	go test ./... -v -cover


# Build manager binary
manager:	fmt	vet
	go build -o tracker main.go

run:	fmt	vet
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

