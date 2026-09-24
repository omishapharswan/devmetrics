package analysis

import "testing"

func edgeExists(edges []DependencyEdge, from, to string) bool {
	for _, e := range edges {
		if e.FromRelPath == from && e.ToRelPath == to {
			return true
		}
	}
	return false
}

func TestBuildDependencyGraphJSRelativeImport(t *testing.T) {
	files := []DependencySource{
		{RelPath: "src/index.js", Language: "JavaScript", Content: `import { helper } from "./utils";`},
		{RelPath: "src/utils.js", Language: "JavaScript", Content: `export function helper() {}`},
	}
	edges := BuildDependencyGraph(files)
	if !edgeExists(edges, "src/index.js", "src/utils.js") {
		t.Errorf("expected edge src/index.js -> src/utils.js, got %+v", edges)
	}
}

func TestBuildDependencyGraphJSIndexResolution(t *testing.T) {
	files := []DependencySource{
		{RelPath: "src/app.ts", Language: "TypeScript", Content: `import x from "./lib";`},
		{RelPath: "src/lib/index.ts", Language: "TypeScript", Content: `export default 1;`},
	}
	edges := BuildDependencyGraph(files)
	if !edgeExists(edges, "src/app.ts", "src/lib/index.ts") {
		t.Errorf("expected edge src/app.ts -> src/lib/index.ts, got %+v", edges)
	}
}

func TestBuildDependencyGraphJSRequireAndParentDir(t *testing.T) {
	files := []DependencySource{
		{RelPath: "src/nested/mod.js", Language: "JavaScript", Content: `const x = require("../shared");`},
		{RelPath: "src/shared.js", Language: "JavaScript", Content: `module.exports = {};`},
	}
	edges := BuildDependencyGraph(files)
	if !edgeExists(edges, "src/nested/mod.js", "src/shared.js") {
		t.Errorf("expected edge src/nested/mod.js -> src/shared.js, got %+v", edges)
	}
}

func TestBuildDependencyGraphJSThirdPartyImportUnresolved(t *testing.T) {
	files := []DependencySource{
		{RelPath: "src/index.js", Language: "JavaScript", Content: `import React from "react";`},
	}
	edges := BuildDependencyGraph(files)
	if len(edges) != 0 {
		t.Errorf("expected no edges for a bare package import, got %+v", edges)
	}
}

func TestBuildDependencyGraphPythonRelativeImport(t *testing.T) {
	files := []DependencySource{
		{RelPath: "pkg/main.py", Language: "Python", Content: "from .utils import helper\n"},
		{RelPath: "pkg/utils.py", Language: "Python", Content: "def helper():\n    pass\n"},
	}
	edges := BuildDependencyGraph(files)
	if !edgeExists(edges, "pkg/main.py", "pkg/utils.py") {
		t.Errorf("expected edge pkg/main.py -> pkg/utils.py, got %+v", edges)
	}
}

func TestBuildDependencyGraphPythonParentRelativeImport(t *testing.T) {
	files := []DependencySource{
		{RelPath: "pkg/sub/mod.py", Language: "Python", Content: "from ..models import Foo\n"},
		{RelPath: "pkg/models.py", Language: "Python", Content: "class Foo: pass\n"},
	}
	edges := BuildDependencyGraph(files)
	if !edgeExists(edges, "pkg/sub/mod.py", "pkg/models.py") {
		t.Errorf("expected edge pkg/sub/mod.py -> pkg/models.py, got %+v", edges)
	}
}

func TestBuildDependencyGraphPythonAbsoluteImport(t *testing.T) {
	files := []DependencySource{
		{RelPath: "app.py", Language: "Python", Content: "import pkg.utils\n"},
		{RelPath: "pkg/utils.py", Language: "Python", Content: "def helper():\n    pass\n"},
	}
	edges := BuildDependencyGraph(files)
	if !edgeExists(edges, "app.py", "pkg/utils.py") {
		t.Errorf("expected edge app.py -> pkg/utils.py, got %+v", edges)
	}
}

func TestBuildDependencyGraphPythonStdlibImportUnresolved(t *testing.T) {
	files := []DependencySource{
		{RelPath: "app.py", Language: "Python", Content: "import os\nimport sys\n"},
	}
	edges := BuildDependencyGraph(files)
	if len(edges) != 0 {
		t.Errorf("expected no edges for stdlib imports, got %+v", edges)
	}
}

func TestBuildDependencyGraphSkipsUnhandledLanguages(t *testing.T) {
	files := []DependencySource{
		{RelPath: "main.go", Language: "Go", Content: `import "fmt"`},
	}
	edges := BuildDependencyGraph(files)
	if len(edges) != 0 {
		t.Errorf("expected Go imports to be skipped, got %+v", edges)
	}
}

func TestBuildDependencyGraphNoSelfEdges(t *testing.T) {
	files := []DependencySource{
		{RelPath: "src/a.js", Language: "JavaScript", Content: `import "./a";`},
	}
	edges := BuildDependencyGraph(files)
	if len(edges) != 0 {
		t.Errorf("expected no self-referencing edge, got %+v", edges)
	}
}
