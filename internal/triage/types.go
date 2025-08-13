package triage

// Finding je jedinstveni tip nalaza (koristi se i u remote i u local režimu).
type Finding struct {
	File      string
	Type      string // bug|enhancement|none
	Title     string
	Details   string
	Severity  string // low|medium|high
	LineHints string
}
