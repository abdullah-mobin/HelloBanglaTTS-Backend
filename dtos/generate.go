package dtos

import validation "github.com/go-ozzo/ozzo-validation"

type GenerateTTSRequest struct {
	Text  string `json:"text"`
	Actor string `json:"actor"`
}

func (obj GenerateTTSRequest) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Text, validation.Required),
		validation.Field(&obj.Actor, validation.Required),
	)
}
