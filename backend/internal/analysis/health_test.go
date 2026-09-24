package analysis

import "testing"

func TestCalculateHealthScoreCleanSimpleFile(t *testing.T) {
	got := CalculateHealthScore(HealthScoreInput{
		Complexity:   3,
		LinesOfCode:  50,
		CommentRatio: 0.1,
		IsDuplicate:  false,
	})
	// 100 base + 5 comment-ratio bonus, clamped to the 100 ceiling.
	if got != 100 {
		t.Errorf("got %v, want 100", got)
	}
}

func TestCalculateHealthScorePenalizesHighComplexity(t *testing.T) {
	low := CalculateHealthScore(HealthScoreInput{Complexity: 5, LinesOfCode: 50})
	high := CalculateHealthScore(HealthScoreInput{Complexity: 60, LinesOfCode: 50})
	if high >= low {
		t.Errorf("expected high complexity (%v) to score lower than low complexity (%v)", high, low)
	}
}

func TestCalculateHealthScorePenalizesLargeFiles(t *testing.T) {
	small := CalculateHealthScore(HealthScoreInput{Complexity: 3, LinesOfCode: 100})
	large := CalculateHealthScore(HealthScoreInput{Complexity: 3, LinesOfCode: 2000})
	if large >= small {
		t.Errorf("expected large file (%v) to score lower than small file (%v)", large, small)
	}
}

func TestCalculateHealthScorePenalizesDuplication(t *testing.T) {
	unique := CalculateHealthScore(HealthScoreInput{Complexity: 3, LinesOfCode: 50, IsDuplicate: false})
	duplicate := CalculateHealthScore(HealthScoreInput{Complexity: 3, LinesOfCode: 50, IsDuplicate: true})
	if duplicate != unique-15 {
		t.Errorf("expected duplicate score to be exactly 15 points lower, got unique=%v duplicate=%v", unique, duplicate)
	}
}

func TestCalculateHealthScoreNeverGoesNegative(t *testing.T) {
	got := CalculateHealthScore(HealthScoreInput{Complexity: 500, LinesOfCode: 5000, IsDuplicate: true})
	if got < 0 {
		t.Errorf("expected score clamped at 0, got %v", got)
	}
	// Complexity and size penalties are each capped (50 and 30), plus the
	// flat 15-point duplication penalty: 100 - 50 - 30 - 15 = 5.
	if got != 5 {
		t.Errorf("got %v, want 5", got)
	}
}
