package utils

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseFilters(queries map[string]string, eligibleFilters []string, objectIDFields []string) (bson.M, error) {
	filters := bson.M{}

	// Convert objectIDFields slice to map for faster lookup
	objectIDFieldsMap := make(map[string]bool)
	for _, field := range objectIDFields {
		objectIDFieldsMap[field] = true
	}

	for _, key := range eligibleFilters {
		val, ok := queries[key]
		if !ok || val == "" {
			continue
		}

		// Skip pagination parameters - they're handled separately
		if key == "beforeId" || key == "afterId" || key == "limit" || key == "id" {
			continue
		}

		// Convert to ObjectID if specified
		if objectIDFieldsMap[key] {
			objectID, err := primitive.ObjectIDFromHex(val)
			if err != nil {
				return nil, fmt.Errorf("invalid ObjectID for %s: %v", key, err)
			}
			filters[key] = objectID
		} else {
			// Regular string filter
			filters[key] = bson.M{
				"$regex":   val,
				"$options": "i", // case-insensitive
			}
		}
	}
	return filters, nil
}

func ParseLimit(queries map[string]string, defaultLimit int) int {
	limitStr, ok := queries["limit"]
	if !ok || limitStr == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return defaultLimit
	}

	return limit
}

func ParsePaginationIDs(queries map[string]string) (primitive.ObjectID, primitive.ObjectID) {
	beforeID := queries["beforeId"]
	afterID := queries["afterId"]

	var beforeObjID, afterObjID primitive.ObjectID
	var err error

	if beforeID != "" {
		beforeObjID, err = primitive.ObjectIDFromHex(beforeID)
		if err != nil {
			return primitive.NilObjectID, primitive.NilObjectID
		}
	}

	if afterID != "" {
		afterObjID, err = primitive.ObjectIDFromHex(afterID)
		if err != nil {
			return primitive.NilObjectID, primitive.NilObjectID
		}
	}

	return beforeObjID, afterObjID
}
