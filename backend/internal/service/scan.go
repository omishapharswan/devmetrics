// Package service orchestrates the analysis packages and the scanner
// package into a single scan operation, persisting the result to SQLite.
package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/omishapharswan/devmetrics/backend/internal/analysis"
	"github.com/omishapharswan/devmetrics/backend/internal/models"
	"github.com/omishapharswan/devmetrics/backend/internal/scanner"
)

// ErrFolderNotFound is returned by RunScan when rootPath does not exist or
// is not a directory.
var ErrFolderNotFound = errors.New("folder not found")

// ScanResult is the response returned to callers of RunScan: the scan
// summary row plus the per-file metrics computed for it.
type ScanResult struct {
	Scan  models.Scan   `json:"scan"`
	Files []models.File `json:"files"`
}

// fileMetrics holds every value computed for one file before it is
// persisted, so duplicate detection and health scoring (which both need
// the full set of files) can run before any database writes happen.
type fileMetrics struct {
	found      scanner.FoundFile
	content    string
	lines      analysis.LineCounts
	complexity int
}

// RunScan walks rootPath, computes line/complexity/duplication/dependency
// metrics for every recognized source file, persists a new scan and its
// files/dependencies/duplicate blocks to conn, and returns the persisted
// result.
func RunScan(rootPath string, conn *sql.DB) (*ScanResult, error) {
	info, err := os.Stat(rootPath)
	if os.IsNotExist(err) {
		return nil, ErrFolderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", rootPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", rootPath)
	}

	found, err := scanner.Walk(rootPath)
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", rootPath, err)
	}

	metrics := make([]fileMetrics, 0, len(found))
	for _, f := range found {
		content, err := os.ReadFile(f.AbsPath)
		if err != nil {
			// Unreadable file (permissions, race with deletion, etc.):
			// skip it rather than failing the whole scan.
			continue
		}
		source := string(content)
		metrics = append(metrics, fileMetrics{
			found:      f,
			content:    source,
			lines:      analysis.CountLines(source, f.Language),
			complexity: analysis.CalculateComplexity(source),
		})
	}

	dupSources := make([]analysis.SourceFile, len(metrics))
	for i, m := range metrics {
		dupSources[i] = analysis.SourceFile{FileID: int64(i), Content: m.content}
	}
	dupBlocks := analysis.DetectDuplicates(dupSources)

	isDuplicate := make([]bool, len(metrics))
	for _, b := range dupBlocks {
		isDuplicate[b.FileID] = true
	}

	healthScores := make([]float64, len(metrics))
	for i, m := range metrics {
		healthScores[i] = analysis.CalculateHealthScore(analysis.HealthScoreInput{
			Complexity:   m.complexity,
			LinesOfCode:  m.lines.LinesOfCode,
			CommentRatio: m.lines.CommentRatio,
			IsDuplicate:  isDuplicate[i],
		})
	}

	depSources := make([]analysis.DependencySource, len(metrics))
	for i, m := range metrics {
		depSources[i] = analysis.DependencySource{
			RelPath:  m.found.RelPath,
			Language: m.found.Language,
			Content:  m.content,
		}
	}
	depEdges := analysis.BuildDependencyGraph(depSources)

	return persistScan(conn, rootPath, metrics, isDuplicate, healthScores, dupBlocks, depEdges)
}

func persistScan(
	conn *sql.DB,
	rootPath string,
	metrics []fileMetrics,
	isDuplicate []bool,
	healthScores []float64,
	dupBlocks []analysis.DuplicateBlock,
	depEdges []analysis.DependencyEdge,
) (*ScanResult, error) {
	tx, err := conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO scans (folder_path) VALUES (?)`, rootPath)
	if err != nil {
		return nil, fmt.Errorf("insert scan: %w", err)
	}
	scanID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("read scan id: %w", err)
	}

	files := make([]models.File, len(metrics))
	fileIDs := make([]int64, len(metrics))
	relPathToID := make(map[string]int64, len(metrics))

	for i, m := range metrics {
		res, err := tx.Exec(
			`INSERT INTO files (scan_id, path, extension, size_bytes, lines_of_code, comment_lines, comment_ratio, complexity, is_duplicate, health_score)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			scanID, m.found.RelPath, m.found.Extension, m.found.SizeBytes,
			m.lines.LinesOfCode, m.lines.CommentLines, m.lines.CommentRatio,
			m.complexity, isDuplicate[i], healthScores[i],
		)
		if err != nil {
			return nil, fmt.Errorf("insert file %s: %w", m.found.RelPath, err)
		}
		fileID, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("read file id for %s: %w", m.found.RelPath, err)
		}

		fileIDs[i] = fileID
		relPathToID[m.found.RelPath] = fileID
		files[i] = models.File{
			ID:           fileID,
			ScanID:       scanID,
			Path:         m.found.RelPath,
			Extension:    m.found.Extension,
			SizeBytes:    m.found.SizeBytes,
			LinesOfCode:  m.lines.LinesOfCode,
			CommentLines: m.lines.CommentLines,
			CommentRatio: m.lines.CommentRatio,
			Complexity:   m.complexity,
			IsDuplicate:  isDuplicate[i],
			HealthScore:  healthScores[i],
		}
	}

	for _, b := range dupBlocks {
		if _, err := tx.Exec(
			`INSERT INTO duplicate_blocks (scan_id, hash, file_id, start_line, end_line) VALUES (?, ?, ?, ?, ?)`,
			scanID, b.Hash, fileIDs[b.FileID], b.StartLine, b.EndLine,
		); err != nil {
			return nil, fmt.Errorf("insert duplicate block: %w", err)
		}
	}

	for _, e := range depEdges {
		fromID, fromOK := relPathToID[e.FromRelPath]
		toID, toOK := relPathToID[e.ToRelPath]
		if !fromOK || !toOK {
			continue
		}
		if _, err := tx.Exec(
			`INSERT INTO dependencies (scan_id, from_file_id, to_file_id) VALUES (?, ?, ?)`,
			scanID, fromID, toID,
		); err != nil {
			return nil, fmt.Errorf("insert dependency edge: %w", err)
		}
	}

	var avgComplexity, avgHealthScore float64
	if len(metrics) > 0 {
		var totalComplexity float64
		var totalHealth float64
		for i, m := range metrics {
			totalComplexity += float64(m.complexity)
			totalHealth += healthScores[i]
		}
		avgComplexity = totalComplexity / float64(len(metrics))
		avgHealthScore = totalHealth / float64(len(metrics))
	}

	if _, err := tx.Exec(
		`UPDATE scans SET total_files = ?, avg_complexity = ?, avg_health_score = ? WHERE id = ?`,
		len(metrics), avgComplexity, avgHealthScore, scanID,
	); err != nil {
		return nil, fmt.Errorf("update scan totals: %w", err)
	}

	var scannedAt time.Time
	if err := tx.QueryRow(`SELECT scanned_at FROM scans WHERE id = ?`, scanID).Scan(&scannedAt); err != nil {
		return nil, fmt.Errorf("read scan timestamp: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit scan: %w", err)
	}

	return &ScanResult{
		Scan: models.Scan{
			ID:             scanID,
			FolderPath:     rootPath,
			ScannedAt:      scannedAt,
			TotalFiles:     len(metrics),
			AvgComplexity:  avgComplexity,
			AvgHealthScore: avgHealthScore,
		},
		Files: files,
	}, nil
}
