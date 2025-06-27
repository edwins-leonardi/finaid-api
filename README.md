# Finaid API

A simple REST API built with Go and Gin framework.

## Setup

1. Make sure you have Go installed (version 1.19+)
2. Navigate to the project directory
3. Dependencies are already managed via Go modules

## Environment Configuration

Create a `.env` file in the root directory with the following content:

```env
# App Configuration
APP_NAME=finaid-api
APP_ENV=development

# HTTP Configuration
HTTP_URL=http://localhost:8080
HTTP_PORT=8080
HTTP_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Database Configuration
DB_CONNECTION=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=finaid_user
DB_PASSWORD=finaid_password
DB_NAME=finaid
```

**Important**: The `HTTP_ALLOWED_ORIGINS` should include your frontend URL (default: `http://localhost:5173`).

## Running the API

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Available Endpoints

- `GET /hello` - Returns a hello message
- `GET /health` - Health check endpoint
- `GET /api/v1/persons` - List persons (with pagination)
- `POST /api/v1/persons` - Create a new person

## Testing

You can test the endpoints using curl:

```bash
# Test hello endpoint
curl http://localhost:8080/hello

# Test health endpoint
curl http://localhost:8080/health

# Test persons endpoint
curl http://localhost:8080/api/v1/persons
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