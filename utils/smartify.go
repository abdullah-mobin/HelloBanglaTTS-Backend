package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// SmartifyBanglaText rewrites a Bengali paragraph with more elegant wording,
// polished structure, and a professional, beautiful tone while preserving meaning.
func SmartifyBanglaText(text string) (string, error) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "", fmt.Errorf("text is required")
	}

	if !HasBanglaScript(trimmed) {
		return trimmed, nil
	}

	baseURL := os.Getenv("TEXT_GENERATION")
	if baseURL == "" {
		baseURL = "https://text.pollinations.ai"
	}

	model := os.Getenv("SMARTIFY_MODEL")
	if model == "" {
		model = "openai-fast"
	}

	instruction := fmt.Sprintf(
		"Rewrite the following Bengali paragraph to make it more elegant, polished, meaningful, professional, and beautiful. Preserve the core meaning and intent, but improve the wording and flow. Keep it in Bengali. Output only the rewritten paragraph, with no quotes, no bullets, and no explanation.\n\nBengali paragraph:\n%s\n\nRewritten paragraph:",
		trimmed,
	)

	q := url.Values{}
	q.Set("model", model)
	q.Set("temperature", "0.8")

	endpoint := fmt.Sprintf(
		"%s/%s?%s",
		strings.TrimRight(baseURL, "/"),
		url.PathEscape(instruction),
		q.Encode(),
	)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("User-Agent", "HelloBanglaTTS-Backend/1.0")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("smartify upstream returned %s: %s", resp.Status, string(body))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	rewritten := strings.TrimSpace(string(raw))
	rewritten = strings.TrimLeft(rewritten, `"'`)
	if i := strings.Index(strings.ToLower(rewritten), "rewritten paragraph:"); i == 0 {
		rewritten = strings.TrimSpace(rewritten[len("rewritten paragraph:"):])
	}

	if rewritten == "" {
		return "", fmt.Errorf("smartify upstream returned an empty response")
	}

	return rewritten, nil
}
