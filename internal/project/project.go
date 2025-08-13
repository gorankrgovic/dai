package project

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Project struct {
	Provider string `yaml:"provider"` // "github"
	Owner    string `yaml:"owner"`
	Repo     string `yaml:"repo"`
}

func Path(root string) string {
	return filepath.Join(root, ".dai", "project.yaml")
}

func Save(root string, p *Project) error {
	dir := filepath.Dir(Path(root))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	b, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	return os.WriteFile(Path(root), b, 0o600)
}

func Load(root string) (*Project, error) {
	b, err := os.ReadFile(Path(root))
	if err != nil {
		return nil, err
	}
	var p Project
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
