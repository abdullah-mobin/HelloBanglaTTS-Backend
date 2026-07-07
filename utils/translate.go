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

// HasBanglaScript reports whether the string contains any character in the
// Bengali Unicode block (U+0980..U+09FF). This is a fast, dependency-free
// detection used to decide whether to auto-translate prompts.
func HasBanglaScript(s string) bool {
	for _, r := range s {
		if r >= 0x0980 && r <= 0x09FF {
			return true
		}
	}
	return false
}

// TranslateBanglaToEnglish converts a Bengali (Bangla) text prompt into an
// English prompt suitable for image models like flux. If the input is already
// English (no Bengali script), it is returned unchanged.
//
// It uses the free Pollinations.ai text API:
//
//	GET https://text.pollinations.ai/{prompt}
//
// No authentication is required for the anonymous tier (1 req / 15s).
func TranslateBanglaToEnglish(text string) (string, error) {
	if !HasBanglaScript(text) {
		return text, nil
	}

	baseURL := os.Getenv("TEXT_GENERATION")
	if baseURL == "" {
		baseURL = "https://text.pollinations.ai"
	}

	model := os.Getenv("TRANSLATE_MODEL")
	if model == "" {
		model = "openai-fast"
	}

	// A focused system-style instruction baked into the prompt. We ask the
	// model for ONLY the translation so we can use the raw response directly
	// without JSON parsing.
	instruction := fmt.Sprintf(
		"Translate the following Bengali text to English. Only output the English translation, with no quotes, no explanation, and no prefix.\n\nBengali: %s\n\nEnglish:",
		text,
	)

	q := url.Values{}
	q.Set("model", model)
	q.Set("temperature", "0.2")

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
		return "", fmt.Errorf("translation upstream returned %s: %s", resp.Status, string(body))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	translated := strings.TrimSpace(string(raw))
	// Strip a stray leading "English:" or quote characters if the model added them.
	translated = strings.TrimLeft(translated, `"'`)
	if i := strings.Index(strings.ToLower(translated), "english:"); i == 0 {
		translated = strings.TrimSpace(translated[len("english:"):])
	}

	if translated == "" {
		return "", fmt.Errorf("translation upstream returned an empty response")
	}

	return translated, nil
}