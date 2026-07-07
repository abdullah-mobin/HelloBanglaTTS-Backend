package dtos

import validation "github.com/go-ozzo/ozzo-validation"

type TranslateRequest struct {
	Text string `json:"text" validate:"required"`
}

func (obj TranslateRequest) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Text, validation.Required),
	)
}