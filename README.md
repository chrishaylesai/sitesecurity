# SiteSecurity

A workforce management platform for the security industry. SiteSecurity handles security workers, their qualifications, shift scheduling, location check-ins, shift reporting, and alarms. Workers may be employed by one or more security companies or operate as freelancers picking up ad-hoc shifts.

## Architecture

The application runs as four services orchestrated with Docker Compose:

```
┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│  Frontend   │   │   Go API    │   │  PostgreSQL  │   │  Keycloak   │
│  Next.js    │──▶│   Chi       │──▶│             │   │  (Auth)     │
│  Port 3000  │   │  Port 8080  │   │  Port 5432  │   │  Port 8180  │
└─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘
```

| Service    | Technology             | Port | Purpose                                    |
|------------|------------------------|------|--------------------------------------------|
| `frontend` | Next.js 16, TypeScript, Tailwind CSS | 3000 | Server-rendered web application    |
| `api`      | Go, Chi router         | 8080 | REST API with layered architecture         |
| `db`       | PostgreSQL 16          | 5432 | Primary data store with SQL migrations     |
| `auth`     | Keycloak 26            | 8180 | OIDC/OAuth2 identity provider              |

Authentication is abstracted behind a Go `auth.Provider` interface so Keycloak can be swapped for another OIDC provider with minimal changes.

## Repository Structure

```
sitesecurity/
├── docker-compose.yml          # All four services
├── api/                        # Go REST API
│   ├── Dockerfile
│   ├── cmd/server/main.go      # Entrypoint — wires repos, services, handlers
│   └── internal/
│       ├── auth/               # Provider interface + Keycloak implementation
│       ├── config/             # Environment-based configuration
│       ├── handler/            # HTTP handlers (one per resource group)
│       ├── middleware/         # Auth, RBAC, logging, CORS
│       ├── model/              # Domain structs
│       ├── repository/         # Database access layer
│       └── service/            # Business logic + validation
├── frontend/                   # Next.js web application
│   ├── Dockerfile
│   └── src/
│       ├── app/                # App Router pages
│       │   ├── companies/      # Company management
│       │   ├── worksites/      # Worksite management
│       │   ├── workers/        # Worker profiles + certificates
│       │   ├── shifts/         # Shift scheduling + assignments
│       │   ├── reports/        # Report templates + submissions
│       │   ├── check-ins/      # GPS location check-ins
│       │   └── alarms/         # Alarm dashboard
│       ├── components/         # Reusable UI components
│       └── lib/                # API client, auth module, types
├── db/
│   ├── init.sh                 # Runs migrations then seeds on first start
│   ├── migrations/             # Numbered SQL scripts (001–011)
│   └── seeds/                  # Realistic test data
└── keycloak/
    └── realm-export.json       # Pre-configured realm for local dev
```

## API Resources

All endpoints are prefixed with `/api/v1/` and require authentication. Write operations are protected by role-based access control.

| Resource          | Endpoint               | Description                              |
|-------------------|------------------------|------------------------------------------|
| Companies         | `/companies`           | Security company CRUD                    |
| Worksites         | `/worksites`           | Worksite CRUD, scoped to company         |
| Workers           | `/workers`             | Profiles, certificates, memberships      |
| Shifts            | `/shifts`              | Scheduling, assignments, status changes  |
| Shift Reports     | `/shift-reports`       | Report templates and submissions         |
| Check-ins         | `/check-ins`           | GPS location recording                   |
| Alarms            | `/alarms`              | Raise, acknowledge, resolve              |

List endpoints support pagination via `?page=1&per_page=25`.

## Running Locally

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose

### Start all services

```bash
docker compose up -d
```

This will:
1. Start PostgreSQL and run all database migrations and seed data
2. Start Keycloak with a pre-configured realm
3. Build and start the Go API
4. Build and start the Next.js frontend

### Verify

```bash
# Check all containers are running
docker compose ps

# Frontend should return 200
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000

# API requires authentication (expect 401)
curl -s http://localhost:8080/api/v1/companies
```

Once running:
- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080
- **Keycloak admin console**: http://localhost:8180 (admin / admin)

### Stop

```bash
docker compose down
```

To also remove the database volume:

```bash
docker compose down -v
```

### Local Development (without Docker)

To run the API and frontend outside Docker while keeping the database and auth in containers:

```bash
# Start only db and auth
docker compose up -d db auth

# Run the API (requires Go 1.24+)
cd api
go run ./cmd/server

# Run the frontend (requires Node 20+)
cd frontend
npm install
npm run dev
```

## Testing

```bash
cd api
go test ./...
```

Tests cover handlers, middleware, and service layers using mocks — no running database required.
