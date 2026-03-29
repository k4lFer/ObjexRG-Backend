# Build stage
FROM golang:1.25-alpine AS builder

ENV GOTOOLCHAIN=local

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.1
RUN swag init -g ./cmd/main.go && mv docs cmd/docs || true

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

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
