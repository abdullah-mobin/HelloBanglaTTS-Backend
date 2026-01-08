package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

// GenerateSlug converts a string to a URL-friendly slug
// - Converts to lowercase
// - Replaces spaces and underscores with dashes
// - Removes special characters except dashes
// - Removes consecutive dashes
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace spaces and underscores with dashes
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove all characters except alphanumeric and dashes
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Replace consecutive dashes with a single dash
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim dashes from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// GenerateRandomKey returns a secure random key of n bytes encoded as hex string
func GenerateRandomKey(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateSerial generates a unique serial number with the given prefix
func GenerateSerial(prefix string) string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return prefix + "-ERROR"
	}
	return prefix + "-" + strings.ToUpper(hex.EncodeToString(b))
}
