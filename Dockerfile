# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./tournament-service ./cmd/app/main.go

# # Stage 2: Create a minimal production image
FROM alpine:latest

# # Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# # Set working directory
WORKDIR /root/

# Copy any static files or templates if needed
COPY --from=builder /app/ .

# # Expose the port your application runs on
EXPOSE 5000

# # Command to run the application
CMD ["./tournament-service"]