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
	stageDB         = "golang_practice"
	stageCollection = "stages"
)

// ErrStageNotFound is returned when a stage document does not exist.
var ErrStageNotFound = errors.New("stage not found")

func stageColl() *mongo.Collection {
	return db.Database(stageDB).Collection(stageCollection)
}

// CreateStage inserts a new stage document.
func CreateStage(ctx context.Context, stage models.Stage) (models.Stage, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stage.ID = bson.NewObjectID()
	if _, err := stageColl().InsertOne(ctx, stage); err != nil {
		return models.Stage{}, err
	}
	return stage, nil
}

// GetStages returns all stage documents.
func GetStages(ctx context.Context) ([]models.Stage, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := stageColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stages := make([]models.Stage, 0)
	if err := cursor.All(ctx, &stages); err != nil {
		return nil, err
	}
	return stages, nil
}

// GetStageByID returns a single stage by its hex id.
func GetStageByID(ctx context.Context, id string) (models.Stage, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Stage{}, ErrStageNotFound
	}

	var stage models.Stage
	err = stageColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&stage)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Stage{}, ErrStageNotFound
		}
		return models.Stage{}, err
	}
	return stage, nil
}

// UpdateStage replaces the mutable fields of an existing stage.
func UpdateStage(ctx context.Context, id string, stage models.Stage) (models.Stage, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Stage{}, ErrStageNotFound
	}

	update := bson.M{"$set": bson.M{
		"name":        stage.Name,
		"description": stage.Description,
		"sequence":    stage.Sequence,
		"status":      stage.Status,
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Stage
	err = stageColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Stage{}, ErrStageNotFound
		}
		return models.Stage{}, err
	}
	return updated, nil
}

// DeleteStage removes a stage by its hex id.
func DeleteStage(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrStageNotFound
	}

	res, err := stageColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrStageNotFound
	}
	return nil
}
