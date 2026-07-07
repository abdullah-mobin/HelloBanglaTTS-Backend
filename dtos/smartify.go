package dtos

import validation "github.com/go-ozzo/ozzo-validation"

type SmartifyRequest struct {
	Text string `json:"text" validate:"required"`
}

func (obj SmartifyRequest) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Text, validation.Required),
	)
}
