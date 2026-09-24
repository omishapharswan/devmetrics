package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/omishapharswan/devmetrics/backend/internal/db"
)

const repeatedJSBlock = `console.log("a");
console.log("b");
console.log("c");
console.log("d");
console.log("e");
console.log("f");`

func TestRunScanEndToEnd(t *testing.T) {
	root := t.TempDir()

	write := func(rel, contents string) {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(contents), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}

	write("index.js", "import { helper } from \"./utils\";\n\n"+repeatedJSBlock+"\n")
	write("utils.js", "export function helper() {\n"+repeatedJSBlock+"\n}\n")
	write("README.md", "# not source\n")

	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	result, err := RunScan(root, conn)
	if err != nil {
		t.Fatalf("RunScan returned error: %v", err)
	}

	if result.Scan.ID == 0 {
		t.Error("expected scan to be persisted with a non-zero ID")
	}
	if result.Scan.TotalFiles != 2 {
		t.Errorf("TotalFiles = %d, want 2", result.Scan.TotalFiles)
	}
	if len(result.Files) != 2 {
		t.Fatalf("expected 2 files in result, got %d: %+v", len(result.Files), result.Files)
	}

	byPath := map[string]bool{}
	dependencyEdgeFound := false
	for _, f := range result.Files {
		byPath[f.Path] = f.IsDuplicate
		if f.ScanID != result.Scan.ID {
			t.Errorf("file %s has ScanID %d, want %d", f.Path, f.ScanID, result.Scan.ID)
		}
	}
	if isDup, ok := byPath["index.js"]; !ok || !isDup {
		t.Errorf("expected index.js to be flagged duplicate, got %+v", byPath)
	}
	if isDup, ok := byPath["utils.js"]; !ok || !isDup {
		t.Errorf("expected utils.js to be flagged duplicate, got %+v", byPath)
	}

	var depCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM dependencies WHERE scan_id = ?`, result.Scan.ID).Scan(&depCount); err != nil {
		t.Fatalf("query dependencies: %v", err)
	}
	if depCount != 1 {
		t.Errorf("expected 1 dependency edge persisted, got %d", depCount)
	} else {
		dependencyEdgeFound = true
	}
	if !dependencyEdgeFound {
		t.Error("dependency edge index.js -> utils.js was not persisted")
	}

	var dupBlockCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM duplicate_blocks WHERE scan_id = ?`, result.Scan.ID).Scan(&dupBlockCount); err != nil {
		t.Fatalf("query duplicate_blocks: %v", err)
	}
	if dupBlockCount == 0 {
		t.Error("expected duplicate_blocks rows to be persisted")
	}

	var fileRowCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM files WHERE scan_id = ?`, result.Scan.ID).Scan(&fileRowCount); err != nil {
		t.Fatalf("query files: %v", err)
	}
	if fileRowCount != 2 {
		t.Errorf("expected 2 file rows persisted, got %d", fileRowCount)
	}
}

func TestRunScanFolderNotFound(t *testing.T) {
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	_, err = RunScan(filepath.Join(t.TempDir(), "does-not-exist"), conn)
	if err != ErrFolderNotFound {
		t.Errorf("expected ErrFolderNotFound, got %v", err)
	}
}

func TestRunScanEmptyFolder(t *testing.T) {
	root := t.TempDir()
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	result, err := RunScan(root, conn)
	if err != nil {
		t.Fatalf("RunScan returned error: %v", err)
	}
	if result.Scan.TotalFiles != 0 || len(result.Files) != 0 {
		t.Errorf("expected an empty scan result, got %+v", result)
	}
}
