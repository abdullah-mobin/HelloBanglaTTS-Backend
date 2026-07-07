package handlers

import (
	"context"
	"log"
	"strings"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/dtos"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/models"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/repository"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/utils"
	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApiKeys godoc
//
//	@Summary		Get all API keys
//	@Description	get all API keys by some filters{id,name,limit,beforeId,afterId}
//	@Tags			ApiKey
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			query	string	false	"Filter by id"
//	@Param			key			query	string	false	"Filter by key"
//	@Param			beforeId	query	string	false	"Filter by beforeId"
//	@Param			afterId		query	string	false	"Filter by afterId"
//	@Param			limit		query	string	false	"Number of items per page (default 10)"
//	@Router			/api-keys/ [get]
func FindApiKeys(c *fiber.Ctx) error {
	var err error
	queries := c.Queries()
	eligibleFilters := []string{"id", "key", "limit", "beforeId", "afterId"}
	objectIDFields := []string{"id"}
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
	repo := repository.NewApiKeyRepository()
	paginatedApiKeys, err := repo.FindApiKeys(context.Background(), filters, beforeID, afterID, limit)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to retrieve API keys", err.Error())
	}
	pagination := response.Pagination{
		HasNextPage:     paginatedApiKeys.HasNextPage,
		HasPreviousPage: paginatedApiKeys.HasPreviousPage,
		NextPage:        paginatedApiKeys.NextPage.Hex(),
		PreviousPage:    paginatedApiKeys.PreviousPage.Hex(),
	}
	return response.PaginatedResponse(c, "API Keys Retrieved Successfully", paginatedApiKeys.ApiKeys, pagination)
}

// CreateApiKey godoc
//
//	@Summary		Create a new API key
//	@Description	Create a new API key for a user
//	@Tags			ApiKey
//	@Accept			json
//	@Produce		json
//	@Param			apiKey	body	dtos.CreateApiKeyDTO	true	"API Key payload"
//	@Security		BearerAuth
//	@Router			/api-keys/new [post]
func CreateApiKey(c *fiber.Ctx) error {
	var payload dtos.CreateApiKeyDTO
	if err := c.BodyParser(&payload); err != nil {
		return response.BadRequestException(c, "Invalid request body", err.Error())
	}

	if err := payload.Validate(); err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid request", errorsArr)
	}

	// Create repository early so we can check key uniqueness
	repo := repository.NewApiKeyRepository()

	// Generate a unique secure random key (32 bytes -> 64 hex chars)
	var key string
	const maxAttempts = 10
	for i := 0; i < maxAttempts; i++ {
		k, err := utils.GenerateRandomKey(32)
		if err != nil {
			return response.InternalServerErrorException(c, "Failed to generate API key", err.Error())
		}
		available, err := repo.IsApiKeyAvailable(context.Background(), k)
		if err != nil {
			return response.InternalServerErrorException(c, "Failed to check API key uniqueness", err.Error())
		}
		if available {
			key = k
			break
		}
	}
	if key == "" {
		return response.InternalServerErrorException(c, "Unable to generate a unique API key", "")
	}

	apiKey := models.ApiKey{
		Name:      payload.Name,
		Key:       key,
		IsActive:  true,
		Count:     0,
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	}

	if err := apiKey.Validate(); err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid API key data", errorsArr)
	}

	if _, err := repo.CreateApiKey(context.Background(), &apiKey); err != nil {
		return response.InternalServerErrorException(c, "Failed to create API key", err.Error())
	}

	return response.Created(c, "API key created successfully", map[string]string{"key": key})
}

// UpdateApiKey godoc
//
//	@Summary		Update API key
//	@Description	Update API key for a user
//	@Tags			ApiKey
//	@Accept			json
//	@Produce		json
//	@Param			key		path	string					true	"API key"
//	@Param			apiKey	body	dtos.UpdateApiKeyDTO	true	"API Key payload"
//	@Security		BearerAuth
//	@Router			/api-keys/{key} [put]
func UpdateApiKey(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return response.BadRequestException(c, "API key is required", "")
	}

	var payload dtos.UpdateApiKeyDTO
	if err := c.BodyParser(&payload); err != nil {
		return response.BadRequestException(c, "Invalid request body", err.Error())
	}

	repo := repository.NewApiKeyRepository()
	if _, err := repo.GetApiKeyByKey(context.Background(), key); err != nil {
		if err == mongo.ErrNoDocuments {
			return response.NotFoundException(c, "API key not found", "")
		}
		return response.InternalServerErrorException(c, "Failed to retrieve API key", err.Error())
	}

	update := bson.M{}

	if payload.Name != "" {
		update["name"] = payload.Name
	}

	if payload.IsActive != nil {
		update["is_active"] = payload.IsActive
	}

	update["updated_at"] = utils.Now()

	log.Println(key)

	if err := repo.UpdateApiKeyByKey(context.Background(), key, update); err != nil {
		return response.InternalServerErrorException(c, "Failed to create API key", err.Error())
	}

	return response.Ok(c, "API key updated successfully", nil)
}
