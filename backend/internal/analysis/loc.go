package analysis

import "strings"

// LineCounts holds the line-based metrics computed for one source file.
type LineCounts struct {
	LinesOfCode  int
	CommentLines int
	CommentRatio float64
}

// commentPrefixes maps a language to the token(s) that mark the start of a
// single-line comment. Block comments are intentionally not tracked here —
// a line-prefix scan is precise enough for a comment-ratio metric without
// needing a full per-language lexer.
var commentPrefixes = map[string][]string{
	"Go":         {"//"},
	"JavaScript": {"//"},
	"TypeScript": {"//"},
	"Java":       {"//"},
	"C":          {"//"},
	"C++":        {"//"},
	"C#":         {"//"},
	"Rust":       {"//"},
	"Swift":      {"//"},
	"Kotlin":     {"//"},
	"Scala":      {"//"},
	"Python":     {"#"},
	"Ruby":       {"#"},
	"Shell":      {"#"},
	"PHP":        {"//", "#"},
}

// CountLines computes lines-of-code (non-blank lines), comment lines (lines
// whose first non-whitespace characters are a line-comment marker for the
// given language), and the ratio of comment lines to lines of code.
func CountLines(source, language string) LineCounts {
	prefixes := commentPrefixes[language]

	var loc, comments int
	for _, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		loc++
		for _, prefix := range prefixes {
			if strings.HasPrefix(trimmed, prefix) {
				comments++
				break
			}
		}
	}

	var ratio float64
	if loc > 0 {
		ratio = float64(comments) / float64(loc)
	}

	return LineCounts{LinesOfCode: loc, CommentLines: comments, CommentRatio: ratio}
}
