package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkFindsSourceFilesAndSkipsIgnoredDirs(t *testing.T) {
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

	write("main.go", "package main\n")
	write("pkg/util.ts", "export const x = 1;\n")
	write("README.md", "# not source\n")
	write("node_modules/dep/index.js", "module.exports = {};\n")
	write(".git/HEAD", "ref: refs/heads/main\n")

	found, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk returned error: %v", err)
	}

	got := map[string]string{}
	for _, f := range found {
		got[f.RelPath] = f.Language
	}

	if lang, ok := got["main.go"]; !ok || lang != "Go" {
		t.Errorf("expected main.go to be found as Go, got %q (found=%v)", lang, ok)
	}
	if lang, ok := got["pkg/util.ts"]; !ok || lang != "TypeScript" {
		t.Errorf("expected pkg/util.ts to be found as TypeScript, got %q (found=%v)", lang, ok)
	}
	if _, ok := got["README.md"]; ok {
		t.Errorf("README.md should not be treated as a source file")
	}
	if _, ok := got["node_modules/dep/index.js"]; ok {
		t.Errorf("node_modules should have been skipped entirely")
	}
	if len(found) != 2 {
		t.Errorf("expected exactly 2 source files, got %d: %+v", len(found), found)
	}
}
