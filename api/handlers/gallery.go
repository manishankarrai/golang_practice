package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
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
	log.Printf("saveGalleryPhoto: start")
	if header == nil {
		log.Printf("saveGalleryPhoto: photo is required")
		return "", "", 0, errors.New("photo is required")
	}

	log.Printf("saveGalleryPhoto: filename=%q size=%d contentType=%q", header.Filename, header.Size, header.Header.Get("Content-Type"))

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		log.Printf("saveGalleryPhoto: rejected extension %q for filename=%q", ext, header.Filename)
		return "", "", 0, errors.New("only jpg, jpeg, png, gif, and webp photos are allowed")
	}

	if err := os.MkdirAll(galleryUploadDir, 0755); err != nil {
		log.Printf("saveGalleryPhoto: mkdir %q failed: %v", galleryUploadDir, err)
		return "", "", 0, err
	}

	fileName := bson.NewObjectID().Hex() + ext
	filePath := filepath.Join(galleryUploadDir, fileName)
	if err := saveUploadedFile(header, filePath); err != nil {
		log.Printf("saveGalleryPhoto: save file %q failed: %v", filePath, err)
		return "", "", 0, err
	}

	log.Printf("saveGalleryPhoto: saved file=%q size=%d", filePath, header.Size)
	return fileName, galleryPublicDir + "/" + fileName, header.Size, nil
}

// saveUploadedFile streams the multipart upload described by header to path.
//
// Unlike a plain io.Copy, it guarantees the destination is either a complete
// file or absent: the copied byte count is verified against header.Size, and
// any failure (open, partial write, close) removes the destination so callers
// never index a truncated or corrupt file. Errors are wrapped with the
// operation and path for easier diagnosis.
func saveUploadedFile(header *multipart.FileHeader, path string) (err error) {
	if header == nil {
		return errors.New("file header is nil")
	}

	src, err := header.Open()
	if err != nil {
		return fmt.Errorf("open uploaded file %q: %w", header.Filename, err)
	}
	defer src.Close()

	// 0644: owner read/write, group and others read-only.
	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}

	// On any error, make sure the writer is closed and the partial file is gone.
	defer func() {
		if err != nil {
			_ = dst.Close()
			_ = os.Remove(path)
		}
	}()

	written, err := io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("write %q: %w", path, err)
	}
	if written != header.Size {
		return fmt.Errorf("incomplete upload for %q: expected %d bytes, wrote %d", header.Filename, header.Size, written)
	}

	// Close explicitly so buffered/OS write errors reach the caller.
	if err := dst.Close(); err != nil {
		return fmt.Errorf("close %q: %w", path, err)
	}

	return nil
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

// UpdateGallery updates a gallery record and optionally replaces its photo.
// It returns 404 when the gallery does not exist, 400 for an invalid upload,
// and removes the previous photo after a successful replacement.
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

// DeleteGallery removes a gallery record and its associated photo file.
// It returns 404 when the gallery does not exist and 500 when deletion fails.
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
