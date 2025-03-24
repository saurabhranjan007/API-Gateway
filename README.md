# ZenEye API Gateway

## Overview

Zeneye Gateway is an API gateway built with Go and Gin that handles user management, sessions, authentication, authorization, logging, and redirects requests to specific microservices.

## Features

- **User Management**: Create, edit, delete, retrieve users, and list all users.
- **Security**: User authentication, authorization, and role-based access control.
- **Logging**: Centralized logging.
- **Rate Limiting**: Rate limiting for incoming requests.
- **Microservice Routing**: Redirects valid requests to specific microservices.

## Getting Started

### Prerequisites

- Go
- Docker
- PostgreSQL

### Installation

1. Clone the repository:

    ```sh
    https://github.com/saurabhranjan007/API-Gateway
    ```

2. Create `.env` files for different environments with the necessary environment variables:

    - `.env.development`
    - `.env.production`

    Example `.env.development`:

    ```env
    DATABASE_URL=postgres://username:password@localhost:5432/authdb
    DATABASE_NAME=authdb
    JWT_SECRET=your_jwt_secret
    JWT_EXPIRATION=2
    REFRESH_TOKEN_EXPIRATION=720
    RATE_LIMIT=100
    AGENT_SERVICE_URL=http://agent-service:8080
    COMPLIANCE_SERVICE_URL=http://compliance-service:8080
    CONFIGURATION_SERVICE_URL=http://configuration-service:8080
    NOTIFICATION_SERVICE_URL=http://notification-service:8080
    BOT_DETECTION_SERVICE_URL=http://bot-detection-service:8080
    WAF_SERVICE_URL=http://waf-service:8080
    BREACH_DETECTION_SERVICE_URL=http://breach-detection-service:8080
    ADMIN_PANEL_SERVICE_URL=http://admin-panel-service:8080
    ```

3. Run the application:

    #### Running in Different Environments

    ##### Running in Development Environment

    ```sh
    go run main.go
    GO_ENV=development go run main.go
    ```
    By default the environment is set to development. 

    ##### Running in Production Environment

    ```sh
    export GO_ENV=production
    ```

    Alternatively, it can be set directly in the command that runs the application:

    ``` sh 
    GO_ENV=production go run main.go
    ```

    ##### Create Build to Run 
    ```sh
    go build -o zeneye-gateway main.go | go build main.go
    ```

4. Run tests:

    #### Running the Complete Test Cases 

    ```sh
    go test -v ./...
    ```

    #### Running the Specific Test Case

    ##### Running Integration Tests

    To run only integration tests, use the following command:
    
    ```sh 
    go test -v ./pkg/tests/integration/...
    ```

    ##### Running Unit Tests

    To run only unit tests, use the following command:
    
    ```sh 
    go test -v ./pkg/tests/unit/...
    ```

## Usage

### Endpoints

#### User Management
- `POST /users`: Create a new user. (Requires authentication)
- `PATCH /users/:id`: Edit an existing user. (Requires authentication)
- `DELETE /users/:id`: Delete a user. (Requires authentication)
- `GET /users/:id`: Retrieve a user. (Requires authentication)
- `GET /users`: List all users. (Requires authentication)

#### Authentication
- `POST /login`: User login to receive JWT and refresh token.
- `POST /refresh-token`: Refresh access token using the refresh token.

#### Superadmin Management
- `GET /superadmin/check`: Check if a superadmin exists.
- `POST /superadmin/create`: Create a superadmin.

#### Health Check
- `GET /health`: Health check endpoint.

### Microservice Routing

Requests to specific paths will be redirected to the corresponding microservices as defined in the `.env` file.

### Authentication and Authorization

All user management endpoints creation require authentication. A valid JWT must be included in the `Authorization` header of the request. The JWT must be prefixed with `Bearer `.

### User Roles

- `superadmin`: Only one in the entire database. Can perform all actions.
- `admin`
- `department_admin`
- `auditor`

Roles other than the specified ones are not allowed.

### Validation Rules

- **Username**: Minimum of 4 characters, containing only lowercase letters and numbers.
- **Password**: Minimum of 8 characters, containing at least one special character, one number, one uppercase letter, and one lowercase letter.
- **Email**: Must be unique and follow standard email format.
