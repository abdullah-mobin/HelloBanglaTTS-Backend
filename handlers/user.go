package handlers

import (
	"context"
	"strings"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/dtos"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/models"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/repository"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/utils"
	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Profile godoc
//
//	@Summary		Get Profile by id
//	@Description	get a hello-bangla-tts users profile by id
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Router			/user/profile [get]
func UserProfile(c *fiber.Ctx) error {
	id := c.Locals("userId")
	repo := repository.NewUserRepository()
	user, err := repo.GetUserByID(context.Background(), id.(string))
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to retrieve user", err.Error())
	}
	return response.Ok(c, "User Profile Retrieved Successfully", user)
}

// Users godoc
//
//	@Summary		Get all hello-bangla-tts users
//	@Description	get all hello-bangla-tts users profile by some filters{id,name,limit,beforeId,afterId}
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			query	string	false	"Filter by id"
//	@Param			name		query	string	false	"Filter by name"
//	@Param			beforeId	query	string	false	"Filter by beforeId"
//	@Param			afterId		query	string	false	"Filter by afterId"
//	@Param			limit		query	string	false	"Number of items per page (default 10)"
//	@Router			/user/ [get]
func Users(c *fiber.Ctx) error {
	var err error
	queries := c.Queries()
	eligibleFilters := []string{"id", "name", "limit", "beforeId", "afterId"}
	objectIDFields := []string{}
	filters, err := utils.ParseFilters(queries, eligibleFilters, objectIDFields)
	if err != nil {
		return response.BadRequestException(c, "Invalid filter parameters", []string{err.Error()})
	}
	limit := utils.ParseLimit(queries, 10)
	// Validate limit
	if limit < 1 || limit > 100 {
		return response.BadRequestException(c, "Limit must be between 1 and 100", "")
	}
	beforeIDStr, afterIDStr := utils.ParsePaginationIDs(queries)
	var afterID, beforeID primitive.ObjectID
	if afterIDStr != primitive.NilObjectID {
		afterID = afterIDStr

	}
	if beforeIDStr != primitive.NilObjectID {
		beforeID = beforeIDStr

	}
	// Validate: cannot have both
	if !afterID.IsZero() && !beforeID.IsZero() {
		return response.BadRequestException(c, "Cannot specify both afterId and beforeId", "")
	}
	repo := repository.NewUserRepository()
	paginatedUsers, err := repo.FindUsers(context.Background(), filters, beforeID, afterID, limit)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to retrieve users", err.Error())
	}
	pagination := response.Pagination{
		HasNextPage:     paginatedUsers.HasNextPage,
		HasPreviousPage: paginatedUsers.HasPreviousPage,
		NextPage:        paginatedUsers.NextPage.Hex(),
		PreviousPage:    paginatedUsers.PreviousPage.Hex(),
	}
	return response.PaginatedResponse(c, "Users Retrieved Successfully", paginatedUsers.Users, pagination)
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	create a new user account in hello-bangla-tts
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			register	body	dtos.RegisterDTO	true	"user register dto"
//	@Router			/auth/register [post]
func RegisterNewUser(c *fiber.Ctx) error {
	var self dtos.RegisterDTO
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid request", errorsArr)
	}
	user := models.User{
		Name:      self.Name,
		Email:     self.Email,
		Status:    models.StatusActive,
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	}
	err = user.Validate()
	if err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid request", errorsArr)
	}

	credentialRepo := repository.NewCredentialRepository()
	isAvailable, err := credentialRepo.IsEmailAvailable(context.Background(), self.Email)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to check email availability", err.Error())
	}
	if !isAvailable {
		return response.ConflictException(c, "Email already in use", "Please use a different email")
	}
	repo := repository.NewUserRepository()
	userID, err := repo.CreateUser(context.Background(), &user)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to create user", err.Error())
	}

	credential := models.Credential{
		Password:  utils.HashPassword(self.Password),
		Email:     self.Email,
		UserID:    userID,
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	}

	err = credential.Validate()
	if err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid request", errorsArr)
	}
	_, err = credentialRepo.CreateCredential(context.Background(), &credential)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to create user credential", err.Error())
	}

	return response.Created(c, "User Registered Successfully", userID)
}
