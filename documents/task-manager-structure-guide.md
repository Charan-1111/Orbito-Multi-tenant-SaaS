# Task Manager Project Structure Guide

## Purpose

This document defines the recommended structure for the `task-manager` Go service and records the main organization decisions for future development.

## Current State

The project currently contains these main areas:

```text
task-manager/
├── cmd/api/              # Application entrypoint
├── configs/local/       # Local configuration
├── internal/config/     # Configuration loading
├── internal/logging/    # Application logging
├── internal/middlewares/ # HTTP middleware
├── internal/server/     # Application and route setup
├── internal/store/      # Database and repository code
├── models/              # Shared data models
├── utils/               # Utility code
├── go.mod
└── go.sum
```

This is a valid starting point, but the implementation and the older structure diagram are not fully synchronized yet.

## Recommended Structure

```text
task-manager/
├── cmd/
│   └── api/
│       └── main.go              # Process entrypoint and OS signals
├── internal/
│   ├── config/                  # Configuration types and loading
│   ├── database/                # Database connection and migrations
│   ├── logging/                 # Logger setup and cleanup
│   ├── middleware/              # HTTP middleware
│   ├── models/                  # Application-only domain models
│   ├── server/                  # HTTP server, routes, and lifecycle
│   ├── task/                    # Task handlers, services, and repositories
│   └── tenant/                  # Tenant handlers, services, and repositories
├── configs/
│   └── local/
│       └── config.json          # Local-only configuration
├── migrations/                  # Database schema migrations
├── tests/                       # Integration and end-to-end tests
├── go.mod
└── go.sum
```

## Package Responsibilities

### `cmd/api`

Contains only process startup concerns:

- Create the application.
- Register `SIGINT` and `SIGTERM` handlers.
- Start the application.
- Wait for graceful shutdown.
- Exit with an appropriate status.

### `internal/config`

Owns configuration models and configuration loading. Configuration models should include every section used by the application, including server and database settings.

### `internal/server`

Owns HTTP routes, middleware registration, and the HTTP server lifecycle. It should receive already-created dependencies instead of constructing every dependency itself as the project grows.

### `internal/database`

Owns database connection setup, health checks, migrations, and database cleanup. Database resources must be closed during application shutdown.

### `internal/task` and `internal/tenant`

Each feature package should keep its related handler, service, and repository code together. This makes feature ownership clear and limits cross-package coupling.

### `internal/models`

Contains models used only by this service. Keep models outside `internal` only when another project must import them directly.

### `internal/middleware`

Contains request ID, authentication, logging, recovery, and similar HTTP middleware. Use the singular package name `middleware` unless there is a strong reason to keep `middlewares`.

## Important Cleanup Items

1. Keep `documents/project-structure.txt` synchronized with the real repository, or replace it with this guide as the canonical structure document.
2. Avoid maintaining duplicate configuration loaders in both `internal/config` and `utils`.
3. Move application-only models from `models` to `internal/models` when external reuse is not required.
4. Keep dependency construction in a bootstrap or application-construction layer, not inside route handlers.
5. Run `go mod tidy` after imports stabilize so direct and indirect dependencies are accurate.
6. Add unit tests for configuration loading, services, and middleware.
7. Add integration tests for database access and HTTP routes.

## Graceful Shutdown

The application should use this lifecycle:

```text
OS signal
   |
   v
cmd/api/main.go
   |
   v
shutdown channel
   |
   v
internal/server application lifecycle
   |
   +--> stop accepting new HTTP requests
   +--> wait for active requests to finish
   +--> close database connections
   +--> close logger resources
   +--> exit
```

Shutdown operations should be safe to call only once. A timeout should be used so the process cannot wait forever for unfinished requests.

## Structure Decision

The current layout is suitable for an early prototype. Before adding substantial task and tenant functionality, adopt the recommended feature-oriented structure and keep this document updated whenever a package boundary changes.
