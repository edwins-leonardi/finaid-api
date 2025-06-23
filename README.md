# Finaid API

A simple REST API built with Go and Gin framework.

## Setup

1. Make sure you have Go installed (version 1.19+)
2. Navigate to the project directory
3. Dependencies are already managed via Go modules

## Running the API

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Available Endpoints

- `GET /hello` - Returns a hello message
- `GET /health` - Health check endpoint

## Testing

You can test the endpoints using curl:

```bash
# Test hello endpoint
curl http://localhost:8080/hello

# Test health endpoint
curl http://localhost:8080/health
```

## Database Setup

The project includes a PostgreSQL database via Docker Compose.

### Start the database:

```bash
docker-compose up -d postgres
```

### Database connection details:
- Host: `localhost`
- Port: `5432`
- Database: `finaid`
- Username: `finaid_user`
- Password: `finaid_password`

### Stop the database:

```bash
docker-compose down
```

## Building

To build the binary:

```bash
go build -o finaid-api main.go
```

Then run:

```bash
./finaid-api
``` 