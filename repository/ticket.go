package repository

import (
	"context"
	"errors"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/database"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TicketRepository struct {
	ticketCollection *mongo.Collection
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{
		ticketCollection: database.DbCollections.TicketCollection,
	}
}

func (r *TicketRepository) CreateTicket(ctx context.Context, createTicket *models.Ticket) (primitive.ObjectID, error) {

	dbRes, err := r.ticketCollection.InsertOne(ctx, createTicket)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := dbRes.InsertedID.(primitive.ObjectID)
	if !ok {
		r.ticketCollection.DeleteOne(ctx, bson.M{"_id": dbRes.InsertedID})
		return primitive.NilObjectID, errors.New("InsertedID is not an ObjectID")

	}

	return insertedID, nil
}

type PaginatedTicket struct {
	Tickets         []models.Ticket
	HasPreviousPage bool
	HasNextPage     bool
	PreviousPage    primitive.ObjectID
	NextPage        primitive.ObjectID
}

func (r *TicketRepository) FindTickets(ctx context.Context, filter bson.M, beforeID, afterID primitive.ObjectID, limit int) (*PaginatedTicket, error) {
	var sortOrder int = -1

	queryFilter := bson.M{}
	for k, v := range filter {
		queryFilter[k] = v
	}

	if !beforeID.IsZero() {
		queryFilter["_id"] = bson.M{"$lt": beforeID}
		sortOrder = -1
	} else if !afterID.IsZero() {
		queryFilter["_id"] = bson.M{"$gt": afterID}
		sortOrder = 1
	}

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: sortOrder}}).SetLimit(int64(limit + 1))

	cursor, err := r.ticketCollection.Find(ctx, queryFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []models.Ticket
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}

	hasMore := len(tickets) > limit
	if hasMore {
		tickets = tickets[:limit]
	}

	if !beforeID.IsZero() && len(tickets) > 0 {
		for i, j := 0, len(tickets)-1; i < j; i, j = i+1, j-1 {
			tickets[i], tickets[j] = tickets[j], tickets[i]
		}
	}

	result := &PaginatedTicket{Tickets: tickets}

	if len(tickets) > 0 {
		if !beforeID.IsZero() {
			result.HasPreviousPage = hasMore
			result.HasNextPage = true
			if hasMore {
				result.PreviousPage = tickets[0].ID
			}
			result.NextPage = tickets[len(tickets)-1].ID
		} else if !afterID.IsZero() {
			result.HasPreviousPage = true
			result.HasNextPage = hasMore
			result.PreviousPage = tickets[0].ID
			if hasMore {
				result.NextPage = tickets[len(tickets)-1].ID
			}
		} else {
			result.HasPreviousPage = false
			result.HasNextPage = hasMore
			if hasMore {
				result.NextPage = tickets[len(tickets)-1].ID
			}
		}
	}

	return result, nil
}

func (r *TicketRepository) UpdateTicketByID(ctx context.Context, id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// wrap in $set if no operator provided
	var updateDoc bson.M
	for k := range update {
		if len(k) > 0 && k[0] == '$' {
			updateDoc = update
			break
		}
	}
	if updateDoc == nil {
		updateDoc = bson.M{"$set": update}
	}

	_, err = r.ticketCollection.UpdateOne(ctx, bson.M{"_id": objID}, updateDoc)
	return err
}

func (r *TicketRepository) GetTicketByID(ctx context.Context, strID string) (*models.Ticket, error) {
	objID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		return nil, err
	}
	var ticket models.Ticket
	err = r.ticketCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) DeleteTicketByID(ctx context.Context, strID string) error {
	objID, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		return err
	}

	_, err = r.ticketCollection.DeleteOne(ctx, bson.M{"_id": objID})

	return err
}
