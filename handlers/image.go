package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/dtos"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func GenerateImage(c *fiber.Ctx) error {
	var payload dtos.GenerateImageRequest
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

	// Auto-translate Bangla prompts so flux (which understands English best)
	// produces accurate images. If the prompt is already English this is a
	// no-op — Bangla detection is a cheap Unicode-range scan.
	originalPrompt := payload.Prompt
	wasTranslated := false
	if utils.HasBanglaScript(payload.Prompt) {
		translated, err := utils.TranslateBanglaToEnglish(payload.Prompt)
		if err != nil {
			return response.BadGatewayException(
				c,
				"Failed to translate Bangla prompt",
				err.Error(),
			)
		}
		payload.Prompt = translated
		wasTranslated = true
	}

	// Pollinations.ai — free, no-auth image generation.
	// Endpoint: GET https://image.pollinations.ai/prompt/{prompt}
	// Docs: https://github.com/pollinations/pollinations/blob/master/APIDOCS.md
	baseURL := os.Getenv("IMAGE_GENERATION")
	if baseURL == "" {
		baseURL = "https://image.pollinations.ai"
	}

	model := payload.Model
	if model == "" {
		model = os.Getenv("IMAGE_MODEL")
	}
	if model == "" {
		model = "flux"
	}

	width := payload.Width
	if width == 0 {
		width = 1024
	}

	height := payload.Height
	if height == 0 {
		height = 1024
	}

	seed := payload.Seed
	if seed == 0 {
		seed = int(time.Now().UnixNano() & 0x7fffffff)
	}

	query := url.Values{}
	query.Set("model", model)
	query.Set("width", fmt.Sprintf("%d", width))
	query.Set("height", fmt.Sprintf("%d", height))
	query.Set("seed", fmt.Sprintf("%d", seed))
	query.Set("nologo", "true")
	query.Set("enhance", "true")

	endpoint := fmt.Sprintf(
		"%s/prompt/%s?%s",
		strings.TrimRight(baseURL, "/"),
		url.PathEscape(payload.Prompt),
		query.Encode(),
	)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to build image request", err.Error())
	}

	req.Header.Set("Accept", "image/jpeg, image/png, image/webp, */*")
	req.Header.Set("User-Agent", "HelloBanglaTTS-Backend/1.0")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return response.InternalServerErrorException(c, "Image generation request failed", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return response.InternalServerErrorException(
			c,
			"Image generation failed",
			fmt.Sprintf("upstream returned %s: %s", resp.Status, string(body)),
		)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to read image response", err.Error())
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	c.Status(fiber.StatusOK)
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", `inline; filename="generated.jpg"`)
	c.Set("Cache-Control", "no-store")
	c.Set("X-Image-Model", model)
	c.Set("X-Image-Seed", fmt.Sprintf("%d", seed))
	if wasTranslated {
		c.Set("X-Image-Prompt-Translated", "true")
		c.Set("X-Image-Original-Prompt", originalPrompt)
	} else {
		c.Set("X-Image-Prompt-Translated", "false")
	}
	c.Set("Access-Control-Expose-Headers", "Content-Type, Content-Disposition, X-Image-Model, X-Image-Seed, X-Image-Prompt-Translated, X-Image-Original-Prompt")

	return c.Send(imageData)
}

// func GenerateImage(c *fiber.Ctx) error {
// 	var payload dtos.GenerateImageRequest
// 	if err := c.BodyParser(&payload); err != nil {
// 		return response.BadRequestException(c, "Invalid request body", err.Error())
// 	}

// 	if err := payload.Validate(); err != nil {
// 		return response.ValidationException(
// 			c,
// 			"Invalid request",
// 			strings.Split(err.Error(), ";"),
// 		)
// 	}

// 	// Auto-translate Bangla prompts so flux (which understands English best)
// 	// produces accurate images. If the prompt is already English this is a
// 	// no-op — Bangla detection is a cheap Unicode-range scan.
// 	// originalPrompt := payload.Prompt
// 	// wasTranslated := false
// 	// if utils.HasBanglaScript(payload.Prompt) {
// 	// 	translated, err := utils.TranslateBanglaToEnglish(payload.Prompt)
// 	// 	if err != nil {
// 	// 		return response.BadGatewayException(
// 	// 			c,
// 	// 			"Failed to translate Bangla prompt",
// 	// 			err.Error(),
// 	// 		)
// 	// 	}
// 	// 	payload.Prompt = translated
// 	// 	wasTranslated = true
// 	// }

// 	// Pollinations.ai — free, no-auth image generation.
// 	// Endpoint: GET https://image.pollinations.ai/prompt/{prompt}
// 	// Docs: https://github.com/pollinations/pollinations/blob/master/APIDOCS.md
// 	baseURL := os.Getenv("IMAGE_GENERATION")
// 	if baseURL == "" {
// 		baseURL = "https://image.pollinations.ai"
// 	}

// 	model := payload.Model
// 	if model == "" {
// 		model = os.Getenv("IMAGE_MODEL")
// 	}
// 	if model == "" {
// 		model = "flux"
// 	}

// 	width := payload.Width
// 	if width == 0 {
// 		width = 1024
// 	}

// 	height := payload.Height
// 	if height == 0 {
// 		height = 1024
// 	}

// 	seed := payload.Seed
// 	if seed == 0 {
// 		seed = int(time.Now().UnixNano() & 0x7fffffff)
// 	}

// 	query := url.Values{}
// 	query.Set("model", model)
// 	query.Set("width", fmt.Sprintf("%d", width))
// 	query.Set("height", fmt.Sprintf("%d", height))
// 	query.Set("seed", fmt.Sprintf("%d", seed))
// 	query.Set("nologo", "true")
// 	query.Set("enhance", "true")

// 	endpoint := fmt.Sprintf(
// 		"%s/prompt/%s?%s",
// 		strings.TrimRight(baseURL, "/"),
// 		url.PathEscape(payload.Prompt),
// 		query.Encode(),
// 	)

// 	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
// 	if err != nil {
// 		return response.InternalServerErrorException(c, "Failed to build image request", err.Error())
// 	}

// 	req.Header.Set("Accept", "image/jpeg, image/png, image/webp, */*")
// 	req.Header.Set("User-Agent", "HelloBanglaTTS-Backend/1.0")

// 	client := &http.Client{Timeout: 120 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return response.InternalServerErrorException(c, "Image generation request failed", err.Error())
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
// 		return response.InternalServerErrorException(
// 			c,
// 			"Image generation failed",
// 			fmt.Sprintf("upstream returned %s: %s", resp.Status, string(body)),
// 		)
// 	}

// 	imageData, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return response.InternalServerErrorException(c, "Failed to read image response", err.Error())
// 	}

// 	contentType := resp.Header.Get("Content-Type")
// 	if contentType == "" {
// 		contentType = "image/jpeg"
// 	}

// 	c.Status(fiber.StatusOK)
// 	c.Set("Content-Type", contentType)
// 	c.Set("Content-Disposition", `inline; filename="generated.jpg"`)
// 	c.Set("Cache-Control", "no-store")
// 	c.Set("X-Image-Model", model)
// 	c.Set("X-Image-Seed", fmt.Sprintf("%d", seed))
// 	c.Set("X-Image-Original-Prompt", payload.Prompt)
// 	c.Set("Access-Control-Expose-Headers", "Content-Type, Content-Disposition, X-Image-Model, X-Image-Seed, X-Image-Prompt-Translated, X-Image-Original-Prompt")

// 	return c.Send(imageData)
// }
