package db

// schema defines every table DevMetrics needs. It is applied with
// CREATE TABLE IF NOT EXISTS so it is safe to run on every startup.
const schema = `
CREATE TABLE IF NOT EXISTS scans (
	id               INTEGER PRIMARY KEY AUTOINCREMENT,
	folder_path      TEXT NOT NULL,
	scanned_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	total_files      INTEGER NOT NULL DEFAULT 0,
	avg_complexity   REAL NOT NULL DEFAULT 0,
	avg_health_score REAL NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS files (
	id                INTEGER PRIMARY KEY AUTOINCREMENT,
	scan_id           INTEGER NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
	path              TEXT NOT NULL,
	extension         TEXT NOT NULL,
	size_bytes        INTEGER NOT NULL DEFAULT 0,
	lines_of_code     INTEGER NOT NULL DEFAULT 0,
	comment_lines     INTEGER NOT NULL DEFAULT 0,
	comment_ratio     REAL NOT NULL DEFAULT 0,
	complexity        INTEGER NOT NULL DEFAULT 0,
	is_duplicate      BOOLEAN NOT NULL DEFAULT 0,
	health_score      REAL NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_files_scan_id ON files(scan_id);

CREATE TABLE IF NOT EXISTS dependencies (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	scan_id        INTEGER NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
	from_file_id   INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
	to_file_id     INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_dependencies_scan_id ON dependencies(scan_id);

CREATE TABLE IF NOT EXISTS duplicate_blocks (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	scan_id        INTEGER NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
	hash           TEXT NOT NULL,
	file_id        INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
	start_line     INTEGER NOT NULL,
	end_line       INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_duplicate_blocks_scan_id ON duplicate_blocks(scan_id);
CREATE INDEX IF NOT EXISTS idx_duplicate_blocks_hash ON duplicate_blocks(hash);
`
