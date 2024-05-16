# Citizens Data Web Service

This project is a web service written in Go that provides functionalities related to citizens' information. It uses the Go-Chi router and SQLite for storage.

## Features

- Validate citizen's IIN (Individual Identification Number)
- Save citizen's information
- Retrieve citizen's information by IIN
- Retrieve citizen's information by name

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
3. Build the project
```bash
go build -o citizens_data_webservice cmd/citizens-data-webservice/main.go
```

### Usage

Start the server:
```bash
./citizens_data_webservice
```

### Linter

To run golang-ci-lint, run the following command:
```bash
golangci-lint run
```

## Run Using Docker

- docker build --pull --no-cache  -t citizen_webservice .
- docker run -p 8082:8082 -d --name citizen_webservice citizen_webservice

## API Endpoints

- `GET /iin_check/{iin}`: Validate a citizen's IIN
- `POST /people/info`: Save a citizen's information
- `GET /people/info/iin/{iin}`: Retrieve a citizen's information by IIN
- `GET /people/info/name/{name}`: Retrieve a citizen's information by name


## Limitations/ Improvements

1. Security - the current implementation uses BasicAuth for authentication. A more secure method should be used.
2. Test coverage should be improved

## License

This project is licensed under the MIT License - see the `LICENSE.md` file for details.