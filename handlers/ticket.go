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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ticket godoc
//
//	@Summary		Get tickets
//	@Description	Retrieve tickets filtered by optional filters
//	@Tags			Ticket
//	@Accept			json
//	@Produce		json
//	@Param			serial		query	string	false	"Serial"
//	@Param			limit		query	int		false	"limit"
//	@Param			beforeId	query	string	false	"beforeId"
//	@Param			afterId		query	string	false	"afterId"
//	@Security		BearerAuth
//	@Router			/support/ticket [get]
func FindTickets(c *fiber.Ctx) error {
	queries := c.Queries()
	eligibleFilters := []string{"id", "serial", "status", "priority", "beforeId", "afterId"}
	objectIDFields := []string{"id"}
	filters, err := utils.ParseFilters(queries, eligibleFilters, objectIDFields)
	if err != nil {
		return response.BadRequestException(c, "Invalid filter parameters", []string{err.Error()})
	}

	limit := utils.ParseLimit(queries, 10)
	if limit < 1 || limit > 100 {
		return response.BadRequestException(c, "Limit must be between 1 and 100", "")
	}

	beforeIDStr, afterIDStr := utils.ParsePaginationIDs(queries)
	var beforeID, afterID primitive.ObjectID
	if afterIDStr != primitive.NilObjectID {
		afterID = afterIDStr
	}
	if beforeIDStr != primitive.NilObjectID {
		beforeID = beforeIDStr
	}

	if !afterID.IsZero() && !beforeID.IsZero() {
		return response.BadRequestException(c, "Cannot specify both afterId and beforeId", "")
	}

	repo := repository.NewTicketRepository()
	result, err := repo.FindTickets(context.Background(), filters, beforeID, afterID, limit)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to retrieve tickets", err.Error())
	}

	pagination := response.Pagination{
		HasNextPage:     result.HasNextPage,
		HasPreviousPage: result.HasPreviousPage,
		NextPage:        result.NextPage.Hex(),
		PreviousPage:    result.PreviousPage.Hex(),
	}

	return response.PaginatedResponse(c, "Tickets Retrieved Successfully", result.Tickets, pagination)
}

// Ticket godoc
//
//	@Summary		Create Tickets
//	@Description	Create a ticket for a support
//	@Tags			Ticket
//	@Accept			json
//	@Produce		json
//	@Param			ticket	body	dtos.CreateTicketDTO	true	"ticket payload"
//	@Security		BearerAuth
//	@Router			/support/ticket/new [post]
func CreateTicket(c *fiber.Ctx) error {
	var payload dtos.CreateTicketDTO
	if err := c.BodyParser(&payload); err != nil {
		return response.BadRequestException(c, "Invalid request body", err.Error())
	}
	if err := payload.Validate(); err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid request", errorsArr)
	}

	ticket := &models.Ticket{
		Serial:    utils.GenerateSerial("TKT"),
		Subject:   payload.Subject,
		Category:  payload.Category,
		Status:    models.TicketStatusOpen,
		Priority:  payload.Priority,
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	}

	repo := repository.NewTicketRepository()
	insertedID, err := repo.CreateTicket(context.Background(), ticket)
	if err != nil {
		return response.InternalServerErrorException(c, "Failed to create ticket", err.Error())
	}

	return response.Created(c, "ticket created successfully", insertedID)
}

// Ticket godoc
//
//	@Summary		Update ticket
//	@Description	Update ticket fields by ID
//	@Tags			Ticket
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"Ticket ID"
//	@Param			ticket	body	dtos.UpdateTicketDTO	true	"ticket payload"
//	@Security		BearerAuth
//	@Router			/support/ticket/{id} [put]
func UpdateTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequestException(c, "Ticket id is required", "")
	}

	var payload dtos.UpdateTicketDTO
	if err := c.BodyParser(&payload); err != nil {
		return response.BadRequestException(c, "Invalid request body", err.Error())
	}

	if err := payload.Validate(); err != nil {
		errorsArr := strings.Split(err.Error(), ";")
		return response.ValidationException(c, "Invalid request", errorsArr)
	}

	repo := repository.NewTicketRepository()
	if _, err := repo.GetTicketByID(context.Background(), id); err != nil {
		return response.NotFoundException(c, "Ticket not found", err.Error())
	}

	update := bson.M{}
	if payload.Subject != "" {
		update["subject"] = payload.Subject
	}
	if payload.Category != "" {
		update["category"] = payload.Category
	}
	if payload.Priority != "" {
		update["priority"] = payload.Priority
	}
	if payload.Status != "" {
		update["status"] = payload.Status
	}

	update["updated_at"] = utils.Now()

	if len(update) == 1 { // only updated_at set
		return response.BadRequestException(c, "No fields to update", "")
	}

	if err := repo.UpdateTicketByID(context.Background(), id, update); err != nil {
		return response.InternalServerErrorException(c, "Failed to update ticket", err.Error())
	}

	return response.Ok(c, "Ticket updated successfully", nil)
}

// Ticket godoc
//
//	@Summary		Delete ticket
//	@Description	Delete ticket by ID
//	@Tags			Ticket
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Ticket ID"
//	@Security		BearerAuth
//	@Router			/support/ticket/{id} [delete]
func DeleteTicket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequestException(c, "Ticket id required", nil)
	}

	repo := repository.NewTicketRepository()
	if _, err := repo.GetTicketByID(context.Background(), id); err != nil {
		return response.NotFoundException(c, "Ticket not found", err.Error())
	}

	if err := repo.DeleteTicketByID(context.Background(), id); err != nil {
		return response.InternalServerErrorException(c, "Failed To Delete Ticket", err.Error())
	}

	return response.Ok(c, "Ticket Deleted Successfully", nil)
}
