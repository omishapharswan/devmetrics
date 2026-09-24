package analysis

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// duplicateWindowSize is the number of consecutive normalized lines hashed
// together to form one comparable block. Small enough to catch copy-pasted
// helper functions, large enough to avoid flagging trivial repeated lines
// like "}" or "return nil".
const duplicateWindowSize = 6

// SourceFile is one file's content to check for duplication against every
// other file in the same batch (including itself, for internal repeats).
type SourceFile struct {
	FileID  int64
	Content string
}

// DuplicateBlock is one window of lines whose normalized content hashes
// the same as a window found elsewhere in the batch.
type DuplicateBlock struct {
	FileID    int64
	Hash      string
	StartLine int // 1-indexed, inclusive
	EndLine   int // 1-indexed, inclusive
}

// normalizedLine is a non-blank source line stripped of surrounding
// whitespace, paired with its original 1-indexed line number so a match
// can still be reported at the right location.
type normalizedLine struct {
	text   string
	lineNo int
}

// DetectDuplicates hashes every duplicateWindowSize-line sliding window of
// non-blank, whitespace-normalized lines across all files and returns every
// window whose hash occurs more than once anywhere in the batch. Results
// are sorted by hash then by file/line so output is deterministic.
func DetectDuplicates(files []SourceFile) []DuplicateBlock {
	occurrencesByHash := map[string][]DuplicateBlock{}

	for _, f := range files {
		lines := normalizeLines(f.Content)
		for start := 0; start+duplicateWindowSize <= len(lines); start++ {
			window := lines[start : start+duplicateWindowSize]
			hash := hashBlock(window)
			occurrencesByHash[hash] = append(occurrencesByHash[hash], DuplicateBlock{
				FileID:    f.FileID,
				Hash:      hash,
				StartLine: window[0].lineNo,
				EndLine:   window[len(window)-1].lineNo,
			})
		}
	}

	var blocks []DuplicateBlock
	for _, occs := range occurrencesByHash {
		if len(occs) < 2 {
			continue
		}
		blocks = append(blocks, occs...)
	}

	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i].Hash != blocks[j].Hash {
			return blocks[i].Hash < blocks[j].Hash
		}
		if blocks[i].FileID != blocks[j].FileID {
			return blocks[i].FileID < blocks[j].FileID
		}
		return blocks[i].StartLine < blocks[j].StartLine
	})

	return blocks
}

func normalizeLines(content string) []normalizedLine {
	var out []normalizedLine
	for i, raw := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		out = append(out, normalizedLine{text: trimmed, lineNo: i + 1})
	}
	return out
}

func hashBlock(lines []normalizedLine) string {
	h := sha256.New()
	for _, l := range lines {
		h.Write([]byte(l.text))
		h.Write([]byte{'\n'})
	}
	return hex.EncodeToString(h.Sum(nil))
}
