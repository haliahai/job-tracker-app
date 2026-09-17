# Job Tracker

A full-stack app for tracking a job search: applications, contacts (recruiters/referrals),
interview stages, prep questions, resume/cover-letter versions, and deadlines — with a
full status history per application.

Stack: Go (stdlib `net/http`, `database/sql`) + MySQL for the API, React for the frontend
(not built yet — backend and schema are done first).

## Status

- [x] Database schema + migrations
- [x] Go API — CRUD for all resources, wired to the schema
- [ ] Go API — build/run verified end-to-end (do this before trusting anything below blindly)
- [ ] React frontend — not started

## Project structure

```
job-tracker-app/
├── backend/
│   ├── cmd/api/main.go          entrypoint: config, DB connection, HTTP server, graceful shutdown
│   ├── internal/config/         env var loading
│   ├── internal/db/             connection pool setup
│   ├── internal/models/         structs mirroring the DB schema
│   ├── internal/store/          DB access layer — parameterized queries, one file per resource group
│   ├── internal/api/            HTTP layer — router, middleware, handlers
│   ├── db/migrations/           golang-migrate SQL migrations (numbered, up/down pairs)
│   ├── go.mod
│   └── .env.example
└── frontend/                    (not created yet)
```

## Prerequisites

- Go (1.22+ — the router uses Go 1.22's method+pattern routing in `net/http`, e.g. `"GET /api/applications/{id}"`)
- MySQL 8.x running locally
- [golang-migrate CLI](https://github.com/golang-migrate/migrate):
  `go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
  (make sure `$(go env GOPATH)/bin` is on your `PATH`, or call it by full path)
- Node.js (for the frontend, once it exists)

## Quick start

### 1. Create the database

```bash
mysql -u root -p -e "CREATE DATABASE job_tracker CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

Create a dedicated app user instead of running the API as `root`:

```bash
mysql -u root -p -e "CREATE USER 'jobtracker'@'localhost' IDENTIFIED BY 'devpass'; \
  GRANT ALL PRIVILEGES ON job_tracker.* TO 'jobtracker'@'localhost'; FLUSH PRIVILEGES;"
```

### 2. Run migrations

```bash
migrate -path backend/db/migrations \
  -database "mysql://jobtracker:devpass@tcp(127.0.0.1:3306)/job_tracker" up
```

Verify:

```bash
mysql -u jobtracker -p job_tracker -e "SHOW TABLES; SELECT * FROM application_statuses;"
```

You should see 8 tables and 10 seeded rows in `application_statuses` (wishlist → applied →
phone_screen → technical_interview → onsite_interview → offer → accepted/rejected/withdrawn/ghosted).

To roll back: `migrate -path backend/db/migrations -database "..." down`

### 3. Configure and run the API

```bash
cd backend
go mod tidy      # resolves github.com/go-sql-driver/mysql, writes go.sum
go build ./...   # compile check before running
```

Set environment variables (see [Configuration](#configuration) below), then:

```bash
DB_USER=jobtracker DB_PASSWORD=devpass go run ./cmd/api
```

Server listens on `:8080` by default. Sanity check:

```bash
curl http://localhost:8080/api/health
curl http://localhost:8080/api/statuses
```

## Configuration

The API reads config from environment variables — there's no `.env` file loader wired in
yet (see `internal/config/config.go`), so either export these in your shell or run with
`env $(cat .env | xargs) go run ./cmd/api` after copying `.env.example` to `.env`.

| Variable        | Default          | Notes                                              |
|-----------------|-------------------|----------------------------------------------------|
| `PORT`          | `8080`            | HTTP port                                          |
| `DATABASE_DSN`  | (built from below)| Full DSN override. If set, `DB_*` vars are ignored. Must include `?parseTime=true` or DATE/DATETIME columns won't scan correctly. |
| `DB_USER`       | `root`            |                                                      |
| `DB_PASSWORD`   | (empty)           |                                                      |
| `DB_HOST`       | `127.0.0.1`       |                                                      |
| `DB_PORT`       | `3306`            |                                                      |
| `DB_NAME`       | `job_tracker`     |                                                      |

CORS is currently hardcoded to allow `http://localhost:5173` (Vite's default dev port) in
`internal/api/middleware.go` — update that (or make it env-driven) before deploying anywhere.

## API reference

All routes are prefixed `/api`. Request/response bodies are JSON. `{id}` is the resource's
numeric ID.

### Health

| Method | Path      | Description       |
|--------|-----------|--------------------|
| GET    | `/health` | Liveness check     |

### Statuses (read-only lookup)

| Method | Path        | Description                                    |
|--------|-------------|--------------------------------------------------|
| GET    | `/statuses` | List all application statuses, ordered by `sort_order` |

### Applications

| Method | Path                        | Description                                                        |
|--------|-----------------------------|----------------------------------------------------------------------|
| GET    | `/applications`             | List applications. Query params: `status_id`, `company` (substring match) |
| POST   | `/applications`              | Create an application                                              |
| GET    | `/applications/{id}`         | Get one application                                                |
| PUT    | `/applications/{id}`         | Update an application (all fields except status)                   |
| PATCH  | `/applications/{id}/status`  | Change status — body: `{"status_id": 3, "notes": "optional"}`. Also writes a row to the status history, in the same transaction. |
| DELETE | `/applications/{id}`         | Delete an application (cascades to all child resources)            |
| GET    | `/applications/{id}/history` | Full status change timeline for an application                     |

`POST`/`PUT` body:

```json
{
  "company_name": "Acme Corp",
  "position_title": "Backend Engineer",
  "job_description": "optional",
  "job_posting_url": "optional",
  "location": "optional",
  "work_mode": "remote | hybrid | onsite",
  "salary_min": 90000,
  "salary_max": 120000,
  "source": "optional, e.g. LinkedIn / referral",
  "current_status_id": 1,
  "priority": "low | medium | high",
  "applied_date": "2026-09-01",
  "notes": "optional"
}
```

### Contacts (recruiters, referrals, interviewers)

| Method | Path                             | Description |
|--------|-----------------------------------|--------------|
| GET    | `/applications/{id}/contacts`     | List contacts for an application |
| POST   | `/applications/{id}/contacts`     | Add a contact |
| PUT    | `/contacts/{id}`                  | Update a contact |
| DELETE | `/contacts/{id}`                  | Delete a contact |

Body: `name`, `role` (`recruiter | referral | hiring_manager | interviewer | other`), `company`,
`email`, `phone`, `linkedin_url`, `notes`.

### Interview stages

| Method | Path                                     | Description |
|--------|--------------------------------------------|--------------|
| GET    | `/applications/{id}/interview-stages`      | List interview stages for an application |
| POST   | `/applications/{id}/interview-stages`      | Add a stage |
| PUT    | `/interview-stages/{id}`                   | Update a stage |
| DELETE | `/interview-stages/{id}`                   | Delete a stage |

Body: `stage_name`, `stage_type` (`phone_screen | technical | behavioral | system_design |
take_home | onsite | panel | final | other`), `scheduled_at`, `completed_at` (RFC3339
timestamps), `format` (`phone | video | onsite | async`), `interviewer_contact_id`,
`outcome` (`pending | passed | failed | cancelled | no_show`, defaults to `pending`),
`feedback_notes`.

### Prep items (questions/topics to prepare, or asked)

| Method | Path                                | Description |
|--------|---------------------------------------|--------------|
| GET    | `/applications/{id}/prep-items`       | List prep items for an application |
| POST   | `/applications/{id}/prep-items`       | Add a prep item |
| PUT    | `/prep-items/{id}`                    | Update a prep item |
| DELETE | `/prep-items/{id}`                    | Delete a prep item |

Body: `interview_stage_id` (optional — omit for general application prep), `category`
(`behavioral | technical | system_design | coding | domain_knowledge | company_specific |
other`, defaults to `other`), `kind` (`to_prepare | asked_to_me | i_asked_them`, defaults to
`to_prepare`), `content` (required), `my_answer_notes`.

### Documents (resume/cover letter versions)

| Method | Path                             | Description |
|--------|-----------------------------------|--------------|
| GET    | `/applications/{id}/documents`    | List documents for an application |
| POST   | `/applications/{id}/documents`    | Add a document record |
| PUT    | `/documents/{id}`                 | Update a document record |
| DELETE | `/documents/{id}`                 | Delete a document record |

Body: `document_type` (`resume | cover_letter | portfolio | other`), `version_label`
(e.g. "Resume v3 - backend focus"), `file_path` (a path or URL — files themselves aren't
stored in MySQL), `notes`.

### Deadlines

| Method | Path                             | Description |
|--------|-----------------------------------|--------------|
| GET    | `/applications/{id}/deadlines`    | List deadlines for an application |
| POST   | `/applications/{id}/deadlines`    | Add a deadline |
| PUT    | `/deadlines/{id}`                 | Update a deadline |
| DELETE | `/deadlines/{id}`                 | Delete a deadline |

Body: `deadline_type` (free text, e.g. `application_deadline`, `assessment_due`,
`offer_decision`), `due_date` (RFC3339, required), `is_completed`, `notes`.

### Errors

Non-2xx responses are `{"error": "message"}`. `404` for missing resources, `400` for bad
input, `500` for anything unexpected (check server logs — the message is deliberately
generic so internals don't leak).

## Database schema

See `backend/db/migrations/` for the source of truth. Eight tables: `application_statuses`
(lookup), `applications`, `application_status_history`, `contacts`, `interview_stages`,
`prep_items`, `documents`, `deadlines`. All child tables cascade-delete when their parent
application is deleted.
