# Movie App Backend

A RESTful API backend for a movie application built with Go, following clean architecture principles.

## Project Architecture

```
movie_app_backend/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── db/
│   │   └── db.go               # Database connection setup
│   ├── migrations/             # SQL migration files
│   ├── model/
│   │   └── models.go           # Domain models and DTOs
│   ├── repo/
│   │   └── movie_repo/         # Repository layer (data access)
│   │       ├── interface.go    # Repository interface
│   │       ├── base_repo.go    # Repository struct
│   │       ├── getMovieBasebyId.go
│   │       ├── searchMovie.go
│   │       ├── discoverMovie.go
│   │       ├── fetchGenres.go
│   │       ├── fetchCompanies.go
│   │       └── fetchCredits.go
│   ├── service/
│   │   ├── movie_service.go      # Business logic layer
│   │   └── movie_service_test.go # Unit tests
│   ├── transport/
│   │   └── http/
│   │       ├── routes.go           # Route definitions
│   │       ├── movie_handler.go    # HTTP handlers
│   │       └── movie_handler_test.go
│   └── seed.go                 # Database seeding script
├── go.mod
├── go.sum
└── README.md
```

### Layer Overview

| Layer | Responsibility |
|-------|----------------|
| **Transport (HTTP)** | Handle HTTP requests/responses, input validation, routing |
| **Service** | Business logic, orchestration, data transformation |
| **Repository** | Database operations, SQL queries, data persistence |
| **Model** | Domain entities and data transfer objects |

### Data Flow

```
HTTP Request → Handler → Service → Repository → Database
                ↓           ↓           ↓
           Validation   Business    SQL Query
                        Logic
```

## Tech Stack

- **Language**: Go 1.25+
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **UUID**: [gofrs/uuid](https://github.com/gofrs/uuid) (UUIDv7)

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 14+(optional if you use docker)
- golang migrate

## Local Setup

### 1. Clone the Repository

```bash
git clone https://github.com/h-raju-arch/movie_app_backend.git
cd movie_app_backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up PostgreSQL

Get Postgredb url 
put the url in the current shell using export DATABASE_URL="<your db url>"
To check if the url exported correctly or not you can use echo DATABASE_URL

### 5. Run Migrations

Apply the database migrations from `internal/migrations/`:

```bash
# Using a migration tool like golang-migrate
migrate -path internal/migrations -database "$DATABASE_URL" up

# Or manually apply SQL files in order
psql $DATABASE_URL -f internal/migrations/20260123155917_add_movies_table.up.sql
# ... apply other migration files
```

### 6. Seed the Database (Optional)

Populate the database with sample data:

```bash
go run internal/seed.go
```

### 7. Run the Application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:3000`.

## API Endpoints

### Get Movie by ID

```http
GET /api/movie/?id={uuid}&lang={language}&append_to_response={fields}
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | UUID | Yes | Movie ID |
| `lang` | string | No | Language code (default: `en`) |
| `append_to_response` | string | No | Comma-separated: `genres`, `companies`, `credits` |

**Example:**
```bash
curl "http://localhost:3000/api/movie/?id=550e8400-e29b-41d4-a716-446655440000&append_to_response=genres,credits"
```

### Search Movies

```http
GET /api/movies/search?query={query}&language={lang}&page={page}&page_size={size}
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `query` | string | Yes | Search term |
| `language` | string | No | Language (default: `en-US`) |
| `include_adult` | boolean | No | Include adult content (default: `false`) |
| `year` | int | No | Filter by release year |
| `region` | string | No | Filter by region/country |
| `page` | int | No | Page number (default: `1`) |
| `page_size` | int | No | Results per page (default: `20`, max: `100`) |

**Example:**
```bash
curl "http://localhost:3000/api/movies/search?query=inception&page=1&page_size=10"
```

### Discover Movies

```http
GET /api/movies/discover?language={lang}&sort_by={sort}&with_genres={genres}
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `language` | string | No | Language (default: `en`) |
| `sort_by` | string | No | Sort order (default: `popularity.DESC`) |
| `with_genres` | string | No | Genre IDs: comma (AND) or pipe (OR) separated |
| `include_adult` | boolean | No | Include adult content (default: `false`) |
| `releaseGTE` | string | No | Release date >= (YYYY-MM-DD) |
| `releaseLTE` | string | No | Release date <= (YYYY-MM-DD) |
| `VoteAvgGTE` | float | No | Vote average >= |
| `VoteAvgLTE` | float | No | Vote average <= |
| `page` | int | No | Page number (default: `1`) |
| `page_size` | int | No | Results per page (default: `20`, max: `100`) |


