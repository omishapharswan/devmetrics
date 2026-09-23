package scanner

// languageByExtension maps a lowercased file extension (including the
// leading dot) to the language DevMetrics reports it as. Only
// extensions present here are treated as source files.
var languageByExtension = map[string]string{
	".go":    "Go",
	".js":    "JavaScript",
	".jsx":   "JavaScript",
	".mjs":   "JavaScript",
	".cjs":   "JavaScript",
	".ts":    "TypeScript",
	".tsx":   "TypeScript",
	".py":    "Python",
	".java":  "Java",
	".c":     "C",
	".h":     "C",
	".cpp":   "C++",
	".cc":    "C++",
	".hpp":   "C++",
	".cs":    "C#",
	".rb":    "Ruby",
	".php":   "PHP",
	".rs":    "Rust",
	".swift": "Swift",
	".kt":    "Kotlin",
	".kts":   "Kotlin",
	".scala": "Scala",
	".sh":    "Shell",
}

// dirsToSkip are directory names that are never descended into,
// regardless of depth, because they hold generated output or
// third-party code rather than a project's own source.
var dirsToSkip = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	".next":        true,
	"out":          true,
	"__pycache__":  true,
	".venv":        true,
	"venv":         true,
	".idea":        true,
	".vscode":      true,
	"bin":          true,
	"obj":          true,
	"target":       true,
	".cache":       true,
}

// LanguageForExtension returns the language name for a given file
// extension and whether that extension is recognized as source code.
func LanguageForExtension(ext string) (string, bool) {
	lang, ok := languageByExtension[ext]
	return lang, ok
}
