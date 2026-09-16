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
	locationDB         = "golang_practice"
	locationCollection = "locations"
)

// ErrLocationNotFound is returned when a location document does not exist.
var ErrLocationNotFound = errors.New("location not found")

func locationColl() *mongo.Collection {
	return db.Database(locationDB).Collection(locationCollection)
}

// CreateLocation inserts a new location document.
func CreateLocation(ctx context.Context, location models.Location) (models.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	location.ID = bson.NewObjectID()
	if _, err := locationColl().InsertOne(ctx, location); err != nil {
		return models.Location{}, err
	}
	return location, nil
}

// GetLocations returns all location documents, optionally filtered by city.
func GetLocations(ctx context.Context, city string) ([]models.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if city != "" {
		filter["city"] = city
	}

	cursor, err := locationColl().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	locations := make([]models.Location, 0)
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, err
	}
	return locations, nil
}

// GetLocationByID returns a single location by its hex id.
func GetLocationByID(ctx context.Context, id string) (models.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Location{}, ErrLocationNotFound
	}

	var location models.Location
	err = locationColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&location)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Location{}, ErrLocationNotFound
		}
		return models.Location{}, err
	}
	return location, nil
}

// UpdateLocation replaces the mutable fields of an existing location.
func UpdateLocation(ctx context.Context, id string, location models.Location) (models.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Location{}, ErrLocationNotFound
	}

	locationDetails := bson.M{
		"name":     location.Name,
		"address":  location.Address,
		"city":     location.City,
		"state":    location.State,
		"country":  location.Country,
		"zip_code": location.ZipCode,
	}

	coordinates := bson.M{
		"latitude":  location.Latitude,
		"longitude": location.Longitude,
	}

	updateFields := bson.M{}
	for key, value := range locationDetails {
		updateFields[key] = value
	}
	for key, value := range coordinates {
		updateFields[key] = value
	}

	update := bson.M{"$set": updateFields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Location
	err = locationColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Location{}, ErrLocationNotFound
		}
		return models.Location{}, err
	}
	return updated, nil
}

// DeleteLocation removes a location by its hex id.
func DeleteLocation(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrLocationNotFound
	}

	res, err := locationColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrLocationNotFound
	}
	return nil
}
