package ignore

import (
	"os"

	igi "github.com/sabhiram/go-gitignore"
)

type Matcher interface {
	MatchesPath(f string) bool
}

func Load(path string) (Matcher, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ign := igi.CompileIgnoreLines(splitLines(string(b))...)
	return ign, nil
}

func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, trimCR(s[start:i]))
			start = i + 1
		}
	}
	if start <= len(s)-1 {
		lines = append(lines, trimCR(s[start:]))
	}
	return lines
}

func trimCR(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\r' {
		return s[:len(s)-1]
	}
	return s
}
