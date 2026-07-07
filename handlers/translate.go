package handlers

import (
	"strings"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/dtos"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/utils"
	"github.com/gofiber/fiber/v2"
)

// TranslateBanglaPrompt exposes the Bangla→English translator as a
// standalone endpoint at POST /api/v1/generate/translate.
//
// If the input is already in English, the response echoes the original
// text and sets `translated: false` so the caller can skip work.
func TranslateBanglaPrompt(c *fiber.Ctx) error {
	var payload dtos.TranslateRequest
	if err := c.BodyParser(&payload); err != nil {
		return response.BadRequestException(c, "Invalid request body", err.Error())
	}

	if err := payload.Validate(); err != nil {
		return response.ValidationException(
			c,
			"Invalid request",
			strings.Split(err.Error(), ";"),
		)
	}

	hadBangla := utils.HasBanglaScript(payload.Text)

	translated, err := utils.TranslateBanglaToEnglish(payload.Text)
	if err != nil {
		return response.InternalServerErrorException(
			c,
			"Translation failed",
			err.Error(),
		)
	}

	return response.Ok(c, "Translation successful", fiber.Map{
		"original":   payload.Text,
		"translated": translated,
		"converted":  hadBangla && translated != payload.Text,
		"wasBangla":  hadBangla,
	})
}