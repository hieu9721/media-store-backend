// Package api provides HTTP handlers for file uploads with metadata extraction
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rwcarlsen/goexif/exif"
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

// LocationInfo - Thông tin vị trí chi tiết
type LocationInfo struct {
	Country     string `json:"country,omitempty"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
	District    string `json:"district,omitempty"`
	Road        string `json:"road,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// ImageMetadata - Thông tin metadata của hình ảnh
type ImageMetadata struct {
	DateTime     int64         `json:"date_time,omitempty"`
	Latitude     float64       `json:"latitude,omitempty"`
	Longitude    float64       `json:"longitude,omitempty"`
	Location     *LocationInfo `json:"location,omitempty"`
	CameraMake   string        `json:"camera_make,omitempty"`
	CameraModel  string        `json:"camera_model,omitempty"`
	Width        int           `json:"width,omitempty"`
	Height       int           `json:"height,omitempty"`
	Orientation  int           `json:"orientation,omitempty"`
	Flash        string        `json:"flash,omitempty"`
	FocalLength  string        `json:"focal_length,omitempty"`
	FNumber      string        `json:"f_number,omitempty"`
	ExposureTime string        `json:"exposure_time,omitempty"`
	ISO          int           `json:"iso,omitempty"`
}

// NominatimResponse - Response từ Nominatim API
type NominatimResponse struct {
	DisplayName string `json:"display_name"`
	Address     struct {
		Country    string `json:"country"`
		State      string `json:"state"`
		City       string `json:"city"`
		Town       string `json:"town"`
		Village    string `json:"village"`
		County     string `json:"county"`
		District   string `json:"district"`
		Road       string `json:"road"`
		PostalCode string `json:"postcode"`
	} `json:"address"`
}

// reverseGeocode - Chuyển đổi tọa độ GPS thành địa chỉ
func reverseGeocode(lat, lon float64) *LocationInfo {
	// Sử dụng Nominatim API (OpenStreetMap)
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&zoom=18&addressdetails=1", lat, lon)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	// Nominatim yêu cầu User-Agent
	req.Header.Set("User-Agent", "MediaStoreBackend/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var result NominatimResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil
	}

	// Xây dựng LocationInfo
	location := &LocationInfo{
		Country:     result.Address.Country,
		State:       result.Address.State,
		PostalCode:  result.Address.PostalCode,
		Road:        result.Address.Road,
		DisplayName: result.DisplayName,
	}

	// Ưu tiên City > Town > Village
	if result.Address.City != "" {
		location.City = result.Address.City
	} else if result.Address.Town != "" {
		location.City = result.Address.Town
	} else if result.Address.Village != "" {
		location.City = result.Address.Village
	}

	// Ưu tiên District > County
	if result.Address.District != "" {
		location.District = result.Address.District
	} else if result.Address.County != "" {
		location.District = result.Address.County
	}

	return location
}

// extractImageMetadata - Trích xuất metadata từ file ảnh
func extractImageMetadata(filePath string) *ImageMetadata {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Đọc EXIF data
	exifData, err := exif.Decode(file)
	if err != nil {
		return nil // Không có EXIF data hoặc lỗi
	}

	metadata := &ImageMetadata{}

	// Lấy thời gian chụp
	if dateTime, err := exifData.DateTime(); err == nil {
		metadata.DateTime = dateTime.Unix()
	}

	// Lấy tọa độ GPS
	if lat, long, err := exifData.LatLong(); err == nil {
		metadata.Latitude = lat
		metadata.Longitude = long

		// Thực hiện reverse geocoding để lấy thông tin địa chỉ
		if location := reverseGeocode(lat, long); location != nil {
			metadata.Location = location
		}
	}

	// Lấy thông tin camera
	if make, err := exifData.Get(exif.Make); err == nil {
		if makeStr, err := make.StringVal(); err == nil {
			metadata.CameraMake = strings.TrimSpace(makeStr)
		}
	}

	if model, err := exifData.Get(exif.Model); err == nil {
		if modelStr, err := model.StringVal(); err == nil {
			metadata.CameraModel = strings.TrimSpace(modelStr)
		}
	}

	// Lấy kích thước ảnh
	if width, err := exifData.Get(exif.PixelXDimension); err == nil {
		if w, err := width.Int(0); err == nil {
			metadata.Width = w
		}
	}

	if height, err := exifData.Get(exif.PixelYDimension); err == nil {
		if h, err := height.Int(0); err == nil {
			metadata.Height = h
		}
	}

	// Lấy orientation
	if orientation, err := exifData.Get(exif.Orientation); err == nil {
		if o, err := orientation.Int(0); err == nil {
			metadata.Orientation = o
		}
	}

	// Lấy thông tin flash
	if flash, err := exifData.Get(exif.Flash); err == nil {
		if f, err := flash.Int(0); err == nil {
			if f == 0 {
				metadata.Flash = "No Flash"
			} else {
				metadata.Flash = "Flash Fired"
			}
		}
	}

	// Lấy focal length
	if focalLength, err := exifData.Get(exif.FocalLength); err == nil {
		if num, denom, err := focalLength.Rat2(0); err == nil && denom != 0 {
			metadata.FocalLength = fmt.Sprintf("%.1fmm", float64(num)/float64(denom))
		}
	}

	// Lấy F-Number (aperture)
	if fNumber, err := exifData.Get(exif.FNumber); err == nil {
		if num, denom, err := fNumber.Rat2(0); err == nil && denom != 0 {
			metadata.FNumber = fmt.Sprintf("f/%.1f", float64(num)/float64(denom))
		}
	}

	// Lấy exposure time
	if exposureTime, err := exifData.Get(exif.ExposureTime); err == nil {
		if num, denom, err := exposureTime.Rat2(0); err == nil && denom != 0 {
			if num < denom {
				metadata.ExposureTime = fmt.Sprintf("1/%d", denom/num)
			} else {
				metadata.ExposureTime = fmt.Sprintf("%d", num/denom)
			}
		}
	}

	// Lấy ISO
	if isoSpeed, err := exifData.Get(exif.ISOSpeedRatings); err == nil {
		if iso, err := isoSpeed.Int(0); err == nil {
			metadata.ISO = iso
		}
	}

	return metadata
}

// UploadAvatar - Upload avatar cho user, lưu trong thư mục uploads/uid_xxx/avatars
func UploadAvatar(c *gin.Context) {
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

	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 5MB limit",
		})
		return
	}

	uploadDir := fmt.Sprintf("uploads/uid_%s/avatars", userIDStr)
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

	imageURL := fmt.Sprintf("%s/uploads/uid_%s/avatars/%s", baseURL, userIDStr, filename)

	// Trích xuất metadata từ ảnh
	metadata := extractImageMetadata(filepath)

	response := gin.H{
		"message":  "Avatar uploaded successfully",
		"user_id":  userID,
		"filename": filename,
		"url":      imageURL,
		"size":     file.Size,
	}

	// Thêm metadata nếu có
	if metadata != nil {
		response["metadata"] = metadata
	}

	c.JSON(http.StatusOK, response)
}

