package repository

import (
	"context"
	"errors"

	"github.com/abdullah-mobin/helloBanglaTTS-backend/database"
	"github.com/abdullah-mobin/helloBanglaTTS-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CredentialRepository struct {
	credentialCollection *mongo.Collection
}

func NewCredentialRepository() *CredentialRepository {
	return &CredentialRepository{
		credentialCollection: database.DbCollections.CredentialCollection,
	}
}
func (r *CredentialRepository) GetCredentialByEmail(ctx context.Context, email string) (*models.Credential, error) {
	var credential models.Credential
	err := r.credentialCollection.FindOne(ctx, bson.M{"email": email}).Decode(&credential)
	if err != nil {
		return nil, err
	}
	return &credential, nil
}
func (r *CredentialRepository) IsEmailAvailable(ctx context.Context, email string) (bool, error) {
	var user models.User
	err := r.credentialCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
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

func (r *CredentialRepository) CreateCredential(ctx context.Context, createCredential *models.Credential) (primitive.ObjectID, error) {

	dbCredential, err := r.credentialCollection.InsertOne(ctx, createCredential)
	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := dbCredential.InsertedID.(primitive.ObjectID)
	if !ok {
		r.credentialCollection.DeleteOne(ctx, bson.M{"_id": dbCredential.InsertedID})
		return primitive.NilObjectID, errors.New("InsertedID is not an ObjectID")
	}

	return insertedID, nil
}

func (r *CredentialRepository) UpdateCredential(ctx context.Context, id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.credentialCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	return err
}

func (r *CredentialRepository) InvalidateRefreshTokens(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.credentialCollection.UpdateMany(ctx, bson.M{"user_id": objID}, bson.M{"$unset": bson.M{"refresh_token": ""}})
	return err
}

func (r *CredentialRepository) FindCredentialUsingRefreshToken(ctx context.Context, refreshToken string) (error, *models.Credential) {
	var credential models.Credential
	err := r.credentialCollection.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&credential)
	if err != nil {
		return err, &credential
	}
	return nil, &credential
}
