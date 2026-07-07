package utils

import "testing"

func TestBuildStoryPromptIncludesBanglaInstruction(t *testing.T) {
	prompt := "একটি ছোট্ট রোমাঞ্চকর গল্প"
	instruction := buildStoryPrompt(prompt)

	if instruction == "" {
		t.Fatal("expected story prompt to be non-empty")
	}

	if !contains(instruction, "বাংলা") {
		t.Fatalf("expected bangla instruction in prompt, got %q", instruction)
	}

	if !contains(instruction, prompt) {
		t.Fatalf("expected original prompt to be included, got %q", instruction)
	}
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && (s == substr || containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
