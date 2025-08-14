package gitutil

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type FileDiff struct {
	Path   string
	Hunks  []Hunk
	Binary bool
}

type Hunk struct {
	OldStart int
	OldLines int
	NewStart int
	NewLines int
	Lines    []string
}

func DiffHunks(dir, commit string, contextLines int) ([]FileDiff, error) {
	args := []string{"show", "--no-color", "--format="}
	if contextLines >= 0 {
		args = append(args, fmt.Sprintf("-U%d", contextLines))
	}
	args = append(args, commit)

	out, err := runGit(dir, args...)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(out) == "" {
		return nil, nil
	}

	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	var diffs []FileDiff
	var cur *FileDiff

	// regexi
	reDiffHeader := regexp.MustCompile(`^diff --git a/(.+) b/(.+)$`)
	reNewFile := regexp.MustCompile(`^\+\+\+ b/(.+)$`)
	reBinary := regexp.MustCompile(`^Binary files .* and .* differ$`)
	reHunk := regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)

	var curHunk *Hunk

	flushHunk := func() {
		if curHunk != nil && cur != nil {
			cur.Hunks = append(cur.Hunks, *curHunk)
			curHunk = nil
		}
	}

	for sc.Scan() {
		line := sc.Text()

		if m := reDiffHeader.FindStringSubmatch(line); m != nil {
			flushHunk()
			if cur != nil {
				diffs = append(diffs, *cur)
			}
			cur = &FileDiff{Path: m[2]}
			continue
		}

		if reBinary.MatchString(line) {
			if cur != nil {
				cur.Binary = true
			}
			continue
		}

		if m := reNewFile.FindStringSubmatch(line); m != nil {
			if cur != nil {
				cur.Path = m[1]
			}
			continue
		}

		if m := reHunk.FindStringSubmatch(line); m != nil {
			flushHunk()
			curHunk = &Hunk{
				OldStart: atoiDefault(m[1]),
				OldLines: atoiDefault(m[2]),
				NewStart: atoiDefault(m[3]),
				NewLines: atoiDefault(m[4]),
				Lines:    make([]string, 0, 64),
			}
			continue
		}

		if curHunk != nil {
			if len(line) > 0 {
				switch line[0] {
				case ' ', '+', '-':
					curHunk.Lines = append(curHunk.Lines, line)
					continue
				}
			}
			flushHunk()
		}
	}
	flushHunk()
	if cur != nil {
		diffs = append(diffs, *cur)
	}
	return diffs, nil
}

func atoiDefault(s string) int {
	if s == "" {
		return 0
	}
	var n int
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < '0' || ch > '9' {
			return n
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
