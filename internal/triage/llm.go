package triage

import (
	"context"
	"encoding/json"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type modelOutput struct {
	Type      string `json:"type"` // bug|enhancement|none
	Title     string `json:"title"`
	Details   string `json:"details"`
	Severity  string `json:"severity"`   // low|medium|high (optional)
	LineHints string `json:"line_hints"` // e.g. "approx lines 120-140"
}

// ------- NEW: diff analiza --------

func analyzeDiff(ctx context.Context, apiKey, model, path string, diffBlocks []string) (Finding, error) {
	sys := `You are a senior code reviewer focused on DIFFS. Output STRICT JSON ONLY (no prose), schema:
{
  "type": "bug" | "enhancement" | "none",
  "title": "short one-line summary",
  "details": "short explanation for developers",
  "severity": "low|medium|high",
  "line_hints": "optional location hints"
}
Rules:
- You are given a unified diff (with minimal context).
- Focus on ADDED code (lines starting with '+'). Use surrounding context to reason.
- If nothing stands out, return "none". Keep it specific.`

	var b strings.Builder
	b.WriteString("FILE PATH: ")
	b.WriteString(path)
	b.WriteString("\n\n")
	b.WriteString("DIFF (unified):\n```diff\n")
	for _, block := range diffBlocks {
		b.WriteString(block)
		if !strings.HasSuffix(block, "\n") {
			b.WriteString("\n")
		}
	}
	b.WriteString("```\n")

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

	// --- robust parsing  ---
	out := modelOutput{}
	raw := ""
	if len(resp.Choices) > 0 {
		raw = strings.TrimSpace(resp.Choices[0].Message.Content)
	}

	// remove code-fence
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

	if !strings.HasPrefix(strings.TrimSpace(raw), "{") {
		if i := strings.Index(raw, "{"); i >= 0 {
			if j := strings.LastIndex(raw, "}"); j >= 0 && j >= i {
				raw = raw[i : j+1]
			}
		}
	}
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
