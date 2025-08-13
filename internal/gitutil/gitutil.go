package gitutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %v: %v (%s)", args, err, strings.TrimSpace(errb.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

func IsRepo(dir string) bool {
	_, err := runGit(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil
}

func RepoRoot(dir string) (string, error) {
	return runGit(dir, "rev-parse", "--show-toplevel")
}

func RemoteOriginURL(dir string) (string, bool) {
	u, err := runGit(dir, "remote", "get-url", "origin")
	if err != nil || strings.TrimSpace(u) == "" {
		return "", false
	}
	return u, true
}

// SSH ili HTTPS GitHub URL → owner/repo
func ParseOwnerRepo(remote string) (owner, repo string) {
	ssh := regexp.MustCompile(`git@[^:]+:([^/]+)/([^\.]+)(?:\.git)?$`)
	http := regexp.MustCompile(`https?://[^/]+/([^/]+)/([^\.]+)(?:\.git)?$`)
	if m := ssh.FindStringSubmatch(remote); len(m) == 3 {
		return m[1], m[2]
	}
	if m := http.FindStringSubmatch(remote); len(m) == 3 {
		return m[1], m[2]
	}
	return "", ""
}
