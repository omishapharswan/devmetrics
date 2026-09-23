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
- Hash-based duplicate code block detection
- Dependency graph construction from import/require statements
- Weighted 0–100 code health score per file
- SQLite-backed scan history for comparing scans over time
- Dashboard with overview cards, complexity distribution, and file size charts
- Sortable/filterable file table with per-file drill-down
- 3D dependency graph visualization (Three.js / React Three Fiber)
- Scan history comparison view

## Tech Stack

| Layer         | Technology                                   |
|---------------|-----------------------------------------------|
| Backend       | Go (Gin)                                       |
| Database      | SQLite (file-based)                            |
| Frontend      | Next.js + React + TypeScript                   |
| UI components | Tailwind CSS + shadcn/ui                       |
| Dashboard     | Tremor                                         |
| 3D graph      | Three.js via React Three Fiber (+ drei)        |

## Project Structure

```
devmetrics/
├── backend/     # Go API server, scanner, analysis engine, SQLite storage
└── frontend/    # Next.js dashboard UI
```

## Prerequisites

- Go 1.21+
- Node.js 18+ and npm

## Setup & Running Locally

### 1. Backend (Go API server)

```bash
cd backend
go mod download
go run ./cmd/server
```

The API server starts on `http://localhost:8080` by default.

### 2. Frontend (Next.js dashboard)

In a separate terminal:

```bash
cd frontend
npm install
npm run dev
```

The dashboard is available at `http://localhost:3000`.

### 3. Run a scan

Open `http://localhost:3000` in your browser, enter the absolute path of a
local folder you want to analyze (you can point it at any codebase,
including this repository itself), and start the scan. Results are stored
in a local SQLite database under `backend/data/` so you can compare scans
over time.

## Status

This project is under active development. See commit history for progress.
