package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"practice/api/db"
	"practice/api/models"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const vehicleLocationCacheTTL = 24 * time.Hour

var ErrVehicleLocationNotFound = errors.New("vehicle location not found")

func vehicleLocationColl() *mongo.Collection {
	return db.Database("golang_practice").Collection("vehicle_locations")
}

func vehicleLocationCacheKey(vehicleID string) string {
	return fmt.Sprintf("vehicle:location:%s", vehicleID)
}

func SaveVehicleLocation(ctx context.Context, vehicleID string, location models.VehicleLocation) (models.VehicleLocation, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	location.VehicleID = vehicleID
	location.UpdatedAt = time.Now().UTC()
	_, err := vehicleLocationColl().UpdateOne(ctx, bson.M{"vehicle_id": vehicleID}, bson.M{"$set": location}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return models.VehicleLocation{}, err
	}
	if err := cacheVehicleLocation(ctx, location); err != nil {
		return models.VehicleLocation{}, err
	}
	return location, nil
}

func GetVehicleLocation(ctx context.Context, vehicleID string) (models.VehicleLocation, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cached, err := db.RedisClient.Get(ctx, vehicleLocationCacheKey(vehicleID)).Result()
	if err == nil {
		var location models.VehicleLocation
		if json.Unmarshal([]byte(cached), &location) == nil {
			return location, true, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		return models.VehicleLocation{}, false, err
	}

	var location models.VehicleLocation
	if err := vehicleLocationColl().FindOne(ctx, bson.M{"vehicle_id": vehicleID}).Decode(&location); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.VehicleLocation{}, false, ErrVehicleLocationNotFound
		}
		return models.VehicleLocation{}, false, err
	}
	if err := cacheVehicleLocation(ctx, location); err != nil {
		return models.VehicleLocation{}, false, err
	}
	return location, false, nil
}

func cacheVehicleLocation(ctx context.Context, location models.VehicleLocation) error {
	payload, err := json.Marshal(location)
	if err != nil {
		return err
	}
	return db.RedisClient.Set(ctx, vehicleLocationCacheKey(location.VehicleID), payload, vehicleLocationCacheTTL).Err()
}
