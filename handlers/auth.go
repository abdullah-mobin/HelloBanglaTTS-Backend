package handlers

import (
	"context"
	"strings"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/dtos"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/repository"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/utils"
	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
)

// Login godoc
//
//	@Summary		Login to hello-bangla-tts
//	@Description	login to user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body	dtos.LoginDTO	true	"user login dto"
//	@Router			/auth/login [post]
func Login(c *fiber.Ctx) error {
	var self dtos.LoginDTO
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid request", errorsArr)
	}
	repo := repository.NewCredentialRepository()
	credential, err := repo.GetCredentialByEmail(context.Background(), self.Email)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to retrieve user", err.Error())
	}
	if credential == nil {
		return response.NotFoundException(c, "User not found", nil)
	}
	if !utils.ComparePassword(credential.Password, self.Password) {
		return response.UnauthorizedException(c, "Invalid credentials", nil)
	}
	tokenPayload := utils.TokenPayload{
		UserID: credential.UserID.Hex(),
	}
	accessToken, err := utils.GenerateAccessToken(tokenPayload)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to generate access token", err.Error())
	}
	refreshToken, err := utils.GenerateRefreshToken(tokenPayload)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to generate refresh token", err.Error())
	}
	err = repo.UpdateCredential(context.Background(), credential.ID.Hex(), bson.M{"refresh_token": refreshToken})
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to save refresh token", err.Error())
	}
	return response.Ok(c, "User Login successfully", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout godoc
//
//	@Summary		logout from hello-bangla-tts
//	@Description	logout from user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/auth/logout [post]
func Logout(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	repo := repository.NewCredentialRepository()
	err := repo.InvalidateRefreshTokens(context.Background(), userId.(string))
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to logout user", err.Error())
	}
	return response.Ok(c, "User Logout successfully", nil)
}

// Refresh token godoc
//
//	@Summary		Refresh token
//	@Description	refresh users access token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body	dtos.RefreshTokenDTO	true	"user login dto"
//	@Router			/auth/refresh [post]
func RefreshToken(c *fiber.Ctx) error {
	var self dtos.RefreshTokenDTO
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid request", errorsArr)
	}
	repo := repository.NewCredentialRepository()
	err, credential := repo.FindCredentialUsingRefreshToken(context.Background(), self.RefreshToken)

	if err != nil {
		return response.InternalServerErrorException(c, "Failed to retrieve credential", err.Error())
	}
	if credential == nil {
		return response.NotFoundException(c, "User not found", nil)
	}
	tokenPayload := utils.TokenPayload{
		UserID: credential.UserID.Hex(),
	}
	accessToken, err := utils.GenerateAccessToken(tokenPayload)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to generate access token", err.Error())
	}
	refreshToken, err := utils.GenerateRefreshToken(tokenPayload)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to generate refresh token", err.Error())
	}
	return response.Ok(c, "Token Refreshed successfully", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
