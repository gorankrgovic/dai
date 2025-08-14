package triage

// Finding is a unique type of finding
type Finding struct {
	File      string
	Type      string // bug|enhancement|none
	Title     string
	Details   string
	Severity  string // low|medium|high
	LineHints string
}
