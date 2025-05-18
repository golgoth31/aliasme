# AliasMe

AliasMe is a gRPC service for managing email aliases using OVH provider. It provides both a gRPC API and CLI commands for managing email aliases.

## Features

- User management
- Email alias creation and management
- OVH integration for alias management
- gRPC API
- CLI interface
- Swagger documentation
- Prometheus metrics

## Prerequisites

- Go 1.21 or later
- OVH API credentials
- SQLite3 (for local development)

## Installation

```bash
go install github.com/golgoth31/aliasme@latest
```

## Configuration

Create a `config.yaml` file in your working directory:

```yaml
logging:
  format: console  # or json
  level: info
  time_format: "2006-01-02 15:04:05"

database:
  path: aliasme.db

ovh:
  endpoint: "https://eu.api.ovh.com/1.0"
  application_key: "your-application-key"
  application_secret: "your-application-secret"
  consumer_key: "your-consumer-key"

smtp:
  host: "smtp.example.com"
  port: "587"
  username: "your-username"
  password: "your-password"
  from_email: "noreply@example.com"
```

## Usage

### Server Commands

Start the server:
```bash
aliasme start
```

Available flags:
- `--config`: Path to config file (default: ./config.yaml)
- `--grpc-port`: gRPC server port (default: 9090)
- `--http-port`: HTTP server port (default: 8080)

### Direct OVH Commands

Create an alias directly in OVH (bypassing the gRPC service):
```bash
aliasme create --domain example.com --source alias@example.com --destination user@example.com
```

Required flags:
- `--domain`: Domain for the alias
- `--source`: Source email address
- `--destination`: Destination email address

### gRPC Client Commands

The client commands allow you to interact with the gRPC service endpoints.

#### Create Alias
```bash
aliasme client create-alias --user-id <user-id> --domain example.com --source alias@example.com --destination user@example.com
```

Required flags:
- `--user-id`: User ID
- `--domain`: Domain for the alias
- `--source`: Source email address
- `--destination`: Destination email address

#### List Aliases
```bash
aliasme client list-aliases --user-id <user-id>
```

Required flags:
- `--user-id`: User ID

## API Documentation

Once the server is running, you can access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## Metrics

Prometheus metrics are available at:
```
http://localhost:8080/metrics
```

## Development

### Project Structure

```
aliasme/
├── cmd/                    # Command-line interface
│   ├── root.go            # Root command
│   ├── start.go           # Start server command
│   ├── create.go          # Direct OVH alias creation
│   └── client.go          # gRPC client commands
├── internal/              # Internal packages
│   ├── config/           # Configuration
│   ├── logger/           # Logging
│   ├── models/           # Data models
│   ├── ovh/              # OVH client
│   ├── server/           # gRPC server
│   ├── service/          # Business logic
│   ├── static/           # Static files
│   └── metrics/          # Prometheus metrics
├── pkg/                  # Public packages
│   └── proto/            # Protocol buffers
├── scripts/              # Build and development scripts
├── static/               # Static assets
├── config.yaml          # Configuration file
├── go.mod               # Go module file
├── go.sum               # Go module checksum
├── LICENSE              # License file
└── README.md            # This file
```

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Generating Protocol Buffers

```bash
make proto
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
