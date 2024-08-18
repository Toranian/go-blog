# Use the official Golang image to build the Go application
FROM golang:1.22.5 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main cmd/main.go

# Use a minimal base image for the final stage
FROM debian:bullseye-slim

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/main .

# Expose the port that the app will run on
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
