package dtos

import (
	"errors"
	"strings"
)

type StoryRequest struct {
	Prompt string `json:"prompt"`
	Topic  string `json:"topic"`
}

func (obj StoryRequest) Validate() error {
	if strings.TrimSpace(obj.Prompt) == "" && strings.TrimSpace(obj.Topic) == "" {
		return errors.New("prompt or topic is required")
	}
	return nil
}
