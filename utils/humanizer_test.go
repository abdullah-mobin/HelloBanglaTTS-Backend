package utils

import "testing"

func TestHumanizeBanglaTextReturnsOriginalTextWithoutBanglaScript(t *testing.T) {
	text := "This is plain english text."

	rewritten, err := HumanizeBanglaText(text)
	if err != nil {
		t.Fatalf("expected no error for non-Bangla input, got %v", err)
	}

	if rewritten != text {
		t.Fatalf("expected original text to be returned, got %q", rewritten)
	}
}
