# Build stage
FROM golang:1.25-alpine AS builder

ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs (cmd/docs is excluded via .dockerignore)
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.1
RUN swag init -g ./cmd/main.go -o ./cmd/docs

# Build the application
RUN GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary
COPY --from=builder /app/main .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 7080

# Set environment variable for storage path
ENV OBJEXRG_STORAGE_PATH=/app/data

# Run the application
CMD ["./main"]
