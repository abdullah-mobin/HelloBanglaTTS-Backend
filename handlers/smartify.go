package handlers

import (
	"strings"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/dtos"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/utils"
	"github.com/gofiber/fiber/v2"
)

// SmartifyBanglaTextHandler rewrites a Bengali paragraph into a more elegant,
// polished, and professional version.
func SmartifyBanglaTextHandler(c *fiber.Ctx) error {
	var payload dtos.SmartifyRequest
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

	rewritten, err := utils.SmartifyBanglaText(payload.Text)
	if err != nil {
		return response.InternalServerErrorException(
			c,
			"Text smartifying failed",
			err.Error(),
		)
	}

	return response.Ok(c, "Text smartified successfully", fiber.Map{
		"original":  payload.Text,
		"rewritten": rewritten,
		"wasBangla": utils.HasBanglaScript(payload.Text),
		"converted": utils.HasBanglaScript(payload.Text) && rewritten != payload.Text,
	})
}
