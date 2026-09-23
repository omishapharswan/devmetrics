package analysis

import "regexp"

// decisionPointPatterns match language-agnostic control-flow keywords and
// short-circuit boolean operators. Each match is one additional branch
// through the code, per McCabe's cyclomatic complexity formula:
// M = decision points + 1.
//
// "if" alone is sufficient to catch "else if" (the "if" substring still
// matches on a word boundary), so no separate else-if pattern is needed.
var decisionPointPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\bif\b`),
	regexp.MustCompile(`\belif\b`),
	regexp.MustCompile(`\bfor\b`),
	regexp.MustCompile(`\bforeach\b`),
	regexp.MustCompile(`\bwhile\b`),
	regexp.MustCompile(`\bcase\b`),
	regexp.MustCompile(`\bcatch\b`),
	regexp.MustCompile(`\bexcept\b`),
	regexp.MustCompile(`\bwhen\b`),
	regexp.MustCompile(`&&`),
	regexp.MustCompile(`\|\|`),
}

// CalculateComplexity computes a McCabe-style cyclomatic complexity score
// for a block of source code by counting decision points and adding 1 for
// the single linear path through the code.
func CalculateComplexity(source string) int {
	complexity := 1
	for _, pattern := range decisionPointPatterns {
		complexity += len(pattern.FindAllStringIndex(source, -1))
	}
	return complexity
}
