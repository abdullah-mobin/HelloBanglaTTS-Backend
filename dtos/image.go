package dtos

import validation "github.com/go-ozzo/ozzo-validation"

type GenerateImageRequest struct {
	Prompt string `json:"prompt" validate:"required"`
	Model  string `json:"model,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	Seed   int    `json:"seed,omitempty"`
}

func (obj GenerateImageRequest) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Prompt, validation.Required),
		validation.Field(&obj.Width, validation.Min(64), validation.Max(2048)),
		validation.Field(&obj.Height, validation.Min(64), validation.Max(2048)),
	)
}
