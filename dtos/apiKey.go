package dtos

import validation "github.com/go-ozzo/ozzo-validation"

type CreateApiKeyDTO struct {
	Name string `json:"name"`
}

func (obj CreateApiKeyDTO) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
}

type UpdateApiKeyDTO struct {
	Name     string `json:"name,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}
