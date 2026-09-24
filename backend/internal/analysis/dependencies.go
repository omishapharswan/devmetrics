package analysis

import (
	"path"
	"regexp"
	"sort"
	"strings"
)

// DependencySource is one file's content to extract import statements from,
// keyed by its project-root-relative path (forward-slash separated,
// matching scanner.FoundFile.RelPath).
type DependencySource struct {
	RelPath  string
	Language string
	Content  string
}

// DependencyEdge is a directed import: the file at FromRelPath imports the
// file at ToRelPath.
type DependencyEdge struct {
	FromRelPath string
	ToRelPath   string
}

// jsResolveExtensions are tried, in order, when a relative JS/TS import
// doesn't already include one.
var jsResolveExtensions = []string{".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs"}

var (
	// jsImportPattern matches `import ... from "x"`, side-effect
	// `import "x"`, `export ... from "x"`, and `require("x")`. Only the
	// text between quotes is captured; it is not itself validated as an
	// import path until resolution.
	jsImportPattern = regexp.MustCompile(`(?:import|export)(?:[^'"]*?from)?\s+['"]([^'"]+)['"]|require\(\s*['"]([^'"]+)['"]\s*\)`)

	pyFromImportPattern  = regexp.MustCompile(`^\s*from\s+(\.*[\w.]*)\s+import\b`)
	pyPlainImportPattern = regexp.MustCompile(`^\s*import\s+([\w.]+(?:\s*,\s*[\w.]+)*)`)
)

// BuildDependencyGraph scans each file's source for import/require
// statements and resolves them against the set of files present in the
// same batch, returning one edge per import that could be matched to a
// file in that set. Imports that resolve outside the batch (third-party
// packages, unresolvable module paths) are silently skipped — this is a
// best-effort, regex-based resolver, not a full module resolution engine.
//
// Go is intentionally not handled: Go import paths are full module paths
// (e.g. "github.com/x/y/pkg") that can't be mapped to a file on disk
// without knowing the scanned project's module name.
func BuildDependencyGraph(files []DependencySource) []DependencyEdge {
	existing := make(map[string]bool, len(files))
	for _, f := range files {
		existing[f.RelPath] = true
	}

	seen := map[DependencyEdge]bool{}
	var edges []DependencyEdge

	for _, f := range files {
		var raw []string
		switch f.Language {
		case "JavaScript", "TypeScript":
			raw = extractJSImports(f.Content)
		case "Python":
			raw = extractPyImports(f.Content)
		default:
			continue
		}

		for _, imp := range raw {
			var target string
			var ok bool
			switch f.Language {
			case "JavaScript", "TypeScript":
				target, ok = resolveJSImport(f.RelPath, imp, existing)
			case "Python":
				target, ok = resolvePyImport(imp, f.RelPath, existing)
			}
			if !ok || target == f.RelPath {
				continue
			}

			edge := DependencyEdge{FromRelPath: f.RelPath, ToRelPath: target}
			if seen[edge] {
				continue
			}
			seen[edge] = true
			edges = append(edges, edge)
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		if edges[i].FromRelPath != edges[j].FromRelPath {
			return edges[i].FromRelPath < edges[j].FromRelPath
		}
		return edges[i].ToRelPath < edges[j].ToRelPath
	})

	return edges
}

func extractJSImports(content string) []string {
	matches := jsImportPattern.FindAllStringSubmatch(content, -1)
	imports := make([]string, 0, len(matches))
	for _, m := range matches {
		switch {
		case m[1] != "":
			imports = append(imports, m[1])
		case m[2] != "":
			imports = append(imports, m[2])
		}
	}
	return imports
}

// resolveJSImport resolves a relative ("./x", "../x") import against the
// importing file's directory, trying each known JS/TS extension and
// directory-index form. Non-relative imports (bare package names) are
// assumed to be third-party and are never resolved.
func resolveJSImport(fromRelPath, importPath string, existing map[string]bool) (string, bool) {
	if !strings.HasPrefix(importPath, ".") {
		return "", false
	}

	dir := path.Dir(fromRelPath)
	joined := path.Clean(path.Join(dir, importPath))

	if existing[joined] {
		return joined, true
	}
	for _, ext := range jsResolveExtensions {
		if candidate := joined + ext; existing[candidate] {
			return candidate, true
		}
	}
	for _, ext := range jsResolveExtensions {
		if candidate := joined + "/index" + ext; existing[candidate] {
			return candidate, true
		}
	}
	return "", false
}

// extractPyImports returns each import target found in content, tagged
// with "from:" or "import:" so resolvePyImport knows which resolution
// rules (relative-dot handling vs. plain) to apply.
func extractPyImports(content string) []string {
	var imports []string
	for _, line := range strings.Split(content, "\n") {
		if m := pyFromImportPattern.FindStringSubmatch(line); m != nil {
			imports = append(imports, "from:"+m[1])
			continue
		}
		if m := pyPlainImportPattern.FindStringSubmatch(line); m != nil {
			for _, mod := range strings.Split(m[1], ",") {
				imports = append(imports, "import:"+strings.TrimSpace(mod))
			}
		}
	}
	return imports
}

// resolvePyImport resolves a tagged Python import (see extractPyImports)
// against the importing file's path. Relative imports ("from .x import y",
// "from ..x import y") are resolved from the importing file's directory;
// absolute imports ("import a.b", "from a.b import c") are resolved from
// the scan root, which only succeeds when the scanned folder is itself the
// top of the package hierarchy.
func resolvePyImport(tagged, fromRelPath string, existing map[string]bool) (string, bool) {
	kind, module, ok := strings.Cut(tagged, ":")
	if !ok {
		return "", false
	}

	if kind == "from" {
		dotCount := 0
		for dotCount < len(module) && module[dotCount] == '.' {
			dotCount++
		}
		remainder := module[dotCount:]

		if dotCount > 0 {
			if remainder == "" {
				// "from . import x" / "from .. import x": can't tell
				// which submodule x refers to without more context.
				return "", false
			}
			baseDir := path.Dir(fromRelPath)
			for i := 1; i < dotCount; i++ {
				baseDir = path.Dir(baseDir)
			}
			target := path.Join(append([]string{baseDir}, strings.Split(remainder, ".")...)...)
			return resolvePyCandidate(target, existing)
		}

		target := path.Join(strings.Split(remainder, ".")...)
		return resolvePyCandidate(target, existing)
	}

	target := path.Join(strings.Split(module, ".")...)
	return resolvePyCandidate(target, existing)
}

func resolvePyCandidate(target string, existing map[string]bool) (string, bool) {
	if existing[target+".py"] {
		return target + ".py", true
	}
	if existing[target+"/__init__.py"] {
		return target + "/__init__.py", true
	}
	return "", false
}
