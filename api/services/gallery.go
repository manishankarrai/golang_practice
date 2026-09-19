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
	galleryDB         = "golang_practice"
	galleryCollection = "gallery"
	galleryTimeout    = 10 * time.Second
)

var ErrGalleryNotFound = errors.New("gallery photo not found")

func galleryColl() *mongo.Collection {
	return db.Database(galleryDB).Collection(galleryCollection)
}

func galleryContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, galleryTimeout)
}

func galleryObjectID(id string) (bson.ObjectID, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, ErrGalleryNotFound
	}
	return objID, nil
}

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

func CreateGallery(ctx context.Context, gallery models.Gallery) (models.Gallery, error) {
	ctx, cancel := galleryContext(ctx)
	defer cancel()

	gallery.ID = bson.NewObjectID()
	if _, err := galleryColl().InsertOne(ctx, gallery); err != nil {
		return models.Gallery{}, err
	}
	return gallery, nil
}

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
	if result.DeletedCount == 0 {
		return ErrGalleryNotFound
	}
	return nil
}
