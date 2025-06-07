# Go Project Folder Structure

## Main Folder Structure

**`bin/`** - Contains compiled application binaries. Go compiles code into
binaries that land in this directory.

**`cmd/`** - Houses executable entry points and main applications:

- `api/` - Contains the main API server application with HTTP handlers,
  middlewares, and server-related code
- `migrate/` - Contains database migration tools and the `migrations/`
  subfolder for migration files

**`internal/`** - Contains internal packages not meant to be exported
externally:

- `store/` - Houses the data storage layer, repositories, and database
  interactions
- This is where your Postgres database interactions, data validations, email
  sending, rate limiter implementations, and testing mocks live
- The `internal` package in Go is automatically not exported to external
  packages

**`docs/`** - Stores auto-generated documentation, particularly Swagger docs

**`scripts/`** - Contains setup scripts for servers and deployment

**`web/`** (Optional) - If building a full-stack application, this would
contain frontend code (React, Svelte, HTML, etc.)

## Key Architecture Principles

The structure implements **clean layered architecture** with three main layers
cascading like an onion:

1. **Transport Layer** (in `cmd/api/`) - HTTP handlers, routing, server setup
2. **Service Layer** - Business logic (would be in `internal/`)
3. **Storage Layer** (in `internal/store/`) - Database interactions,
   repositories

## Repository Pattern Implementation

The `internal/store/` directory implements the repository pattern:

- `storage.go` - Contains the main storage interface/struct
- `posts.go` - Post-specific storage operations
- `users.go` - User-specific storage operations

This structure allows for:

- **Dependency injection** - layers depend on abstractions, not
  implementations
- **Loose coupling** - easy to swap implementations
- **Testability** - interfaces make mocking straightforward
- **Separation of concerns** - each layer has a clear responsibility

The folder structure prioritizes maintainability and follows Go conventions
while implementing clean architecture principles.
