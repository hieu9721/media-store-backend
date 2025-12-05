package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hieu9721/media-store-backend/config"
	"github.com/hieu9721/media-store-backend/models"
	"github.com/hieu9721/media-store-backend/utils"
)

func CreateAlbum(c *gin.Context) {
	var input models.CreateAlbumInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	collection := config.DB.Collection("albums")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	albumID := utils.GenerateID("alb")

	album := models.Album{
		ID:          albumID,
		UserID:      userID,
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	_, err := collection.InsertOne(ctx, album)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}

	c.JSON(http.StatusCreated, album)
}

func GetAlbums(c *gin.Context) {
	userID := c.GetString("user_id")

	collection := config.DB.Collection("albums")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, gin.H{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch albums"})
		return
	}
	defer cursor.Close(ctx)

	var albums []models.Album
	if err = cursor.All(ctx, &albums); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse albums"})
		return
	}

	c.JSON(http.StatusOK, albums)
}
