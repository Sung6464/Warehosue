# Use the official Go image as a builder
FROM golang:1.24.2 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./main.go

# Use a minimal base image for the final stage
FROM alpine:latest

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Set working directory for the final stage
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# !!! IMPORTANT: Make the executable file executable !!!
RUN chmod +x /app/main

# Expose the port your service listens on (Customer Service usually 8085)
EXPOSE 8085

# Command to run the executable
CMD ["./main"]