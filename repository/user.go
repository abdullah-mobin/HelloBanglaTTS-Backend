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

type UserRepository struct {
	userCollection       *mongo.Collection
	credentialCollection *mongo.Collection
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		userCollection: database.DbCollections.UserCollection,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, createUser *models.User) (primitive.ObjectID, error) {
	result, err := r.userCollection.InsertOne(ctx, createUser)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		r.userCollection.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
		return primitive.NilObjectID, errors.New("InsertedID is not an ObjectID")
	}

	return insertedID, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.userCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetCredentialByID(ctx context.Context, id string) (*models.Credential, error) {
	var credential models.Credential
	err := r.credentialCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&credential)
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

func (r *UserRepository) GetCredentialByUsername(ctx context.Context, username string) (*models.Credential, error) {
	var credential models.Credential
	err := r.credentialCollection.FindOne(ctx, bson.M{"username": username}).Decode(&credential)
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

func (r *UserRepository) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	var user models.User
	err := r.credentialCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true, nil
		}
		return false, err
	}
	if user.ID != primitive.NilObjectID {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, update bson.M) error {
	_, err := r.userCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.userCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

type PaginatedUsers struct {
	Users           []models.User
	HasPreviousPage bool
	HasNextPage     bool
	PreviousPage    primitive.ObjectID
	NextPage        primitive.ObjectID
}

func (r *UserRepository) FindUsers(ctx context.Context, filter bson.M, beforeID, afterID primitive.ObjectID, limit int) (*PaginatedUsers, error) {
	var sortOrder int = -1 // Default: show latest users first

	// Clone filter to avoid mutation
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

	// Fetch limit+1 to check if there are more items
	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: sortOrder}}).
		SetLimit(int64(limit + 1)).
		SetProjection(bson.M{"password": 0}) // Exclude sensitive fields if needed

	cursor, err := r.userCollection.Find(ctx, queryFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// Determine if there are more pages
	hasMore := len(users) > limit
	if hasMore {
		users = users[:limit]
	}

	// For backward pagination, reverse to maintain descending order
	if !beforeID.IsZero() && len(users) > 0 {
		for i, j := 0, len(users)-1; i < j; i, j = i+1, j-1 {
			users[i], users[j] = users[j], users[i]
		}
	}

	result := &PaginatedUsers{
		Users: users,
	}

	if len(users) > 0 {
		if !beforeID.IsZero() {
			// Navigating backward
			result.HasPreviousPage = hasMore
			result.HasNextPage = true

			if hasMore {
				result.PreviousPage = users[0].ID
			}
			result.NextPage = users[len(users)-1].ID
		} else if !afterID.IsZero() {
			// Navigating forward
			result.HasPreviousPage = true
			result.HasNextPage = hasMore

			result.PreviousPage = users[0].ID
			if hasMore {
				result.NextPage = users[len(users)-1].ID
			}
		} else {
			// First page
			result.HasPreviousPage = false
			result.HasNextPage = hasMore

			if hasMore {
				result.NextPage = users[len(users)-1].ID
			}
		}
	}

	return result, nil
}
