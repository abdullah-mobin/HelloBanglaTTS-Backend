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

func buildStoryPrompt(prompt string) string {
	trimmed := strings.TrimSpace(prompt)
	if trimmed == "" {
		trimmed = "একটি সুন্দর, ছোট, মজার বাংলা গল্প"
	}

	return fmt.Sprintf(
		"নিচের বিষয়ে একটি নতুন, সুন্দর, ছোট এবং আকর্ষণীয় বাংলা গল্প লেখো। গল্পটি 150 থেকে 250 শব্দের মধ্যে হওয়া উচিত। শুধু গল্পের লেখা দাও, কোনো শিরোনাম, ব্যাখ্যা বা তালিকা দিও না।\n\nবিষয়: %s\n\nগল্প:",
		trimmed,
	)
}

// GenerateBanglaStory generates a Bangla story using the free Pollinations.ai text API.
func GenerateBanglaStory(prompt string) (string, error) {
	baseURL := os.Getenv("TEXT_GENERATION")
	if baseURL == "" {
		baseURL = "https://text.pollinations.ai"
	}

	model := os.Getenv("STORY_MODEL")
	if model == "" {
		model = "openai-fast"
	}

	instruction := buildStoryPrompt(prompt)

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
		return "", fmt.Errorf("story upstream returned %s: %s", resp.Status, string(body))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	story := strings.TrimSpace(string(raw))
	story = strings.Trim(story, "\"'\n\r")
	if story == "" {
		return "", fmt.Errorf("story upstream returned an empty response")
	}

	return story, nil
}
