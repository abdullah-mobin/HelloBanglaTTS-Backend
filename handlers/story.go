package handlers

import (
	"strings"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/dtos"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/utils"
	"github.com/gofiber/fiber/v2"
)

// GenerateStory exposes a Bangla story generation endpoint at POST /api/v1/generate/story.
func GenerateStory(c *fiber.Ctx) error {
	var payload dtos.StoryRequest
	if err := c.BodyParser(&payload); err != nil {
		return response.BadRequestException(c, "Invalid request body", err.Error())
	}

	if err := payload.Validate(); err != nil {
		return response.ValidationException(c, "Invalid request", []string{err.Error()})
	}

	topic := strings.TrimSpace(payload.Prompt)
	if topic == "" {
		topic = strings.TrimSpace(payload.Topic)
	}

	story, err := utils.GenerateBanglaStory(topic)
	if err != nil {
		return response.InternalServerErrorException(c, "Story generation failed", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"story": story,
		},
	})
}
