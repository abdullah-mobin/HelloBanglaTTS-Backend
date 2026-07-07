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

type ApiKeyRepository struct {
	apiKeyCollection *mongo.Collection
}

func NewApiKeyRepository() *ApiKeyRepository {
	repo := &ApiKeyRepository{
		apiKeyCollection: database.DbCollections.ApiKeyCollection,
	}

	// Ensure unique index on the `key` field
	ctx := context.Background()
	idxModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "key", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	// create index (ignore error if index already exists)
	_, _ = repo.apiKeyCollection.Indexes().CreateOne(ctx, idxModel)

	return repo
}
func (r *ApiKeyRepository) GetApiKeyByName(ctx context.Context, name string) (*models.ApiKey, error) {
	var apiKeyObj models.ApiKey
	err := r.apiKeyCollection.FindOne(ctx, bson.M{"name": name}).Decode(&apiKeyObj)
	if err != nil {
		return nil, err
	}
	return &apiKeyObj, nil
}
func (r *ApiKeyRepository) GetApiKeyByKey(ctx context.Context, key string) (*models.ApiKey, error) {
	var apiKeyObj models.ApiKey
	err := r.apiKeyCollection.FindOne(ctx, bson.M{"key": key}).Decode(&apiKeyObj)
	if err != nil {
		return nil, err
	}
	return &apiKeyObj, nil
}
func (r *ApiKeyRepository) IsApiKeyAvailable(ctx context.Context, key string) (bool, error) {
	var apiKeyObj models.ApiKey
	err := r.apiKeyCollection.FindOne(ctx, bson.M{"key": key}).Decode(&apiKeyObj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true, nil
		}
		return false, err
	}
	if apiKeyObj.ID != primitive.NilObjectID {
		return false, nil
	}
	return true, nil
}

func (r *ApiKeyRepository) CreateApiKey(ctx context.Context, createApiKey *models.ApiKey) (primitive.ObjectID, error) {
	result, err := r.apiKeyCollection.InsertOne(ctx, createApiKey)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		r.apiKeyCollection.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
		return primitive.NilObjectID, errors.New("InsertedID is not an ObjectID")
	}

	return insertedID, nil
}

type PaginatedApiKeys struct {
	ApiKeys         []models.ApiKey
	HasPreviousPage bool
	HasNextPage     bool
	PreviousPage    primitive.ObjectID
	NextPage        primitive.ObjectID
}

func (r *ApiKeyRepository) FindApiKeys(ctx context.Context, filter bson.M, beforeID, afterID primitive.ObjectID, limit int) (*PaginatedApiKeys, error) {
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

	cursor, err := r.apiKeyCollection.Find(ctx, queryFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var apiKeys []models.ApiKey
	if err := cursor.All(ctx, &apiKeys); err != nil {
		return nil, err
	}

	// Determine if there are more pages
	hasMore := len(apiKeys) > limit
	if hasMore {
		apiKeys = apiKeys[:limit]
	}

	// For backward pagination, reverse to maintain descending order
	if !beforeID.IsZero() && len(apiKeys) > 0 {
		for i, j := 0, len(apiKeys)-1; i < j; i, j = i+1, j-1 {
			apiKeys[i], apiKeys[j] = apiKeys[j], apiKeys[i]
		}
	}

	result := &PaginatedApiKeys{
		ApiKeys: apiKeys,
	}

	if len(apiKeys) > 0 {
		if !beforeID.IsZero() {
			// Navigating backward
			result.HasPreviousPage = hasMore
			result.HasNextPage = true

			if hasMore {
				result.PreviousPage = apiKeys[0].ID
			}
			result.NextPage = apiKeys[len(apiKeys)-1].ID
		} else if !afterID.IsZero() {
			// Navigating forward
			result.HasPreviousPage = true
			result.HasNextPage = hasMore

			result.PreviousPage = apiKeys[0].ID
			if hasMore {
				result.NextPage = apiKeys[len(apiKeys)-1].ID
			}
		} else {
			// First page
			result.HasPreviousPage = false
			result.HasNextPage = hasMore

			if hasMore {
				result.NextPage = apiKeys[len(apiKeys)-1].ID
			}
		}
	}

	return result, nil
}

// UpdateApiKeyByKey updates an API key document identified by its key
func (r *ApiKeyRepository) UpdateApiKeyByKey(ctx context.Context, key string, update bson.M) error {
	if key == "" {
		return errors.New("API key is required")
	}

	// Ensure update is wrapped in $set when appropriate
	// If caller already provided an operator (e.g. $set), use it as-is
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

	_, err := r.apiKeyCollection.UpdateOne(ctx, bson.M{"key": key}, updateDoc)
	return err
}
