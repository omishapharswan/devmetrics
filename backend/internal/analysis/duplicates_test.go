package analysis

import "testing"

const repeatedBlock = `line a
line b
line c
line d
line e
line f`

func TestDetectDuplicatesAcrossFiles(t *testing.T) {
	fileA := SourceFile{FileID: 1, Content: "func one() {\n" + repeatedBlock + "\n}\n"}
	fileB := SourceFile{FileID: 2, Content: "func two() {\n" + repeatedBlock + "\n}\n"}

	blocks := DetectDuplicates([]SourceFile{fileA, fileB})

	if len(blocks) == 0 {
		t.Fatal("expected at least one duplicate block, got none")
	}

	fileIDs := map[int64]bool{}
	for _, b := range blocks {
		fileIDs[b.FileID] = true
	}
	if !fileIDs[1] || !fileIDs[2] {
		t.Errorf("expected duplicate blocks reported in both files, got %+v", blocks)
	}
}

func TestDetectDuplicatesNoMatchWhenContentDiffers(t *testing.T) {
	fileA := SourceFile{FileID: 1, Content: "alpha\nbeta\ngamma\ndelta\nepsilon\nzeta\n"}
	fileB := SourceFile{FileID: 2, Content: "one\ntwo\nthree\nfour\nfive\nsix\n"}

	blocks := DetectDuplicates([]SourceFile{fileA, fileB})
	if len(blocks) != 0 {
		t.Errorf("expected no duplicate blocks, got %+v", blocks)
	}
}

func TestDetectDuplicatesIgnoresShortFiles(t *testing.T) {
	fileA := SourceFile{FileID: 1, Content: "only\nthree\nlines\n"}
	fileB := SourceFile{FileID: 2, Content: "only\nthree\nlines\n"}

	blocks := DetectDuplicates([]SourceFile{fileA, fileB})
	if len(blocks) != 0 {
		t.Errorf("expected no duplicate blocks below the window size, got %+v", blocks)
	}
}

func TestDetectDuplicatesWithinSameFile(t *testing.T) {
	content := repeatedBlock + "\n\n" + repeatedBlock + "\n"
	blocks := DetectDuplicates([]SourceFile{{FileID: 1, Content: content}})

	if len(blocks) < 2 {
		t.Fatalf("expected repeated block within the same file to be flagged, got %+v", blocks)
	}
	for _, b := range blocks {
		if b.FileID != 1 {
			t.Errorf("expected all blocks to belong to file 1, got %+v", b)
		}
	}
}

func TestDetectDuplicatesResultsAreSorted(t *testing.T) {
	fileA := SourceFile{FileID: 5, Content: repeatedBlock + "\n"}
	fileB := SourceFile{FileID: 2, Content: repeatedBlock + "\n"}

	blocks := DetectDuplicates([]SourceFile{fileA, fileB})
	if len(blocks) != 2 {
		t.Fatalf("expected exactly 2 blocks, got %d: %+v", len(blocks), blocks)
	}
	if blocks[0].FileID != 2 || blocks[1].FileID != 5 {
		t.Errorf("expected blocks sorted by file ID (2 before 5), got %+v", blocks)
	}
}
