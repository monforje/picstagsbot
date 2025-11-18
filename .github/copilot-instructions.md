## Goal

Make productive edits to the PicstagsBot Go service with minimal setup. This file focuses on project structure, common patterns, and exact places to change when adding features or fixing bugs.

## Big-picture architecture (what to know first)

- This is a small Go service organized under `internal/` using classic DI: domain interfaces in `internal/domain`, Postgres implementations in `internal/postgres/repoimpl`, services in `internal/service`, and Telegram handlers in `internal/tg`.
- App bootstrap lives in `cmd/picstagsbot/main.go` which calls `app.New()` (`internal/app/app.go`) to wire: config -> postgres -> repoimpl -> services -> handlers -> router -> bot.
- The telegram integration uses `gopkg.in/telebot.v4` with a `Bot` wrapper at `internal/tg/bot/bot.go`. Routing/handlers are under `internal/tg/router` and `internal/tg/handler`.
- DB migrations are stored in `migrations/` and there is a `cmd/migrate` entrypoint. Postgres connection and pooling are managed in `internal/postgres`.
- Logging and small infra helpers are in `pkg/` (e.g., `pkg/logx`, `pkg/middleware`, `pkg/constants`).

## Key files to open when changing behavior

- App wiring: `internal/app/app.go` (shows exact order of initialization, graceful shutdown, rate limiter config)
- Services factory: `internal/service/service.go` (where `Reg`, `Upload`, `Search` are created)
- DB layer: `internal/postgres/*` and `internal/postgres/repoimpl/*` (schema changes + repo impls)
- Telegram bot wrapper: `internal/tg/bot/bot.go` and handlers under `internal/tg/handler/` (search, upload, reg folders)
- Entrypoints: `cmd/picstagsbot/main.go` (run) and `cmd/migrate/main.go` (migrations)

## Project-specific conventions and patterns

- Dependency wiring: constructors return concrete structs and are assembled centrally in `app.New`. Prefer adding a constructor and wiring it there rather than using globals.
- Domain vs implementation: domain interfaces and models are in `internal/domain/*`, implementations live in `internal/postgres` and `internal/postgres/repoimpl`. When adding DB-backed behavior, update the domain interface first, then implement in the repoimpl package.
- Services are thin: `internal/service` exposes logical services (e.g., `RegService`, `UploadService`, `SearchService`). Handlers call service methods; business logic should live in services, not handlers.
- Handlers and router: new Telegram commands or callbacks should be added under `internal/tg/handler/*` and wired in the router (`internal/tg/router/*`). Look at `internal/tg/handler/search/search_handlers.go` for patterns.
- Config: loaded via `config.New(".env")` in `app.New`. Use the `.env` file for dev overrides; production uses env vars.

## Integration points & external deps

- Postgres (pgx) — connection created in `internal/postgres`. Update pooling parameters in `config` and `app.New` where `postgres.Config` is created.
- Migrations: SQL files live in `migrations/` and `cmd/migrate` exists to apply them (goose is used in go.mod). Add new numbered SQL files to `migrations/`.
- Telegram: `gopkg.in/telebot.v4` — the `Bot` wrapper (`internal/tg/bot/bot.go`) is thin; configuration (token, poller timeout) comes from `config.TG`.

## How to build & run (developer commands)

Build or run the service locally from the repo root. These are the minimal commands used by developers.

```pwsh
# Run the bot (reads .env via config.New(".env"))
go run ./cmd/picstagsbot

# Build binary
go build -o bin/picstagsbot ./cmd/picstagsbot

# Run DB migrations (uses cmd/migrate)
go run ./cmd/migrate
```

If you need to use `goose` directly, the project depends on `pressly/goose` (see `go.mod`) but `cmd/migrate` is the preferred entrypoint.

## Quick examples for common tasks

- Add a new API/handler:
  - Add handler file under `internal/tg/handler/<area>` following existing handlers (see `internal/tg/handler/upload/upload_handlers.go`).
  - Add wiring in `internal/tg/router` so the handler is registered with the telegram bot.

- Add DB-backed functionality:
  - Add domain model and repo interface in `internal/domain/*`.
  - Add SQL migration in `migrations/000NN_description.sql`.
  - Implement the interface in `internal/postgres/repoimpl` (use `repoimpl.New(pg.Pool)` as a pattern).
  - Use the service layer (`internal/service`) to expose business logic to handlers.

## Code style & small conventions

- Prefer explicit constructors (NewX) and avoid package-level state.
- Use `pkg/logx` for logging; follow the existing log key/value format (e.g., `logx.Info("msg", "key", val)`).
- Graceful shutdown is expected: respect contexts and check how `App.Stop()` uses timeouts (`cfg.App.ShutdownTimeout`).

## What I looked at to create this guidance

- `cmd/picstagsbot/main.go`
- `internal/app/app.go`
- `internal/service/service.go`
- `internal/tg/bot/bot.go`
- `migrations/` and `cmd/migrate`

If you'd like, I can also:
- Add a short checklist for PR reviewers specific to this repo (migrations, config, tests, router changes).
- Expand examples to include exact function signatures to implement for a new repo method.

Please tell me if you'd prefer a longer version (with more file examples) or a shorter, checklist-only variant.