// UploadUserImage - Upload ảnh vào thư viện cá nhân của user, lưu trong thư mục uploads/uid_xxx/gallery
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

	uploadDir := fmt.Sprintf("uploads/uid_%s/gallery", userIDStr)
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

	imageURL := fmt.Sprintf("%s/uploads/uid_%s/gallery/%s", baseURL, userIDStr, filename)

	// Trích xuất metadata từ ảnh
	metadata := extractImageMetadata(filepath)

	response := gin.H{
		"message":  "Image uploaded to gallery successfully",
		"user_id":  userID,
		"filename": filename,
		"url":      imageURL,
		"size":     file.Size,
	}

	// Thêm metadata nếu có
	if metadata != nil {
		response["metadata"] = metadata
	}

	c.JSON(http.StatusOK, response)
}

// UploadVideo - Upload video vào thư viện của user, lưu trong thư mục uploads/uid_xxx/videos
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

	uploadDir := fmt.Sprintf("uploads/uid_%s/videos", userIDStr)
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

	videoURL := fmt.Sprintf("%s/uploads/uid_%s/videos/%s", baseURL, userIDStr, filename)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Video uploaded successfully",
		"user_id":  userID,
		"filename": filename,
		"url":      videoURL,
		"size":     file.Size,
	})
}
