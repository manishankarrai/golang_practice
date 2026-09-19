package services

import (
	"context"
	"errors"
	"math"
	"time"

	"practice/api/db"
	"practice/api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrRideNotFound = errors.New("ride not found")
var ErrVehicleUnavailable = errors.New("vehicle is unavailable")

func rideColl() *mongo.Collection {
	return db.Database("golang_practice").Collection("rides")
}

func rideObjectID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrRideNotFound
	}
	return objectID, nil
}

func distanceKM(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKM = 6371.0
	lat1Rad, lat2Rad := lat1*math.Pi/180, lat2*math.Pi/180
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return earthRadiusKM * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func GetAvailableVehicles(ctx context.Context, latitude, longitude, radiusKM float64) ([]models.AvailableVehicle, error) {
	vehicles, err := GetMotorVehicles(ctx)
	if err != nil {
		return nil, err
	}
	available := make([]models.AvailableVehicle, 0)
	for _, vehicle := range vehicles {
		if vehicle.RiderID == nil {
			continue
		}
		location, _, err := GetVehicleLocation(ctx, vehicle.ID.Hex())
		if err != nil {
			if errors.Is(err, ErrVehicleLocationNotFound) {
				continue
			}
			return nil, err
		}
		distance := distanceKM(latitude, longitude, location.Latitude, location.Longitude)
		if distance <= radiusKM {
			available = append(available, models.AvailableVehicle{Vehicle: vehicle, Location: location, Distance: distance})
		}
	}
	return available, nil
}

func BookRide(ctx context.Context, customerID, vehicleID string, latitude, longitude float64) (models.Ride, error) {
	if _, err := GetCustomerByID(ctx, customerID); err != nil {
		return models.Ride{}, err
	}
	vehicle, err := GetMotorVehicleByID(ctx, vehicleID)
	if err != nil {
		return models.Ride{}, err
	}
	if vehicle.RiderID == nil {
		return models.Ride{}, ErrVehicleUnavailable
	}
	location, _, err := GetVehicleLocation(ctx, vehicleID)
	if err != nil {
		return models.Ride{}, ErrVehicleUnavailable
	}
	if distanceKM(latitude, longitude, location.Latitude, location.Longitude) > 50 {
		return models.Ride{}, ErrVehicleUnavailable
	}
	customerObjectID, err := customerObjectID(customerID)
	if err != nil {
		return models.Ride{}, err
	}
	vehicleObjectID, err := bson.ObjectIDFromHex(vehicleID)
	if err != nil {
		return models.Ride{}, ErrVehicleUnavailable
	}
	ride := models.Ride{ID: bson.NewObjectID(), CustomerID: customerObjectID, VehicleID: vehicleObjectID, Latitude: latitude, Longitude: longitude, Status: "requested", CreatedAt: time.Now().UTC()}
	if _, err := rideColl().InsertOne(ctx, ride); err != nil {
		return models.Ride{}, err
	}
	return ride, nil
}

func GetRides(ctx context.Context) ([]models.Ride, error) {
	cursor, err := rideColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	rides := make([]models.Ride, 0)
	if err := cursor.All(ctx, &rides); err != nil {
		return nil, err
	}
	return rides, nil
}

func GetRideByID(ctx context.Context, id string) (models.Ride, error) {
	objectID, err := rideObjectID(id)
	if err != nil {
		return models.Ride{}, err
	}
	var ride models.Ride
	if err := rideColl().FindOne(ctx, bson.M{"_id": objectID}).Decode(&ride); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Ride{}, ErrRideNotFound
		}
		return models.Ride{}, err
	}
	return ride, nil
}
