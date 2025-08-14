package triage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type LocalFinding = Finding

func AnalyzeLocal(ctx context.Context, apiKey, model, absPath string, maxBytes int64) (LocalFinding, bool, error) {
	code, truncated, err := readWithLimit(absPath, maxBytes)
	if err != nil {
		return LocalFinding{}, false, err
	}
	ff, err := analyzeSingleFile(ctx, apiKey, model, absPath, code, truncated)
	return LocalFinding(ff), truncated, err
}

func readWithLimit(path string, maxBytes int64) (string, bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}
	if maxBytes > 0 && int64(len(b)) > maxBytes {
		return string(b[:maxBytes]), true, nil
	}
	return string(b), false, nil
}

func analyzeSingleFile(ctx context.Context, apiKey, model, path, code string, truncated bool) (Finding, error) {
	sys := `You are a senior code reviewer. Output STRICT JSON ONLY (no prose), following schema:
{
  "type": "bug" | "enhancement" | "none",
  "title": "short one-line summary",
  "details": "short explanation for developers",
  "severity": "low|medium|high",
  "line_hints": "optional location hints"
}
Important:
- Treat any syntax/parse error, typo (unknown identifier, misplaced token), missing import, wrong API usage, or type error as a "bug".
- If the code would not compile/run as-is (e.g., malformed arrow function or callback), classify it as "bug".
- Prefer being specific and reference approximate lines in "line_hints".
- Return "none" ONLY if nothing problematic is present.`

	var b strings.Builder
	b.WriteString("FILE PATH: ")
	b.WriteString(filepath.ToSlash(path))
	b.WriteString("\n\n")
	if truncated {
		b.WriteString("(Note: content truncated to size limit)\n")
	}
	b.WriteString("CODE:\n```")
	b.WriteString(detectFence(path))
	b.WriteString("\n")
	b.WriteString(code)
	b.WriteString("\n```")

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: sys},
			{Role: openai.ChatMessageRoleUser, Content: b.String()},
		},
		Temperature: 0.1,
	})
	if err != nil {
		return Finding{File: path, Type: "none"}, err
	}
	out := modelOutput{}

	// raw text
	raw := ""
	if len(resp.Choices) > 0 {
		raw = strings.TrimSpace(resp.Choices[0].Message.Content)
	}

	// 1) code fences
	low := strings.ToLower(raw)
	if strings.HasPrefix(low, "```json") || strings.HasPrefix(low, "```") {
		raw = strings.TrimPrefix(raw, "```json")
		raw = strings.TrimPrefix(raw, "```JSON")
		raw = strings.TrimPrefix(raw, "```")
		raw = strings.TrimSpace(raw)
		if strings.HasSuffix(raw, "```") {
			raw = strings.TrimSuffix(raw, "```")
			raw = strings.TrimSpace(raw)
		}
	}

	// 2) remove JSON block
	if !strings.HasPrefix(strings.TrimSpace(raw), "{") {
		if i := strings.Index(raw, "{"); i >= 0 {
			if j := strings.LastIndex(raw, "}"); j >= 0 && j >= i {
				raw = raw[i : j+1]
			}
		}
	}

	// 3)  unmarshal
	_ = json.Unmarshal([]byte(raw), &out)

	typ := strings.ToLower(strings.TrimSpace(out.Type))
	if typ != "bug" && typ != "enhancement" && typ != "none" {
		typ = "none"
	}
	sev := strings.ToLower(strings.TrimSpace(out.Severity))
	if sev != "low" && sev != "medium" && sev != "high" {
		sev = ""
	}

	return Finding{
		File:      path,
		Type:      typ,
		Title:     strings.TrimSpace(out.Title),
		Details:   strings.TrimSpace(out.Details),
		Severity:  sev,
		LineHints: strings.TrimSpace(out.LineHints),
	}, nil
}

func detectFence(path string) string {
	l := strings.ToLower(path)
	switch {
	case strings.HasSuffix(l, ".go"):
		return "go"
	case strings.HasSuffix(l, ".js"), strings.HasSuffix(l, ".jsx"):
		return "javascript"
	case strings.HasSuffix(l, ".ts"), strings.HasSuffix(l, ".tsx"):
		return "ts"
	case strings.HasSuffix(l, ".vue"):
		return "vue"
	case strings.HasSuffix(l, ".php"):
		return "php"
	case strings.HasSuffix(l, ".py"):
		return "python"
	case strings.HasSuffix(l, ".rb"):
		return "ruby"
	case strings.HasSuffix(l, ".java"):
		return "java"
	case strings.HasSuffix(l, ".cs"):
		return "csharp"
	case strings.HasSuffix(l, ".rs"):
		return "rust"
	case strings.HasSuffix(l, ".c"), strings.HasSuffix(l, ".h"):
		return "c"
	case strings.HasSuffix(l, ".cpp"), strings.HasSuffix(l, ".cc"), strings.HasSuffix(l, ".cxx"), strings.HasSuffix(l, ".hpp"), strings.HasSuffix(l, ".hh"):
		return "cpp"
	default:
		return ""
	}
}
