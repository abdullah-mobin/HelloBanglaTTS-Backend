package dtos

import validation "github.com/go-ozzo/ozzo-validation"

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (obj LoginDTO) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Email, validation.Required),
		validation.Field(&obj.Password, validation.Required),
	)
}

type RegisterDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (obj RegisterDTO) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Email, validation.Required),
		validation.Field(&obj.Password, validation.Required),
	)
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

func (obj RefreshTokenDTO) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.RefreshToken, validation.Required),
	)
}
