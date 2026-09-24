package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/omishapharswan/devmetrics/backend/internal/db"
)

func TestListAndGetScan(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	first, err := RunScan(root, conn)
	if err != nil {
		t.Fatalf("first RunScan: %v", err)
	}
	second, err := RunScan(root, conn)
	if err != nil {
		t.Fatalf("second RunScan: %v", err)
	}

	scans, err := ListScans(conn)
	if err != nil {
		t.Fatalf("ListScans: %v", err)
	}
	if len(scans) != 2 {
		t.Fatalf("expected 2 scans, got %d", len(scans))
	}
	// Most recent first.
	if scans[0].ID != second.Scan.ID || scans[1].ID != first.Scan.ID {
		t.Errorf("ListScans order = %+v, want most recent (%d) first", scans, second.Scan.ID)
	}

	got, err := GetScan(conn, first.Scan.ID)
	if err != nil {
		t.Fatalf("GetScan: %v", err)
	}
	if got.Scan.ID != first.Scan.ID {
		t.Errorf("GetScan returned scan %d, want %d", got.Scan.ID, first.Scan.ID)
	}
	if len(got.Files) != 1 || got.Files[0].Path != "a.go" {
		t.Errorf("GetScan files = %+v, want a single a.go entry", got.Files)
	}
}

func TestGetScanNotFound(t *testing.T) {
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	_, err = GetScan(conn, 999)
	if err != ErrScanNotFound {
		t.Errorf("expected ErrScanNotFound, got %v", err)
	}
}
