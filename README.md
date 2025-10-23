# Citizens Data Web Service

This project is a web service written in Go that provides functionalities related to citizens' information. It uses the Go-Chi router and SQLite for storage.

## Features

- Validate citizen's IIN (Individual Identification Number)
- Save citizen's information with validation
- Retrieve citizen's information by IIN
- Retrieve citizen's information by name
- Health check endpoint for monitoring
- Input validation for phone numbers and names

## Getting Started

### Prerequisites

- Go 1.21 or later
- SQLite

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/citizen_webservice.git
```

2. Navigate to the project directory
```bash
cd citizens_data_webservice
```

3. Install dependencies
```bash
go mod download
```

4. Set up configuration
```bash
export CONFIG_PATH=config/prod.yaml
```

5. Build the project
```bash
go build -o citizens_data_webservice cmd/citizens-data-webservice/main.go
```

### Usage

Start the server:
```bash
export CONFIG_PATH=config/prod.yaml
./citizens_data_webservice
```

### Linter

To run golang-ci-lint:
```bash
golangci-lint run
```

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Run Using Docker

Build the image:
```bash
docker build --pull --no-cache -t citizen_webservice .
```

Run the container:
```bash
docker run -p 8082:8082 -d --name citizen_webservice citizen_webservice
```

The Docker image uses multi-stage builds for optimal size and runs as a non-root user for security.

## API Endpoints

### Health Check
- `GET /health`: Health check endpoint (no authentication required)

### Authenticated Endpoints
All endpoints below require BasicAuth credentials.

- `GET /iin_check/{iin}`: Validate a citizen's IIN
  - Returns: `{correct: bool, sex: string, date_of_birth: string}`

- `POST /people/info`: Save a citizen's information
  - Body: `{iin: string, name: string, phone: string}`
  - Validation:
    - IIN: must be exactly 12 digits and valid
    - Name: 2-100 characters, letters/spaces/hyphens/apostrophes only
    - Phone: Kazakhstan format (+7XXXXXXXXXX, 8XXXXXXXXXX, or 7XXXXXXXXXX)

- `GET /people/info/iin/{iin}`: Retrieve a citizen's information by IIN
  - Returns: `{success: bool, IIN: string, Name: string, Phone: string}`

- `GET /people/info/name/{name}`: Retrieve citizens by name (partial match)
  - Returns: `{success: bool, people: [{IIN, Name, Phone}]}`

- `DELETE /people/delete/{iin}`: Delete a citizen's information by IIN

## Input Validation

The service includes comprehensive input validation:

- **IIN Validation**: Checks length, format, date of birth validity, gender digit, and checksum
- **Phone Validation**: Accepts Kazakhstan phone formats with automatic cleanup of separators
- **Name Validation**: Ensures valid characters, proper length (2-100), and prevents injection attacks


## Configuration

The application uses YAML configuration files. Example (`config/prod.yaml`):

```yaml
env: "prod"
storage_path: "./storage.db"
http_server:
  address: "0.0.0.0:8082"
  timeout: 4s
  idle_timeout: 30s
  user: "user"
  password: "password"
```

Set the config path via environment variable:
```bash
export CONFIG_PATH=config/prod.yaml
```

## Known Limitations

1. **Security**: Currently uses BasicAuth. Consider implementing JWT or OAuth2 for production use.
2. **Test Coverage**: Limited to IIN validator and integration tests. Handler unit tests needed.
3. **Rate Limiting**: No rate limiting on endpoints.
4. **Pagination**: Name search returns all results without pagination.
5. **Database**: SQLite is suitable for development/small deployments. Consider PostgreSQL for production.

## Future Improvements

- Implement JWT authentication
- Add comprehensive unit tests for all handlers
- Add rate limiting middleware
- Implement pagination for search results
- Add database migration tool
- Add metrics and observability (Prometheus, tracing)
- Implement structured error codes

## License

This project is licensed under the MIT License - see the `LICENSE.md` file for details.