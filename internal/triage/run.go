package triage

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gorankrgovic/dai/internal/gh"
	"github.com/gorankrgovic/dai/internal/gitutil"
	"github.com/gorankrgovic/dai/internal/ignore"
)

func Run(ctx context.Context, opt Options) (*Result, error) {
	commit := strings.TrimSpace(opt.Commit)
	if commit == "" {
		h, err := gitutil.HeadCommit(opt.Root)
		if err != nil {
			return nil, fmt.Errorf("resolve HEAD: %w", err)
		}
		commit = h
	}

	ign, _ := ignore.Load(opt.IgnoreFile)

	// --- DIFF HUNKS ---
	fileDiffs, err := gitutil.DiffHunks(opt.Root, commit, opt.DiffContext)
	if err != nil {
		return nil, fmt.Errorf("diff hunks: %w", err)
	}
	// filter po ekstenzijama i ignore pravilima
	filtered := make([]gitutil.FileDiff, 0, len(fileDiffs))
	for _, fd := range fileDiffs {
		if fd.Binary {
			continue
		}
		if !hasAllowedExt(fd.Path, opt.IncludeExts) {
			continue
		}
		if ign != nil && ign.MatchesPath(fd.Path) {
			continue
		}
		if len(fd.Hunks) == 0 {
			continue
		}
		filtered = append(filtered, fd)
	}

	// LLM analiza po fajlu, ali šaljemo SAMO hunk-ove (kao jedan concatenated diff)
	findings := make([]Finding, 0, len(filtered))
	for _, fd := range filtered {
		blocks := make([]string, 0, len(fd.Hunks))
		for _, h := range fd.Hunks {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", h.OldStart, h.OldLines, h.NewStart, h.NewLines))
			for _, ln := range h.Lines {
				// ostavi i '-' i ' ' i '+' linije — modelu treba minimalni kontekst,
				// ali najveći fokus je na '+'
				sb.WriteString(ln)
				sb.WriteString("\n")
			}
			blocks = append(blocks, sb.String())
		}
		ff, err := analyzeDiff(ctx, opt.OpenAIKey, opt.Model, fd.Path, blocks)
		if err != nil {
			// non-fatal: preskačemo konkretan fajl
			continue
		}
		findings = append(findings, ff)
	}

	title, body, labels := summarize(commit, findings)
	if opt.DryRun {
		return &Result{Body: body}, nil
	}
	if len(findings) == 0 && !opt.AlwaysOpen {
		return &Result{Body: body, Skipped: true}, nil
	}

	if err := gh.EnsureLabels(ctx, opt.Owner, opt.Repo, opt.GitHubToken, labels); err != nil {
		return nil, fmt.Errorf("ensure labels: %w", err)
	}
	url, num, err := gh.CreateIssue(ctx, opt.Owner, opt.Repo, opt.GitHubToken, title, body, labels)
	if err != nil {
		return nil, fmt.Errorf("create issue: %w", err)
	}
	return &Result{URL: url, Number: num, Body: body}, nil
}

func hasAllowedExt(path string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}
	l := strings.ToLower(path)
	dot := strings.LastIndex(l, ".")
	if dot < 0 {
		return false
	}
	ext := l[dot:]
	for _, e := range exts {
		if strings.ToLower(strings.TrimSpace(e)) == ext {
			return true
		}
	}
	return false
}

func summarize(commit string, findings []Finding) (title, body string, labels []string) {
	if len(findings) == 0 {
		title = fmt.Sprintf("DAI Triage: commit %.8s (no candidate findings)", commit)
		body = fmt.Sprintf("Automated triage for commit `%s` at %s\n\n_No findings from diff hunks._\n", commit, time.Now().Format(time.RFC3339))
		labels = []string{"question"}
		return
	}

	var bugs, enh []Finding
	for _, f := range findings {
		switch strings.ToLower(f.Type) {
		case "bug":
			bugs = append(bugs, f)
		case "enhancement":
			enh = append(enh, f)
		}
	}
	sevRank := map[string]int{"high": 0, "medium": 1, "low": 2, "": 3}
	sort.SliceStable(bugs, func(i, j int) bool {
		return sevRank[strings.ToLower(bugs[i].Severity)] < sevRank[strings.ToLower(bugs[j].Severity)]
	})

	var sb strings.Builder
	fmt.Fprintf(&sb, "Automated triage for commit `%s` at %s\n\n", commit, time.Now().Format(time.RFC3339))
	if len(bugs) > 0 {
		fmt.Fprintf(&sb, "## 🐞 Bugs (%d)\n", len(bugs))
		for i, f := range bugs {
			fmt.Fprintf(&sb, "%d) **%s** — `%s`\n", i+1, safeText(f.Title), f.File)
			if f.Severity != "" {
				fmt.Fprintf(&sb, "   - Severity: %s\n", strings.ToUpper(f.Severity))
			}
			if f.LineHints != "" {
				fmt.Fprintf(&sb, "   - Lines: %s\n", f.LineHints)
			}
			if f.Details != "" {
				fmt.Fprintf(&sb, "   - Details: %s\n", f.Details)
			}
		}
		fmt.Fprintln(&sb)
	}
	if len(enh) > 0 {
		fmt.Fprintf(&sb, "## ✨ Enhancements / Suggestions (%d)\n", len(enh))
		for i, f := range enh {
			fmt.Fprintf(&sb, "%d) **%s** — `%s`\n", i+1, safeText(f.Title), f.File)
			if f.Details != "" {
				fmt.Fprintf(&sb, "   - Details: %s\n", f.Details)
			}
		}
		fmt.Fprintln(&sb)
	}
	labels = nil
	if len(bugs) > 0 {
		labels = append(labels, "bug")
	}
	if len(enh) > 0 {
		labels = append(labels, "enhancement")
	}
	if len(labels) == 0 {
		labels = []string{"question"}
	}
	title = fmt.Sprintf("DAI Triage: commit %.8s — %d bug(s), %d suggestion(s)", commit, len(bugs), len(enh))
	return title, sb.String(), labels
}

func safeText(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
}
