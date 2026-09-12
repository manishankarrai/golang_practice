package services

import (
	"context"
	"errors"
	"time"

	"practice/api/db"
	"practice/api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	stateDB         = "golang_practice"
	stateCollection = "states"
)

// ErrStateNotFound is returned when a state document does not exist.
var ErrStateNotFound = errors.New("state not found")

func stateColl() *mongo.Collection {
	return db.Database(stateDB).Collection(stateCollection)
}

// CreateState inserts a new state document.
func CreateState(ctx context.Context, state models.State) (models.State, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	state.ID = bson.NewObjectID()
	if _, err := stateColl().InsertOne(ctx, state); err != nil {
		return models.State{}, err
	}
	return state, nil
}

// GetStates returns all state documents, optionally filtered by country.
func GetStates(ctx context.Context, country string) ([]models.State, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if country != "" {
		filter["country"] = country
	}

	cursor, err := stateColl().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	states := make([]models.State, 0)
	if err := cursor.All(ctx, &states); err != nil {
		return nil, err
	}
	return states, nil
}

// GetStateByID returns a single state by its hex id.
func GetStateByID(ctx context.Context, id string) (models.State, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.State{}, ErrStateNotFound
	}

	var state models.State
	err = stateColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&state)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.State{}, ErrStateNotFound
		}
		return models.State{}, err
	}
	return state, nil
}

// UpdateState replaces the mutable fields of an existing state.
func UpdateState(ctx context.Context, id string, state models.State) (models.State, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.State{}, ErrStateNotFound
	}

	update := bson.M{"$set": bson.M{
		"name":    state.Name,
		"code":    state.Code,
		"country": state.Country,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.State
	err = stateColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.State{}, ErrStateNotFound
		}
		return models.State{}, err
	}
	return updated, nil
}

// DeleteState removes a state by its hex id.
func DeleteState(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrStateNotFound
	}

	res, err := stateColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrStateNotFound
	}
	return nil
}
