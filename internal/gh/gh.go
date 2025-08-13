package gh

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type issueReq struct {
	Title  string   `json:"title"`
	Body   string   `json:"body,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

type issueResp struct {
	HTMLURL string `json:"html_url"`
	Number  int    `json:"number"`
}

func CreateIssue(ctx context.Context, owner, repo, token, title, body string, labels []string) (string, int, error) {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
	reqBody := issueReq{Title: title, Body: body, Labels: labels}
	b, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return "", 0, err
	}
	setCommonHeaders(req, token)

	hc := &http.Client{Timeout: 20 * time.Second}
	resp, err := hc.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		return "", 0, fmt.Errorf("github create issue: %d: %s", resp.StatusCode, buf.String())
	}
	var out issueResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", 0, err
	}
	return out.HTMLURL, out.Number, nil
}

func EnsureLabels(ctx context.Context, owner, repo, token string, labels []string) error {
	if len(labels) == 0 {
		return nil
	}
	existing, err := listLabels(ctx, owner, repo, token)
	if err != nil {
		return err
	}
	need := missing(labels, existing)
	for _, name := range need {
		if err := createLabel(ctx, owner, repo, token, name, defaultColor(name)); err != nil {
			return err
		}
	}
	return nil
}

func listLabels(ctx context.Context, owner, repo, token string) (map[string]struct{}, error) {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/labels?per_page=100", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	setCommonHeaders(req, token)
	hc := &http.Client{Timeout: 15 * time.Second}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		return nil, fmt.Errorf("github list labels: %d: %s", resp.StatusCode, buf.String())
	}
	var arr []label
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		return nil, err
	}
	set := make(map[string]struct{}, len(arr))
	for _, l := range arr {
		set[strings.ToLower(strings.TrimSpace(l.Name))] = struct{}{}
	}
	return set, nil
}

func createLabel(ctx context.Context, owner, repo, token, name, color string) error {
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/labels", owner, repo)
	body := label{Name: name, Color: color}
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	setCommonHeaders(req, token)
	hc := &http.Client{Timeout: 15 * time.Second}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 201 Created or 422 Unprocessable (exists with different case) both OK
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusUnprocessableEntity {
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("github create label: %d: %s", resp.StatusCode, buf.String())
	}
	return nil
}

func missing(want []string, have map[string]struct{}) []string {
	out := []string{}
	for _, w := range want {
		if _, ok := have[strings.ToLower(strings.TrimSpace(w))]; !ok {
			out = append(out, w)
		}
	}
	return out
}

func defaultColor(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "bug":
		return "d73a4a"
	case "enhancement":
		return "a2eeef"
	case "question":
		return "d876e3"
	default:
		return "cccccc"
	}
}

func setCommonHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dai-cli/triage")
}
