package handlers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	galleryUploadDir = "uploads/gallery"
	galleryPublicDir = "/uploads/gallery"
)

func saveGalleryPhoto(header *multipart.FileHeader) (string, string, int64, error) {
	if header == nil {
		return "", "", 0, errors.New("photo is required")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		return "", "", 0, errors.New("only jpg, jpeg, png, gif, and webp photos are allowed")
	}

	if err := os.MkdirAll(galleryUploadDir, 0755); err != nil {
		return "", "", 0, err
	}

	fileName := bson.NewObjectID().Hex() + ext
	filePath := filepath.Join(galleryUploadDir, fileName)
	if err := saveUploadedFile(header, filePath); err != nil {
		return "", "", 0, err
	}

	return fileName, galleryPublicDir + "/" + fileName, header.Size, nil
}

func saveUploadedFile(header *multipart.FileHeader, path string) error {
	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	output, err := os.Create(path)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, file)
	return err
}

func CreateGallery(c *gin.Context) {
	header, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "photo is required"})
		return
	}

	fileName, url, size, err := saveGalleryPhoto(header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	gallery := models.Gallery{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		FileName:    fileName,
		URL:         url,
		ContentType: header.Header.Get("Content-Type"),
		Size:        size,
	}

	created, err := services.CreateGallery(c.Request.Context(), gallery)
	if err != nil {
		_ = os.Remove(filepath.Join(galleryUploadDir, fileName))
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "gallery photo created successfully", "data": created})
}

func GetGalleries(c *gin.Context) {
	galleries, err := services.GetGalleries(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "gallery photos fetched successfully", "data": galleries})
}

func GetGallery(c *gin.Context) {
	gallery, err := services.GetGalleryByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrGalleryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "gallery photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "gallery photo fetched successfully", "data": gallery})
}

func UpdateGallery(c *gin.Context) {
	gallery, err := services.GetGalleryByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrGalleryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "gallery photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	updated := gallery
	if title := c.PostForm("title"); title != "" {
		updated.Title = title
	}
	updated.Description = c.PostForm("description")

	if header, fileErr := c.FormFile("photo"); fileErr == nil {
		fileName, url, size, saveErr := saveGalleryPhoto(header)
		if saveErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": saveErr.Error()})
			return
		}
		updated.FileName = fileName
		updated.URL = url
		updated.ContentType = header.Header.Get("Content-Type")
		updated.Size = size
	}

	result, err := services.UpdateGallery(c.Request.Context(), c.Param("id"), updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if result.FileName != gallery.FileName {
		_ = os.Remove(filepath.Join(galleryUploadDir, gallery.FileName))
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "gallery photo updated successfully", "data": result})
}

func DeleteGallery(c *gin.Context) {
	gallery, err := services.GetGalleryByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrGalleryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "gallery photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := services.DeleteGallery(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	_ = os.Remove(filepath.Join(galleryUploadDir, gallery.FileName))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "gallery photo deleted successfully"})
}
