package analysis

import (
	"regexp"
	"strings"
)

// FunctionComplexity is the cyclomatic complexity computed for a single
// detected function or method within a source file.
type FunctionComplexity struct {
	Name       string
	StartLine  int // 1-indexed, inclusive
	EndLine    int // 1-indexed, inclusive
	Complexity int
}

var (
	declPattern = regexp.MustCompile(
		`^\s*(?:[\w.<>\[\]]+\s+)*(?:func|function|def|fn)\s+(?:\([^)]*\)\s*)?([A-Za-z_]\w*)\s*[(:]`,
	)
	arrowPattern = regexp.MustCompile(
		`^\s*(?:export\s+)?(?:const|let|var)\s+([A-Za-z_]\w*)\s*=\s*(?:async\s*)?\([^)]*\)\s*(?::[^=]+)?=>\s*\{?`,
	)
	pyDefPattern = regexp.MustCompile(`^(\s*)def\s+([A-Za-z_]\w*)\s*\(`)
)

// DetectFunctions splits source into per-function slices — using brace
// nesting for brace-delimited languages and indentation for Python — and
// computes each function's cyclomatic complexity independently.
func DetectFunctions(source, language string) []FunctionComplexity {
	lines := strings.Split(source, "\n")
	if language == "Python" {
		return detectByIndentation(lines)
	}
	return detectByBraces(lines)
}

func matchFunctionName(line string) (string, bool) {
	if m := declPattern.FindStringSubmatch(line); m != nil {
		return m[1], true
	}
	if m := arrowPattern.FindStringSubmatch(line); m != nil {
		return m[1], true
	}
	return "", false
}

// detectByBraces finds function signatures line-by-line, then tracks
// brace depth from the signature line to locate the matching closing
// brace as the function's end.
func detectByBraces(lines []string) []FunctionComplexity {
	var results []FunctionComplexity

	i := 0
	for i < len(lines) {
		name, ok := matchFunctionName(lines[i])
		if !ok {
			i++
			continue
		}

		start := i
		depth := 0
		opened := false
		end := start

		for j := start; j < len(lines); j++ {
			for _, ch := range lines[j] {
				switch ch {
				case '{':
					depth++
					opened = true
				case '}':
					depth--
				}
			}
			end = j
			if opened && depth <= 0 {
				break
			}
		}

		body := strings.Join(lines[start:end+1], "\n")
		results = append(results, FunctionComplexity{
			Name:       name,
			StartLine:  start + 1,
			EndLine:    end + 1,
			Complexity: CalculateComplexity(body),
		})

		i = end + 1
	}

	return results
}

// detectByIndentation handles Python, where a function body is every
// following line indented further than its "def" line.
func detectByIndentation(lines []string) []FunctionComplexity {
	var results []FunctionComplexity

	i := 0
	for i < len(lines) {
		m := pyDefPattern.FindStringSubmatch(lines[i])
		if m == nil {
			i++
			continue
		}
		indent := len(m[1])
		name := m[2]
		start := i
		end := i

		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "" {
				end = j
				continue
			}
			lineIndent := len(lines[j]) - len(strings.TrimLeft(lines[j], " \t"))
			if lineIndent <= indent {
				break
			}
			end = j
		}

		body := strings.Join(lines[start:end+1], "\n")
		results = append(results, FunctionComplexity{
			Name:       name,
			StartLine:  start + 1,
			EndLine:    end + 1,
			Complexity: CalculateComplexity(body),
		})

		i = end + 1
	}

	return results
}
