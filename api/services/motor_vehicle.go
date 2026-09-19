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

const vehicleTimeout = 10 * time.Second

var ErrMotorVehicleNotFound = errors.New("motor vehicle not found")

func vehicleColl() *mongo.Collection {
	return db.Database("golang_practice").Collection("motor_vehicles")
}

func vehicleObjectID(id string) (bson.ObjectID, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrMotorVehicleNotFound
	}
	return objID, nil
}

func vehicleFields(vehicle models.MotorVehicle) bson.M {
	return bson.M{"make": vehicle.Make, "model": vehicle.Model, "year": vehicle.Year, "color": vehicle.Color, "vehicle_type": vehicle.VehicleType, "vin": vehicle.VIN, "engine_no": vehicle.EngineNo}
}

func CreateMotorVehicle(ctx context.Context, vehicle models.MotorVehicle) (models.MotorVehicle, error) {
	ctx, cancel := context.WithTimeout(ctx, vehicleTimeout)
	defer cancel()
	vehicle.ID = bson.NewObjectID()
	_, err := vehicleColl().InsertOne(ctx, vehicle)
	if err != nil {
		return models.MotorVehicle{}, err
	}
	return vehicle, nil
}

func GetMotorVehicles(ctx context.Context) ([]models.MotorVehicle, error) {
	ctx, cancel := context.WithTimeout(ctx, vehicleTimeout)
	defer cancel()
	cursor, err := vehicleColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	vehicles := make([]models.MotorVehicle, 0)
	if err := cursor.All(ctx, &vehicles); err != nil {
		return nil, err
	}
	return vehicles, nil
}

func GetMotorVehicleByID(ctx context.Context, id string) (models.MotorVehicle, error) {
	ctx, cancel := context.WithTimeout(ctx, vehicleTimeout)
	defer cancel()
	objID, err := vehicleObjectID(id)
	if err != nil {
		return models.MotorVehicle{}, err
	}
	var vehicle models.MotorVehicle
	if err := vehicleColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&vehicle); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.MotorVehicle{}, ErrMotorVehicleNotFound
		}
		return models.MotorVehicle{}, err
	}
	return vehicle, nil
}

func UpdateMotorVehicle(ctx context.Context, id string, vehicle models.MotorVehicle) (models.MotorVehicle, error) {
	ctx, cancel := context.WithTimeout(ctx, vehicleTimeout)
	defer cancel()
	objID, err := vehicleObjectID(id)
	if err != nil {
		return models.MotorVehicle{}, err
	}
	var updated models.MotorVehicle
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := vehicleColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": vehicleFields(vehicle)}, opts).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.MotorVehicle{}, ErrMotorVehicleNotFound
		}
		return models.MotorVehicle{}, err
	}
	return updated, nil
}

func SetVehicleRider(ctx context.Context, vehicleID, riderID string) (models.MotorVehicle, error) {
	ctx, cancel := context.WithTimeout(ctx, vehicleTimeout)
	defer cancel()
	vehicleObjectID, err := vehicleObjectID(vehicleID)
	if err != nil {
		return models.MotorVehicle{}, err
	}
	riderObjectID, err := riderObjectID(riderID)
	if err != nil {
		return models.MotorVehicle{}, err
	}
	var vehicle models.MotorVehicle
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := vehicleColl().FindOneAndUpdate(ctx, bson.M{"_id": vehicleObjectID}, bson.M{"$set": bson.M{"rider_id": riderObjectID}}, opts).Decode(&vehicle); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.MotorVehicle{}, ErrMotorVehicleNotFound
		}
		return models.MotorVehicle{}, err
	}
	return vehicle, nil
}

func ClearVehicleRider(ctx context.Context, vehicleID string) (models.MotorVehicle, error) {
	ctx, cancel := context.WithTimeout(ctx, vehicleTimeout)
	defer cancel()
	objectID, err := vehicleObjectID(vehicleID)
	if err != nil {
		return models.MotorVehicle{}, err
	}
	var vehicle models.MotorVehicle
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := vehicleColl().FindOneAndUpdate(ctx, bson.M{"_id": objectID}, bson.M{"$unset": bson.M{"rider_id": ""}}, opts).Decode(&vehicle); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.MotorVehicle{}, ErrMotorVehicleNotFound
		}
		return models.MotorVehicle{}, err
	}
	return vehicle, nil
}

func DeleteMotorVehicle(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, vehicleTimeout)
	defer cancel()
	objID, err := vehicleObjectID(id)
	if err != nil {
		return err
	}
	result, err := vehicleColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrMotorVehicleNotFound
	}
	return nil
}
