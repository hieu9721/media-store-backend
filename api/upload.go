package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

var allowedVideoExtensions = map[string]bool{
	".mp4":  true,
	".avi":  true,
	".mov":  true,
	".mkv":  true,
	".webm": true,
}

// UploadAvatar - Upload avatar cho user, lưu trong thư mục uploads/avatars
func UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded or invalid field name. Use 'image' as field name",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Only PNG, JPG, JPEG, and GIF are allowed",
		})
		return
	}

	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 5MB limit",
		})
		return
	}

	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
		if baseURL == "http://localhost:" {
			baseURL = "http://localhost:8080"
		}
	}

	imageURL := fmt.Sprintf("%s/uploads/avatars/%s", baseURL, filename)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Avatar uploaded successfully",
		"user_id":  userID,
		"filename": filename,
		"url":      imageURL,
		"size":     file.Size,
	})
}

// UploadUserImage - Upload ảnh vào thư viện cá nhân của user, lưu trong thư mục uploads/gallery/{user_id}
func UploadUserImage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded or invalid field name. Use 'image' as field name",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Only PNG, JPG, JPEG, and GIF are allowed",
		})
		return
	}

	maxSize := int64(10 * 1024 * 1024) // 10MB for gallery images
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 10MB limit",
		})
		return
	}

	uploadDir := fmt.Sprintf("uploads/gallery/%s", userIDStr)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
		if baseURL == "http://localhost:" {
			baseURL = "http://localhost:8080"
		}
	}

	imageURL := fmt.Sprintf("%s/uploads/gallery/%s/%s", baseURL, userIDStr, filename)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Image uploaded to gallery successfully",
		"user_id":  userID,
		"filename": filename,
		"url":      imageURL,
		"size":     file.Size,
	})
}

// UploadVideo - Upload video vào thư viện của user, lưu trong thư mục uploads/videos/{user_id}
func UploadVideo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded or invalid field name. Use 'video' as field name",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedVideoExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Only MP4, AVI, MOV, MKV, and WEBM are allowed",
		})
		return
	}

	maxSize := int64(500 * 1024 * 1024) // 100MB for videos
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 500MB limit",
		})
		return
	}

	uploadDir := fmt.Sprintf("uploads/videos/%s", userIDStr)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create upload directory",
		})
		return
	}

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL := fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
		if baseURL == "http://localhost:" {
			baseURL = "http://localhost:8080"
		}
	}

	videoURL := fmt.Sprintf("%s/uploads/videos/%s/%s", baseURL, userIDStr, filename)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Video uploaded successfully",
		"user_id":  userID,
		"filename": filename,
		"url":      videoURL,
		"size":     file.Size,
	})
}
