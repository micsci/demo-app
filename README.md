# Marketplace API

A Go microservice for managing a marketplace of partner applications, built with gRPC and gRPC-Gateway.

## Architecture

```
cmd/main.go                          → entrypoint
internal/
  server.go                          → gRPC server + gRPC-Gateway HTTP server wiring
  config/                            → environment-based configuration
  model/                             → domain models
  handler/                           → gRPC service implementation
  service/                           → business logic + validation
  storage/                           → repository interface
    postgres/                        → PostgreSQL implementation (sqlx + squirrel)
proto/marketplace/v1/                → protobuf service definition with HTTP annotations
gen/go/marketplace/v1/               → generated gRPC + gateway code
migrations/                          → goose SQL migrations
```

**Request flow:** HTTP client → gRPC-Gateway (`:8080`) → gRPC server (`:9090`) → service → storage → PostgreSQL

## Quick Start

```bash
# Start everything with Docker Compose
make up

# Check health
curl http://localhost:8080/health

# View logs
make logs
```

## API (HTTP via gRPC-Gateway)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| POST | `/v1/categories` | Create category |
| GET | `/v1/categories` | List categories |
| DELETE | `/v1/categories/{id}` | Delete category |
| POST | `/v1/applications` | Create application |
| GET | `/v1/applications` | List applications (`?activeOnly=true`) |
| GET | `/v1/applications/{id}` | Get application |
| PUT | `/v1/applications/{id}` | Update application |
| DELETE | `/v1/applications/{id}` | Delete application |
| GET | `/v1/applications/{id}/connections` | List connections for app |
| POST | `/v1/connections` | Create connection |
| PATCH | `/v1/connections/{id}` | Update connection status |
| DELETE | `/v1/connections/{id}` | Delete connection |

gRPC is also available directly on port `9090`.

## Running Tests

```bash
make test
```

## Development

```bash
# Regenerate proto code (requires buf)
make generate

# Run locally (requires Postgres on localhost:5432)
make run

# Docker Compose
make up
make down
```
