# go-api App

A simple Go application for Backstage demonstration with REST API endpoints.

## Getting Started

The app comes with a Dockerfile and Jenkinsfile for easy deployment and CI/CD integration.

## API Endpoints

- **GET** `/health`: Returns a OK response to check if the service is running.
- **GET** `/metrics`: Returns information about the CPU usage of the system.
- **GET** `/metrics/mem`: Returns information about the memory usage of the system.
- **GET** `/metrics/disk`: Returns information about the storage usage of the system.

## Running Locally

```bash
# Install dependencies
go mod download
go mod tidy

# Run the application
go run .

```

The application will start on port 8080 by default. You can access the API endpoints at `http://localhost:8080`.

### Testing the API

```bash
# Test the health endpoint
curl http://localhost:8080/health

# Test the metrics endpoint
curl http://localhost:8080/metrics

# Test the memory metrics endpoint
curl http://localhost:8080/metrics/mem

# Test the disk metrics endpoint
curl http://localhost:8080/metrics/disk
```

## Running with Docker

### Building the Docker Image

```bash
docker build -t go-api .
```

### Running the Docker Container

```bash
docker run -d -p 8080:8080 go-api
```

### Running with Docker Compose

```yaml
# docker-compose.yml
services:
  go-api:
    image: go-api
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
```

```bash
docker-compose up -d
```
