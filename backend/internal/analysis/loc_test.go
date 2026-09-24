package analysis

import "testing"

func TestCountLinesGo(t *testing.T) {
	source := `package main

// Add returns the sum of a and b.
func Add(a, b int) int {
	return a + b
}
`
	got := CountLines(source, "Go")
	if got.LinesOfCode != 5 {
		t.Errorf("LinesOfCode = %d, want 5", got.LinesOfCode)
	}
	if got.CommentLines != 1 {
		t.Errorf("CommentLines = %d, want 1", got.CommentLines)
	}
	wantRatio := 1.0 / 5.0
	if got.CommentRatio != wantRatio {
		t.Errorf("CommentRatio = %v, want %v", got.CommentRatio, wantRatio)
	}
}

func TestCountLinesPython(t *testing.T) {
	source := `# greet prints a greeting
def greet(name):
    print("hi " + name)
`
	got := CountLines(source, "Python")
	if got.LinesOfCode != 3 {
		t.Errorf("LinesOfCode = %d, want 3", got.LinesOfCode)
	}
	if got.CommentLines != 1 {
		t.Errorf("CommentLines = %d, want 1", got.CommentLines)
	}
}

func TestCountLinesIgnoresBlankLines(t *testing.T) {
	source := "line one\n\n\nline two\n"
	got := CountLines(source, "Go")
	if got.LinesOfCode != 2 {
		t.Errorf("LinesOfCode = %d, want 2", got.LinesOfCode)
	}
}

func TestCountLinesEmptySourceHasZeroRatio(t *testing.T) {
	got := CountLines("", "Go")
	if got.LinesOfCode != 0 || got.CommentLines != 0 || got.CommentRatio != 0 {
		t.Errorf("got %+v, want all zero", got)
	}
}

func TestCountLinesUnknownLanguageCountsNoComments(t *testing.T) {
	source := "// looks like a comment\ncode here\n"
	got := CountLines(source, "Cobol")
	if got.LinesOfCode != 2 {
		t.Errorf("LinesOfCode = %d, want 2", got.LinesOfCode)
	}
	if got.CommentLines != 0 {
		t.Errorf("CommentLines = %d, want 0 for unrecognized language", got.CommentLines)
	}
}
