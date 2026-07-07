package dtos

import validation "github.com/go-ozzo/ozzo-validation"

type HumanizeRequest struct {
	Text string `json:"text" validate:"required"`
}

func (obj HumanizeRequest) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Text, validation.Required),
	)
}
