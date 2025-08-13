package config

import (
	"errors"
	"os"
	"path/filepath"
)

func githubTokenPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".dai", "github_token"), nil
}

func SaveGitHubToken(tok string) error {
	p, err := githubTokenPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(tok), 0o600)
}

func LoadGitHubToken() (string, error) {
	p, err := githubTokenPath()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DeleteGitHubToken removes the stored token file entirely.
func DeleteGitHubToken() error {
	p, err := githubTokenPath()
	if err != nil {
		return err
	}
	// If it doesn't exist, return a friendly error so caller can surface it.
	if _, statErr := os.Stat(p); errors.Is(statErr, os.ErrNotExist) {
		return os.ErrNotExist
	}
	return os.Remove(p)
}

// GitHubTokenExists returns true if a non-empty token file is present.
func GitHubTokenExists() (bool, error) {
	p, err := githubTokenPath()
	if err != nil {
		return false, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return len(b) > 0, nil
}
