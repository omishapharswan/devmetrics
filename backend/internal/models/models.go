package models

import "time"

// Scan represents one run of the analyzer over a folder.
type Scan struct {
	ID              int64     `json:"id"`
	FolderPath      string    `json:"folderPath"`
	ScannedAt       time.Time `json:"scannedAt"`
	TotalFiles      int       `json:"totalFiles"`
	AvgComplexity   float64   `json:"avgComplexity"`
	AvgHealthScore  float64   `json:"avgHealthScore"`
}

// File represents the computed metrics for a single scanned source file.
type File struct {
	ID           int64   `json:"id"`
	ScanID       int64   `json:"scanId"`
	Path         string  `json:"path"`
	Extension    string  `json:"extension"`
	SizeBytes    int64   `json:"sizeBytes"`
	LinesOfCode  int     `json:"linesOfCode"`
	CommentLines int     `json:"commentLines"`
	CommentRatio float64 `json:"commentRatio"`
	Complexity   int     `json:"complexity"`
	IsDuplicate  bool    `json:"isDuplicate"`
	HealthScore  float64 `json:"healthScore"`
}

// Dependency represents a directed edge: FromFileID imports/requires ToFileID.
type Dependency struct {
	ID         int64 `json:"id"`
	ScanID     int64 `json:"scanId"`
	FromFileID int64 `json:"fromFileId"`
	ToFileID   int64 `json:"toFileId"`
}

// DuplicateBlock represents one code block whose hash matches another
// block elsewhere in the same scan.
type DuplicateBlock struct {
	ID        int64  `json:"id"`
	ScanID    int64  `json:"scanId"`
	Hash      string `json:"hash"`
	FileID    int64  `json:"fileId"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
}
