# ObjexRG Backend

A simple S3-like storage API built with Go and Gin.

## Features

- Create and delete buckets (folders)
- Upload files with streaming
- Download files with streaming
- List buckets and objects
- Docker support with persistent storage
- Swagger API documentation

## Quick Start

### Using Docker (Recommended)

1. Build and run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

2. The API will be available at `http://localhost:7080`

3. View Swagger docs at `http://localhost:7080/swagger/index.html`

### Local Development

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Install swag for docs:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

3. Generate Swagger docs (from project root):
   ```bash
   swag init -g ./cmd/main.go
   # Then move docs to cmd/docs folder
   mv docs cmd/docs
   ```

4. Run the application:
   ```bash
   go run cmd/main.go
   ```

5. View Swagger docs at: `http://localhost:7080/swagger/index.html`

## API Endpoints

### Buckets

- `POST /bucket/{name}` - Create bucket
- `DELETE /bucket/{name}` - Delete bucket
- `GET /bucket` - List all buckets
- `GET /bucket/{name}` - List objects in bucket
- `GET /bucket/{name}/storage` - List storage in bucket

### Files

- `POST /upload/{bucket}` - Upload file (multipart/form-data, field: `file`)
- `GET /download/{bucket}/{filename}` - Download file

## Swagger Documentation

View interactive API documentation at: `http://localhost:7080/swagger/index.html`

The documentation includes all endpoints with parameters, responses, and examples.

## Storage

Files are stored in `./data` directory (or `/app/data` in Docker).

Set `OBJEXRG_STORAGE_PATH` environment variable to change storage location.

## Streaming

Download uses streaming (`io.Copy`) - does not load entire file into memory.

## Docker

- Multi-stage build for optimized image
- Persistent volume for data storage
- Alpine Linux base for small footprint