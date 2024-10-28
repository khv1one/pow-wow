# Use the official Golang image as the builder
FROM golang:1.22.6-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build the application
RUN go build -o main ./cmd/app/main.go

# Use a minimal image to run the compiled application
FROM alpine:latest

# Set the working directory in the minimal image
WORKDIR /app

# Install necessary packages for the migration script
RUN apk add --no-cache postgresql-client

# Copy the compiled binary from the builder image
COPY --from=builder /app/main .

# Copy the goose binary
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy the migration files
COPY migrations/migrate /app/migrations/migrate

# Copy the entrypoint script
COPY ./build/scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose the port your application listens on
EXPOSE 8080

# Run the entrypoint script
ENTRYPOINT ["/entrypoint.sh"]