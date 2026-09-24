package analysis

// HealthScoreInput bundles the per-file metrics used to compute a single
// 0-100 code health score.
type HealthScoreInput struct {
	Complexity   int
	LinesOfCode  int
	CommentRatio float64
	IsDuplicate  bool
}

// CalculateHealthScore combines complexity, size, documentation, and
// duplication into a single 0-100 score, where 100 is healthiest:
//   - Complexity above 10 is penalized, tapering linearly to a 50-point
//     cap at a complexity of 50 (McCabe's own "consider refactoring"
//     threshold is commonly cited as 10).
//   - Files over 400 lines are penalized, tapering linearly to a 30-point
//     cap at 1000 lines — long files are harder to hold in your head
//     regardless of their per-line complexity.
//   - A comment ratio between 5% and 40% earns a small bonus; outside that
//     band (undocumented, or mostly comments) earns none.
//   - A file that participates in a detected duplicate block takes a flat
//     penalty, since duplication is a maintenance cost independent of any
//     single file's own complexity or size.
func CalculateHealthScore(in HealthScoreInput) float64 {
	score := 100.0

	if in.Complexity > 10 {
		penalty := float64(in.Complexity-10) / 40.0 * 50.0
		if penalty > 50 {
			penalty = 50
		}
		score -= penalty
	}

	if in.LinesOfCode > 400 {
		penalty := float64(in.LinesOfCode-400) / 600.0 * 30.0
		if penalty > 30 {
			penalty = 30
		}
		score -= penalty
	}

	if in.CommentRatio >= 0.05 && in.CommentRatio <= 0.4 {
		score += 5
	}

	if in.IsDuplicate {
		score -= 15
	}

	switch {
	case score < 0:
		score = 0
	case score > 100:
		score = 100
	}

	return score
}
