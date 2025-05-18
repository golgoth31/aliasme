# AliasMe

AliasMe is a gRPC service that allows users to create email aliases using the OVH API. It provides a user management system and email alias creation functionality.

## Features

- User management (CRUD operations)
- Email address registration with verification
- Email alias creation using OVH API
- gRPC and HTTP (REST) endpoints
- SQLite database for data persistence

## Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)
- OVH API credentials
- SMTP server for email verification

## Installation

1. Clone the repository:
```bash
git clone https://github.com/golgoth31/aliasme.git
cd aliasme
```

2. Install dependencies:
```bash
make deps
```

3. Generate protobuf code:
```bash
make proto
```

4. Build the project:
```bash
make build
```

## Configuration

The service requires the following environment variables:

```bash
# OVH API Configuration
OVH_ENDPOINT=ovh-eu
OVH_APPLICATION_KEY=your_application_key
OVH_APPLICATION_SECRET=your_application_secret
OVH_CONSUMER_KEY=your_consumer_key

# SMTP Configuration
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password
FROM_EMAIL=noreply@example.com

# Service Configuration
BASE_URL=http://localhost:8080
```

## Running the Service

```bash
make run
```

The service will start with:
- gRPC server on port 9090
- HTTP server on port 8080

## API Documentation

### gRPC Endpoints

The service provides the following gRPC endpoints:

#### User Service
- CreateUser
- GetUser
- UpdateUser
- DeleteUser

#### Email Service
- RegisterEmail
- VerifyEmail
- CreateAlias
- ListAliases

### HTTP Endpoints

The service also provides HTTP endpoints through gRPC-Gateway:

- POST /v1/users - Create a new user
- GET /v1/users/{id} - Get user by ID
- PUT /v1/users/{id} - Update user
- DELETE /v1/users/{id} - Delete user
- POST /v1/emails - Register a new email
- POST /v1/emails/verify - Verify email address
- POST /v1/aliases - Create a new alias
- GET /v1/aliases - List aliases for a user

## Development

### Running Tests

```bash
make test
```

### Building for Linux

```bash
make build-linux
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
