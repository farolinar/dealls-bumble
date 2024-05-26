# dealls-bumble
A dating app service for dealls.

## Milestones

- [x] Register/Login endpoint
- [ ] AWS S3 bucket integration for image uploads
- [ ] Email verification endpoint
- [ ] Login by email
- [ ] Users matching feature
- [ ] Premium package perk
- [ ] Push docker image to AWS ECR for deployment to AWS ECS
- [ ] Languages supported

## Table of Contents
1. [Project Structure](#project-structure)
2. [Requirement](#requirement)
3. [Getting Started](#getting-started)

## Project Structure

<details>
  <summary>Project structure tree</summary>

  ```bash
├── Dockerfile
├── LICENSE
├── README.md
├── build
│   └── postgres
│       ├── init.sql
│       └── testdata
│           └── init.sql
├── cmd
│   ├── app
│   │   └── server.go
│   ├── main.go
│   └── readiness
│       └── readiness.go
├── config
│   ├── config.go
│   └── postgres
│       └── database.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   └── common
│       ├── auth
│       │   └── token.go
│       ├── db
│       │   ├── db.go
│       │   └── test
│       │       └── db.go
│       ├── jwt
│       │   └── jwt.go
│       ├── middleware
│       │   ├── auth.go
│       │   ├── logging.go
│       │   └── recoverer.go
│       ├── parser
│       │   └── time.go
│       ├── password
│       │   └── bcrypt.go
│       ├── request
│       │   └── request.go
│       ├── response
│       │   └── json.go
│       └── uid
│           └── uid.go
└── services
    ├── base
    │   ├── message.go
    │   ├── request.go
    │   └── response.go
    └── v1
        └── user
            ├── entity.go
            ├── errors.go
            ├── handler.go
            ├── integration_test.go
            ├── message.go
            ├── repository.go
            ├── request.go
            ├── response.go
            ├── service.go
            └── unit_test.go
```
</details>

- `.github/workflows` contains github actions
- `build` contains scripts that will be executed when starting the docker container
- `cmd` is the main folder to execute the service
- `config` contains the configuration for the service
- `internal/common` contains the functionality of the service
    - `auth` contains the functionality for authentication outside the middleware
    - `db` contains the functionality for database-related purposes
        - `test` contains the functionality for database-related integration tests
    - `jwt` contains the functionality for JWT-related functionality
    - `middleware` contains the middleware
    - `parser` contains helper for parsing
    - `password` contains password encryption/decryption functionality
    - `request` contains request parsing functionality
    - `response` contains response parsing functionality
    - `uid` contains the unique identifier generator functionality
- `services` contains available services
- `services/{version}` contains the version 1 of the services
- `services/{version}/{service_name}`
    - `entity.go` is a file for data representations as a struct
    - `errors.go` is a file that contains possible specific errors for the service
    - `handler.go` is a file for http.Handler functions
    - `integration_test.go` is a file for integration tests
    - `messages.go` is a file for response messages
    - `repository.go` is a file for repository
    - `request.go` is a file for request structs and validations
    - `response.go` is a file for response structs
    - `service.go` is a file for service functionality
    - `unit_test.go` is a file for unit tests

## Requirement
In order to run the service, you need
- Docker & Docker Compose v2++ (recommended, but there's an alternative in Getting Started)
- Go (minimum v1.22 recommended)

## Getting Started

Make a new file called `.env` inside the root of the project. Copy the contents of `.env.example` to the `.env` file

### Use Docker
1. Run the docker compose to start the database and the service
```bash
docker-compose up -d
```

### Without Docker Alternative

> [!NOTE]  
> change the value of `POSTGRES_HOST` inside `.env` to `localhost`

1. Start the database
    + Run the sql file script `init.sql` inside `build/postgres` folder to your postgres database
2. Run the service
```bash
go run cmd/main.go
```

