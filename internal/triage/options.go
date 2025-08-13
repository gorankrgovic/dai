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
	AlwaysOpen   bool // NEW: forsiraj kreiranje issue-a i kad nema nalaza
	DiffContext  int  // NEW: broj kontekst linija po hunku (npr. 3)
}

type Result struct {
	URL     string
	Number  int
	Body    string
	Skipped bool // NEW: true ako smo preskočili otvaranje issue-a (nema nalaza)
}
