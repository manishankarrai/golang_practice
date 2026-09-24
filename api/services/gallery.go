// Package services contains the application services that coordinate gallery data access.
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
	// Keep gallery data isolated from the other application collections.
	galleryDB         = "golang_practice"
	galleryCollection = "gallery"

	// Bound every gallery database operation so a stalled MongoDB request does
	// not keep an HTTP request waiting indefinitely.
	galleryTimeout = 10 * time.Second
)

// ErrGalleryNotFound is returned for invalid IDs and IDs that do not match a
// stored gallery document, allowing handlers to map both cases consistently.
var ErrGalleryNotFound = errors.New("gallery photo not found")

// galleryColl returns the MongoDB collection used by all gallery operations.
func galleryColl() *mongo.Collection {
	return db.Database(galleryDB).Collection(galleryCollection)
}

// galleryContext applies the service-wide timeout to a database operation.
func galleryContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, galleryTimeout)
}

// galleryObjectID converts an API ID into MongoDB's ObjectID format and
// normalizes malformed IDs to the service's not-found error.
func galleryObjectID(id string) (bson.ObjectID, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrGalleryNotFound
	}
	return objID, nil
}

// galleryUpdate builds the fields that may be changed without replacing the
// document or altering its MongoDB-generated identifier.
func galleryUpdate(gallery models.Gallery) bson.M {
	return bson.M{"$set": bson.M{
		"title":        gallery.Title,
		"description":  gallery.Description,
		"file_name":    gallery.FileName,
		"url":          gallery.URL,
		"content_type": gallery.ContentType,
		"size":         gallery.Size,
	}}
}

// CreateGallery assigns the document ID in the service layer before storing
// the new gallery record.
func CreateGallery(ctx context.Context, gallery models.Gallery) (models.Gallery, error) {
	ctx, cancel := galleryContext(ctx)
	defer cancel()

	gallery.ID = bson.NewObjectID()
	if _, err := galleryColl().InsertOne(ctx, gallery); err != nil {
		return models.Gallery{}, err
	}
	return gallery, nil
}

// GetGalleries loads every gallery document and returns an empty, non-nil
// slice when the collection contains no records.
func GetGalleries(ctx context.Context) ([]models.Gallery, error) {
	ctx, cancel := galleryContext(ctx)
	defer cancel()

	cursor, err := galleryColl().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	galleries := make([]models.Gallery, 0)
	if err := cursor.All(ctx, &galleries); err != nil {
		return nil, err
	}
	return galleries, nil
}

// GetGalleryByID validates the API ID, fetches the matching document, and
// translates MongoDB's no-document result into the service error.
func GetGalleryByID(ctx context.Context, id string) (models.Gallery, error) {
	ctx, cancel := galleryContext(ctx)
	defer cancel()

	objID, err := galleryObjectID(id)
	if err != nil {
		return models.Gallery{}, err
	}

	var gallery models.Gallery
	err = galleryColl().FindOne(ctx, bson.M{"_id": objID}).Decode(&gallery)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Gallery{}, ErrGalleryNotFound
		}
		return models.Gallery{}, err
	}
	return gallery, nil
}

// UpdateGallery changes the editable gallery fields and decodes the document
// returned by MongoDB after the update so callers receive the saved values.
func UpdateGallery(ctx context.Context, id string, gallery models.Gallery) (models.Gallery, error) {
	ctx, cancel := galleryContext(ctx)
	defer cancel()

	objID, err := galleryObjectID(id)
	if err != nil {
		return models.Gallery{}, err
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Gallery
	err = galleryColl().FindOneAndUpdate(ctx, bson.M{"_id": objID}, galleryUpdate(gallery), opts).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Gallery{}, ErrGalleryNotFound
		}
		return models.Gallery{}, err
	}
	return updated, nil
}

// DeleteGallery removes one gallery document and reports not-found when the
// filter matched no records.
func DeleteGallery(ctx context.Context, id string) error {
	ctx, cancel := galleryContext(ctx)
	defer cancel()

	objID, err := galleryObjectID(id)
	if err != nil {
		return err
	}

	result, err := galleryColl().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	// A zero DeletedCount means no document matched the filter, so the
	// requested gallery photo never existed; report it as not found.
	if result.DeletedCount == 0 {
		return ErrGalleryNotFound
	}
	return nil
}
