package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/omishapharswan/devmetrics/backend/internal/models"
)

// ErrScanNotFound is returned by GetScan when no scan with the given ID exists.
var ErrScanNotFound = errors.New("scan not found")

// ListScans returns every persisted scan, most recently scanned first.
func ListScans(conn *sql.DB) ([]models.Scan, error) {
	// scanned_at has only second-level resolution (SQLite's
	// CURRENT_TIMESTAMP), so tie-break on id — which is monotonically
	// increasing — to keep "most recent first" correct even when two
	// scans land in the same second.
	rows, err := conn.Query(
		`SELECT id, folder_path, scanned_at, total_files, avg_complexity, avg_health_score
		 FROM scans ORDER BY scanned_at DESC, id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("query scans: %w", err)
	}
	defer rows.Close()

	scans := []models.Scan{}
	for rows.Next() {
		var s models.Scan
		if err := rows.Scan(&s.ID, &s.FolderPath, &s.ScannedAt, &s.TotalFiles, &s.AvgComplexity, &s.AvgHealthScore); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		scans = append(scans, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scans: %w", err)
	}

	return scans, nil
}

// GetScan returns the full persisted result (scan, files, dependencies)
// for a single scan ID.
func GetScan(conn *sql.DB, scanID int64) (*ScanResult, error) {
	var s models.Scan
	err := conn.QueryRow(
		`SELECT id, folder_path, scanned_at, total_files, avg_complexity, avg_health_score
		 FROM scans WHERE id = ?`,
		scanID,
	).Scan(&s.ID, &s.FolderPath, &s.ScannedAt, &s.TotalFiles, &s.AvgComplexity, &s.AvgHealthScore)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrScanNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query scan: %w", err)
	}

	fileRows, err := conn.Query(
		`SELECT id, scan_id, path, extension, size_bytes, lines_of_code, comment_lines, comment_ratio, complexity, is_duplicate, health_score
		 FROM files WHERE scan_id = ? ORDER BY path`,
		scanID,
	)
	if err != nil {
		return nil, fmt.Errorf("query files: %w", err)
	}
	defer fileRows.Close()

	files := []models.File{}
	for fileRows.Next() {
		var f models.File
		if err := fileRows.Scan(
			&f.ID, &f.ScanID, &f.Path, &f.Extension, &f.SizeBytes,
			&f.LinesOfCode, &f.CommentLines, &f.CommentRatio, &f.Complexity,
			&f.IsDuplicate, &f.HealthScore,
		); err != nil {
			return nil, fmt.Errorf("file row: %w", err)
		}
		files = append(files, f)
	}
	if err := fileRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate files: %w", err)
	}

	depRows, err := conn.Query(
		`SELECT id, scan_id, from_file_id, to_file_id FROM dependencies WHERE scan_id = ?`,
		scanID,
	)
	if err != nil {
		return nil, fmt.Errorf("query dependencies: %w", err)
	}
	defer depRows.Close()

	dependencies := []models.Dependency{}
	for depRows.Next() {
		var d models.Dependency
		if err := depRows.Scan(&d.ID, &d.ScanID, &d.FromFileID, &d.ToFileID); err != nil {
			return nil, fmt.Errorf("dependency row: %w", err)
		}
		dependencies = append(dependencies, d)
	}
	if err := depRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate dependencies: %w", err)
	}

	return &ScanResult{Scan: s, Files: files, Dependencies: dependencies}, nil
}
