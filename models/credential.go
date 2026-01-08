package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credential struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"password"`
	RefreshToken string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

func (obj Credential) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Email, validation.Required),
		validation.Field(&obj.Password, validation.Required),
		validation.Field(&obj.UserID, validation.Required),
	)
}

type CredentialSearch struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty"`
	UserId   primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

func (obj CredentialSearch) GetCredentialSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.ID != primitive.NilObjectID {
		self["_id"] = obj.ID
	}

	if obj.Username != "" {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Username)
		self["username"] = bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.UserId != primitive.NilObjectID {
		self["user_id"] = obj.UserId
	}

	return self
}
