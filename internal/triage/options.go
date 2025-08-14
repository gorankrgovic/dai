package triage

type Options struct {
	Root         string
	Owner        string
	Repo         string
	GitHubToken  string
	OpenAIKey    string
	Model        string
	Commit       string
	IncludeExts  []string
	MaxFileBytes int64
	IgnoreFile   string
	DryRun       bool
	AlwaysOpen   bool
	DiffContext  int
}

type Result struct {
	URL     string
	Number  int
	Body    string
	Skipped bool
}
