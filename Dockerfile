FROM golang:1.21-alpine

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup folder
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY go.mod ./

RUN go mod download

COPY . .

# Build the Go app
ENV CGO_ENABLED=1
ENV CONFIG_PATH=config/prod.yaml
RUN go build -a /app/cmd/citizens-data-webservice

# Expose port 8080 to the outside world
EXPOSE 8082

# Run the executable, path to main.go file
ENTRYPOINT exec go run cmd/citizens-data-webservice/main.go