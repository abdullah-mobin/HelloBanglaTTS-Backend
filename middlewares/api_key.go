package middlewares

// import (
// 	"context"
// 	"os"

// 	"github.com/abdullah-mobin/helloBanglaTTS-backend/repository"
// 	"github.com/abdullah-mobin/helloBanglaTTS-backend/response"
// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func ValidateApiKey(c *fiber.Ctx) error {

// 	apiKeyStr := c.Get("SenderBee-API-KEY")
// 	if apiKeyStr == "" {
// 		return c.Next()
// 	}

// 	apiRepo := repository.NewApiKeyRepository()
// 	apiKey, err := apiRepo.GetApiKeyByKey(context.Background(), apiKeyStr)
// 	if err == mongo.ErrNoDocuments || apiKey == nil || !apiKey.IsActive {
// 		return response.UnauthorizedException(c, "failed to validate api key", err.Error())
// 	}

// 	rmqPublisher, err := rabbitmq.NewPublisher()
// 	if err != nil {
// 		return response.InternalServerErrorException(c, "failed to publish email", err.Error())
// 	}
// 	defer rmqPublisher.Close()

// 	var emailData dtos.SendEmailDTO
// 	if err := c.BodyParser(&emailData); err != nil {
// 		return response.BadRequestException(c, "Invalid request body", err.Error())
// 	}
// 	if err := emailData.Validate(); err != nil {
// 		return response.BadRequestException(c, "Validation failed", err.Error())
// 	}

// 	data := map[string]any{
// 		"payload": emailData,
// 		"api_key": apiKeyStr,
// 		"count":   apiKey.Count + 1,
// 	}

// 	rmqPublisher.PublishJSON(os.Getenv("API_BASED_EMAIL_QUEUE"), data)

// 	return response.Ok(c, "email queued successfully", nil)
// }
