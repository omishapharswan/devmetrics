package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// FoundFile is one recognized source file discovered while walking a
// folder, before any complexity/LOC/duplication analysis has run.
type FoundFile struct {
	AbsPath   string // absolute path on disk
	RelPath   string // path relative to the scanned root, forward-slash separated
	Extension string
	Language  string
	SizeBytes int64
}

// Walk recursively scans rootPath and returns every recognized source
// file found, skipping directories in dirsToSkip and any hidden
// directory (name starting with '.') other than the root itself.
func Walk(rootPath string) ([]FoundFile, error) {
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("resolve root path: %w", err)
	}

	var found []FoundFile

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip entries we can't read (permission errors, broken
			// symlinks) instead of failing the whole scan.
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		name := d.Name()

		if d.IsDir() {
			if path != absRoot && (dirsToSkip[name] || strings.HasPrefix(name, ".")) {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(name))
		lang, ok := LanguageForExtension(ext)
		if !ok {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(absRoot, path)
		if err != nil {
			relPath = path
		}

		found = append(found, FoundFile{
			AbsPath:   path,
			RelPath:   filepath.ToSlash(relPath),
			Extension: ext,
			Language:  lang,
			SizeBytes: fileInfo.Size(),
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", absRoot, err)
	}

	return found, nil
}
