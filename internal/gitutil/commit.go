package gitutil

import (
	"fmt"
	"strings"
)

// HeadCommit vraća hash HEAD commita, radeći u zadatom repo root-u.
func HeadCommit(dir string) (string, error) {
	out, err := runGit(dir, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// FilesInCommit vraća listu fajlova iz commita.
func FilesInCommit(dir, commit string) ([]string, error) {
	out, err := runGit(dir, "show", "--name-only", "--pretty=", commit)
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	lines := strings.Split(out, "\n")
	files := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			files = append(files, l)
		}
	}
	return files, nil
}

// FileAtCommit vraća sadržaj fajla iz commita, uz opciono trunciranje po veličini.
func FileAtCommit(dir, commit, path string, maxBytes int64) (content string, truncated bool, err error) {
	out, err := runGit(dir, "show", fmt.Sprintf("%s:%s", commit, path))
	if err != nil {
		return "", false, err
	}
	b := []byte(out)
	if maxBytes > 0 && int64(len(b)) > maxBytes {
		return string(b[:maxBytes]), true, nil
	}
	return out, false, nil
}
