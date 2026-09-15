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
	// stateDB is the name of the MongoDB database that stores states.
	stateDB = "golang_practice"
	// stateCollection is the MongoDB collection that stores state documents.
	stateCollection = "states"
)

// ErrStateNotFound is returned when a state document does not exist.
var ErrStateNotFound = errors.New("state not found")

// stateColl returns a handle to the states collection in the configured database.
func stateColl() *mongo.Collection {
	return db.Database(stateDB).Collection(stateCollection)
}

// CreateState inserts a new state document.
func CreateState(ctx context.Context, state models.State) (models.State, error) {
	// Limit the database operation to 10 seconds and release the context resources when done.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Generate a unique MongoDB object ID for the new state.
	state.ID = bson.NewObjectID()
	// Insert the complete state document into the states collection.
	if _, err := stateColl().InsertOne(ctx, state); err != nil {
		// Return an empty state and the database error if insertion fails.
		return models.State{}, err
	}
	// Return the state, including the newly generated ID.
	return state, nil
}

// GetStates returns all state documents, optionally filtered by country.
func GetStates(ctx context.Context, country string) ([]models.State, error) {
	// Limit the database operation to 10 seconds and release the context resources when done.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Start with an empty filter so the query returns every state by default.
	filter := bson.M{}
	if country != "" {
		// Add a country condition when the caller provides a country value.
		filter["country"] = country
	}

	// Query MongoDB using the optional filter and receive a cursor for the results.
	cursor, err := stateColl().Find(ctx, filter)
	if err != nil {
		// Return the database error if the query cannot be created or executed.
		return nil, err
	}
	// Close the cursor after all query results have been processed.
	defer cursor.Close(ctx)

	// Create an empty slice so a successful query returns [] instead of nil when there are no states.
	states := make([]models.State, 0)
	if err := cursor.All(ctx, &states); err != nil {
		// Return the error if MongoDB cannot decode all cursor results.
		return nil, err
	}
	// Return all states found by the query.
	return states, nil
}

// GetStateByID returns a single state by its hex id.
func GetStateByID(ctx context.Context, id string) (models.State, error) {
	// Limit the database operation to 10 seconds and release the context resources when done.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Convert the text ID from the request into MongoDB's ObjectID type.
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// Treat an invalid ID format as a missing state.
		return models.State{}, ErrStateNotFound
	}

	// Declare a variable that MongoDB will populate with the matching state.
	var state models.State
	// Find the document by its _id field and decode it into the state variable.
	err = stateColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&state)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Convert MongoDB's no-document error into the service's standard error.
			return models.State{}, ErrStateNotFound
		}
		// Return any other database error unchanged.
		return models.State{}, err
	}
	// Return the state that was found.
	return state, nil
}

// UpdateState replaces the mutable fields of an existing state.
func UpdateState(ctx context.Context, id string, state models.State) (models.State, error) {
	// Limit the database operation to 10 seconds and release the context resources when done.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Convert the text ID from the request into MongoDB's ObjectID type.
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// Treat an invalid ID format as a missing state.
		return models.State{}, ErrStateNotFound
	}

	// Build an update that changes only the state's name, code, and country.
	update := bson.M{"$set": bson.M{
		"name":    state.Name,
		"code":    state.Code,
		"country": state.Country,
	}}

	// Configure the update to return the document after the changes are applied.
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	// Declare a variable that will receive the updated state document.
	var updated models.State
	// Find the state by ID, apply the update, and decode the resulting document.
	err = stateColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Convert MongoDB's no-document error into the service's standard error.
			return models.State{}, ErrStateNotFound
		}
		// Return any other database error unchanged.
		return models.State{}, err
	}
	// Return the updated state document.
	return updated, nil
}

// DeleteState removes a state by its hex id.
func DeleteState(ctx context.Context, id string) error {
	// Limit the database operation to 10 seconds and release the context resources when done.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Convert the text ID from the request into MongoDB's ObjectID type.
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		// Treat an invalid ID format as a missing state.
		return ErrStateNotFound
	}

	// Delete the state document that matches the requested ID.
	res, err := stateColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		// Return the database error if deletion fails.
		return err
	}
	if res.DeletedCount == 0 {
		// Report that the state was not found when no document was deleted.
		return ErrStateNotFound
	}
	// A nil error indicates that the state was deleted successfully.
	return nil
}
