package buildinfo

var (
	// Overridden via -ldflags at build time.
	Version = "dev"
	Commit  = ""
	Date    = "" // ISO-8601 UTC
)
