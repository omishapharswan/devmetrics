# DevMetrics

DevMetrics is a local codebase analytics and visualization platform. Point it
at any folder on disk and it scans the source files, computes code quality
metrics, and renders the results as an interactive dashboard — including a
3D dependency graph of how files reference each other.

Everything runs on `localhost`. There are no external APIs, cloud services,
or AI/ML components involved — all analysis is done with plain algorithmic
logic (parsing, counting, graph building) in Go.

## Features

- Recursive file-tree scanner with source-file detection by extension
- Cyclomatic complexity calculation per function/file
- Lines-of-code and comment-ratio counting
- Hash-based duplicate code block detection (within and across files)
- Dependency graph construction from relative JS/TS and Python imports
- Weighted 0–100 code health score per file
- SQLite-backed scan history for comparing scans over time
- Dashboard with overview cards and complexity/health/extension charts
- Sortable/filterable file table with per-file drill-down
- 3D dependency graph visualization (Three.js / React Three Fiber)
- Scan history view with two-scan comparison (aggregate + per-file deltas)
- Light/dark/system theme toggle

## Tech Stack

| Layer         | Technology                                   |
|---------------|-----------------------------------------------|
| Backend       | Go (Gin)                                       |
| Database      | SQLite (file-based)                            |
| Frontend      | Next.js (App Router) + React + TypeScript      |
| UI components | Tailwind CSS + shadcn/ui                       |
| Charts        | Recharts (via shadcn/ui's chart component)     |
| 3D graph      | Three.js via React Three Fiber (+ drei)        |
| Testing       | Go's testing package (backend); Vitest + React Testing Library (frontend) |

## Project Structure

```
devmetrics/
├── backend/     # Go API server, scanner, analysis engine, SQLite storage
└── frontend/    # Next.js dashboard UI
```

## Prerequisites

- Go 1.27+ (matches the `go` directive in `backend/go.mod`)
- Node.js 20.9+ and npm (Next.js 16's minimum)

## Setup & Running Locally

### 1. Backend (Go API server)

```bash
cd backend
go mod download
go run ./cmd/server
```

The API server starts on `http://localhost:8080` by default. Override with
the `PORT` and `DEVMETRICS_DB_PATH` environment variables (the SQLite file
defaults to `backend/data/devmetrics.db` and its parent directory is
created automatically).

### 2. Frontend (Next.js dashboard)

In a separate terminal:

```bash
cd frontend
npm install
npm run dev
```

The dashboard is available at `http://localhost:3000`. It expects the API
at `http://localhost:8080` by default — copy `frontend/.env.example` to
`frontend/.env.local` and set `NEXT_PUBLIC_API_BASE_URL` to point it
elsewhere. Note the backend's CORS policy currently only allows the
`http://localhost:3000` origin.

### 3. Run a scan

Open `http://localhost:3000` in your browser, enter the absolute path of a
local folder you want to analyze (you can point it at any codebase,
including this repository itself), and start the scan. Results are stored
in a local SQLite database under `backend/data/` so you can compare scans
over time from the History tab.

## API

| Method | Path              | Description                                      |
|--------|-------------------|---------------------------------------------------|
| GET    | `/api/health`     | Server + database status                          |
| POST   | `/api/scan`       | Run a scan of `{ "folderPath": "..." }`, persist it, and return the full result |
| GET    | `/api/scans`      | List all persisted scans, most recent first        |
| GET    | `/api/scans/:id`  | Full result (scan, files, dependencies) for one scan |

## Testing

```bash
# Backend: unit tests for the analysis packages plus HTTP-level
# integration tests (router_test.go boots the real router against a
# real temp SQLite DB).
cd backend && go test ./...

# Frontend: pure-logic unit tests plus component integration tests
# (React Testing Library + user-event, network mocked at the fetch
# boundary).
cd frontend && npm run test
```

## Status

Feature-complete per the original project plan: scanning, all analysis
metrics, the full dashboard (overview, charts, file table, 3D dependency
graph, scan history/comparison), and both backend and frontend test
suites. See commit history for the build-out.
