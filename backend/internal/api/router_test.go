package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/omishapharswan/devmetrics/backend/internal/api"
	"github.com/omishapharswan/devmetrics/backend/internal/db"
)

// newTestServer boots the real router (full middleware stack, real
// handlers) against a fresh temp SQLite database, so these tests exercise
// the actual HTTP request/response cycle rather than calling service
// functions directly.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	srv := httptest.NewServer(api.NewRouter(conn))
	t.Cleanup(srv.Close)
	return srv
}

func TestHealthEndpoint(t *testing.T) {
	srv := newTestServer(t)

	res, err := http.Get(srv.URL + "/api/health")
	if err != nil {
		t.Fatalf("GET /api/health: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" || body["db"] != "ok" {
		t.Errorf("body = %+v, want status=ok db=ok", body)
	}
}

func TestScanLifecycleEndToEnd(t *testing.T) {
	srv := newTestServer(t)

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "index.js"), []byte(`import { helper } from "./utils";`), 0o644); err != nil {
		t.Fatalf("write index.js: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "utils.js"), []byte(`export function helper() {}`), 0o644); err != nil {
		t.Fatalf("write utils.js: %v", err)
	}

	// POST /api/scan
	scanBody := `{"folderPath": "` + strings.ReplaceAll(root, `\`, `\\`) + `"}`
	res, err := http.Post(srv.URL+"/api/scan", "application/json", strings.NewReader(scanBody))
	if err != nil {
		t.Fatalf("POST /api/scan: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("POST /api/scan status = %d, want 201", res.StatusCode)
	}
	if origin := res.Header.Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin = %q, want http://localhost:3000", origin)
	}

	var scanResult struct {
		Scan struct {
			ID         int64 `json:"id"`
			TotalFiles int   `json:"totalFiles"`
		} `json:"scan"`
		Files        []map[string]any `json:"files"`
		Dependencies []map[string]any `json:"dependencies"`
	}
	if err := json.NewDecoder(res.Body).Decode(&scanResult); err != nil {
		t.Fatalf("decode scan response: %v", err)
	}
	res.Body.Close()

	if scanResult.Scan.TotalFiles != 2 {
		t.Errorf("TotalFiles = %d, want 2", scanResult.Scan.TotalFiles)
	}
	if len(scanResult.Files) != 2 {
		t.Errorf("len(Files) = %d, want 2", len(scanResult.Files))
	}
	if len(scanResult.Dependencies) != 1 {
		t.Errorf("len(Dependencies) = %d, want 1", len(scanResult.Dependencies))
	}

	scanID := scanResult.Scan.ID

	// GET /api/scans
	res, err = http.Get(srv.URL + "/api/scans")
	if err != nil {
		t.Fatalf("GET /api/scans: %v", err)
	}
	var list struct {
		Scans []struct {
			ID int64 `json:"id"`
		} `json:"scans"`
	}
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		t.Fatalf("decode scans list: %v", err)
	}
	res.Body.Close()
	if len(list.Scans) != 1 || list.Scans[0].ID != scanID {
		t.Errorf("GET /api/scans = %+v, want single scan with id %d", list.Scans, scanID)
	}

	// GET /api/scans/:id
	res, err = http.Get(srv.URL + "/api/scans/" + strconv.FormatInt(scanID, 10))
	if err != nil {
		t.Fatalf("GET /api/scans/:id: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET /api/scans/:id status = %d, want 200", res.StatusCode)
	}
	res.Body.Close()

	// GET /api/scans/:id for a missing scan
	res, err = http.Get(srv.URL + "/api/scans/999999")
	if err != nil {
		t.Fatalf("GET /api/scans/999999: %v", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("GET /api/scans/999999 status = %d, want 404", res.StatusCode)
	}
	res.Body.Close()
}

func TestScanEndpointRejectsMissingFolder(t *testing.T) {
	srv := newTestServer(t)

	res, err := http.Post(srv.URL+"/api/scan", "application/json", strings.NewReader(`{"folderPath": "/does/not/exist"}`))
	if err != nil {
		t.Fatalf("POST /api/scan: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}
