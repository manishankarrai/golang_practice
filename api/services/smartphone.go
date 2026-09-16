// Package services holds the business logic between the HTTP handlers and
// MongoDB. Every function in this file follows the same shape: create a
// deadline-bounded context, convert the caller's hex id into a BSON ObjectID,
// run exactly one database call, then translate driver errors into the domain
// error ErrSmartphoneNotFound.
package services

// Imports used in this file:
//   - context: carries deadlines/cancellation down into the database driver.
//   - errors:  errors.New for the sentinel error, errors.Is to match it.
//   - time:    builds the 10s operationTimeout duration.
//   - db:      db.Database(name) returns the *mongo.Database handle.
//   - models:  the Smartphone struct that maps to a BSON document.
//   - bson:    BSON helpers - bson.M (a query map) and bson.ObjectID.
//   - mongo:   the driver package, exposes the ErrNoDocuments sentinel.
//   - options: per-call options such as "return the updated document".
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
	// Where the documents live. MongoDB creates the database and collection
	// lazily on the first write, so no explicit setup is needed here.
	smartphoneDB         = "golang_practice"
	smartphoneCollection = "smartphones"

	// operationTimeout bounds every database round trip in this service, so a
	// slow or unreachable database can never hang a request forever.
	operationTimeout = 10 * time.Second
)

// ErrSmartphoneNotFound is returned when a smartphone document does not exist.
// It is a sentinel value: callers detect it with errors.Is(err, ErrSmartphoneNotFound)
// whether the miss came from a malformed id or from a missing document.
var ErrSmartphoneNotFound = errors.New("smartphone not found")

// smartphoneColl resolves the collection handle on each call. The underlying
// *mongo.Client is safe for concurrent use and pools connections, so building
// this lightweight handle per call is cheap and avoids shared mutable state.
func smartphoneColl() *mongo.Collection {
	return db.Database(smartphoneDB).Collection(smartphoneCollection)
}

// CreateSmartphone inserts a new smartphone document.
func CreateSmartphone(ctx context.Context, phone models.Smartphone) (models.Smartphone, error) {
	// Derive a child context that is cancelled after operationTimeout - or
	// earlier if the caller's context is cancelled. cancel must always run to
	// release the timer, which is why it is deferred right away.
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	// Generate the primary key up front so the caller receives the id that was
	// actually persisted, without a second round trip.
	phone.ID = bson.NewObjectID()
	if _, err := smartphoneColl().InsertOne(ctx, phone); err != nil {
		// Return a zero value on failure so callers never use a partial doc.
		return models.Smartphone{}, err
	}

	return phone, nil
}

// GetSmartphones returns all smartphone documents, optionally filtered by category.
func GetSmartphones(ctx context.Context, category string) ([]models.Smartphone, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	// An empty bson.M matches every document. When a category is supplied we
	// add a field, which is the BSON equivalent of SQL's WHERE category = ?.
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}

	// Find returns a cursor, not the documents themselves; the query is only
	// executed lazily as the cursor is iterated below.
	cursor, err := smartphoneColl().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	// The cursor holds a server-side resource, so it must always be closed.
	defer cursor.Close(ctx)

	// Decode into a non-nil slice so the JSON response is [] rather than null
	// when there are no matching documents.
	phones := make([]models.Smartphone, 0)
	if err := cursor.All(ctx, &phones); err != nil {
		return nil, err
	}

	return phones, nil
}

// GetSmartphoneByID returns a single smartphone by its hex id.
func GetSmartphoneByID(ctx context.Context, id string) (models.Smartphone, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	// Convert the string id into a BSON ObjectID. A malformed id can never
	// match a real document, so the helper reports it as "not found".
	objID, err := smartphoneObjectID(id)
	if err != nil {
		return models.Smartphone{}, err
	}

	// FindOne returns at most one document; Decode copies it into the struct,
	// mapping BSON fields onto their `bson:"..."` struct tags.
	var phone models.Smartphone
	if err := smartphoneColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&phone); err != nil {
		return models.Smartphone{}, smartphoneError(err)
	}

	return phone, nil
}

// UpdateSmartphone replaces the mutable fields of an existing smartphone.
func UpdateSmartphone(ctx context.Context, id string, phone models.Smartphone) (models.Smartphone, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	objID, err := smartphoneObjectID(id)
	if err != nil {
		return models.Smartphone{}, err
	}

	// $set updates only the listed fields, leaving _id and any other fields
	// untouched. The id itself is deliberately not part of the update.
	update := bson.M{"$set": bson.M{
		"name":     phone.Name,
		"brand":    phone.Brand,
		"category": phone.Category,
		"price":    phone.Price,
		"stock":    phone.Stock,
	}}
	// By default FindOneAndUpdate returns the document as it was BEFORE the
	// update; requesting options.After returns the freshly updated version.
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated models.Smartphone
	if err := smartphoneColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updated); err != nil {
		// A filter that matches nothing surfaces as ErrNoDocuments, which
		// smartphoneError converts into ErrSmartphoneNotFound.
		return models.Smartphone{}, smartphoneError(err)
	}

	return updated, nil
}

// DeleteSmartphone removes a smartphone by its hex id.
func DeleteSmartphone(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	objID, err := smartphoneObjectID(id)
	if err != nil {
		return err
	}

	// DeleteOne reports what happened through the result, not through an error:
	// deleting an id that does not exist is not a driver error, it simply
	// affects 0 documents.
	res, err := smartphoneColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	// Translate "matched nothing" into the domain not-found error.
	if res.DeletedCount == 0 {
		return ErrSmartphoneNotFound
	}

	return nil
}

// smartphoneObjectID parses a hex id into a BSON ObjectID, mapping malformed
// ids to ErrSmartphoneNotFound so callers never leak driver errors.
func smartphoneObjectID(id string) (bson.ObjectID, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, ErrSmartphoneNotFound
	}

	return objID, nil
}

// smartphoneError maps driver-level errors to domain errors.
func smartphoneError(err error) error {
	// errors.Is unwraps the error chain, so it is safer than a plain ==.
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrSmartphoneNotFound
	}

	return err
}
