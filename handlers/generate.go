package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/dtos"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/gofiber/fiber/v2"
)

func GenerateTTS(c *fiber.Ctx) error {
	var payload dtos.GenerateTTSRequest
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

	// Prepare Python request
	body, _ := json.Marshal(map[string]string{
		"text": payload.Text,
	})

	req, err := http.NewRequest(
		http.MethodPost,
		"http://127.0.0.1:8000/female",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response.InternalServerErrorException(
			c,
			"TTS generation failed",
			resp.Status,
		)
	}

	// 🔥 READ FULL FILE
	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 🔥 Now Go owns the response
	c.Status(fiber.StatusOK)
	c.Set("Content-Type", "audio/wav")
	c.Set("Content-Disposition", `inline; filename="tts.wav"`)
	c.Set("Cache-Control", "no-store")
	c.Set("Access-Control-Expose-Headers", "Content-Type, Content-Disposition")

	return c.Send(audioBytes)
}
