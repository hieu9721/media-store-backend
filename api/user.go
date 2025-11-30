package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hieu9721/media-store-backend/config"
	"github.com/hieu9721/media-store-backend/models"
	"github.com/hieu9721/media-store-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return config.GetCollection("users")
}

// CreateUser - Tạo user mới
func CreateUser(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var user models.User

    // Validate input
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check email exist
    var existUser models.User
    err := getUserCollection().FindOne(ctx, bson.M{"email": user.Email}).Decode(&existUser)
    if err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
        return
    }

    user.ID = utils.GenerateUserID()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    _, err = getUserCollection().InsertOne(ctx, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "data":    user,
    })
}

// GetUsers - Lấy danh sách users với pagination
func GetUsers(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var users []models.User

    cursor, err := getUserCollection().Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user"})
            return
        }
        users = append(users, user)
    }

    if users == nil {
        users = []models.User{}
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Users fetched successfully",
        "count":   len(users),
        "data":    users,
    })
}

// GetUser - Lấy user theo ID
func GetUser(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    userId := c.Param("id")
    if !utils.IsValidUserID(userId) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    var user models.User
    err := getUserCollection().FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "User fetched successfully",
        "data":    user,
    })
}

// UpdateUser - Cập nhật user
func UpdateUser(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    userId := c.Param("id")
    if !utils.IsValidUserID(userId) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    var updateData models.UpdateUser
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if user exists
    var existUser models.User
    err := getUserCollection().FindOne(ctx, bson.M{"_id": userId}).Decode(&existUser)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
        return
    }

    // Check if email already used by another user
    if updateData.Email != "" && updateData.Email != existUser.Email {
        var emailCheck models.User
        err := getUserCollection().FindOne(ctx, bson.M{"email": updateData.Email}).Decode(&emailCheck)
        if err == nil {
            c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
            return
        }
    }

    // Build update document
    updateDoc := bson.M{
        "$set": bson.M{
            "updated_at": time.Now(),
        },
    }

    if updateData.Name != "" {
        updateDoc["$set"].(bson.M)["name"] = updateData.Name
    }
    if updateData.Email != "" {
        updateDoc["$set"].(bson.M)["email"] = updateData.Email
    }
    if updateData.Phone != "" {
        updateDoc["$set"].(bson.M)["phone"] = updateData.Phone
    }
    if updateData.Avatar != "" {
        updateDoc["$set"].(bson.M)["avatar"] = updateData.Avatar
    }

    _, err = getUserCollection().UpdateOne(ctx, bson.M{"_id": userId}, updateDoc)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    // Get updated user
    var updatedUser models.User
    getUserCollection().FindOne(ctx, bson.M{"_id": userId}).Decode(&updatedUser)

    c.JSON(http.StatusOK, gin.H{
        "message": "User updated successfully",
        "data":    updatedUser,
    })
}

// DeleteUser - Xóa user
func DeleteUser(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    userId := c.Param("id")
    if !utils.IsValidUserID(userId) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
        return
    }

    result, err := getUserCollection().DeleteOne(ctx, bson.M{"_id": userId})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    if result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "User deleted successfully",
    })
}

// SearchUsers - Tìm kiếm users
func SearchUsers(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    searchTerm := c.Query("q")
    if searchTerm == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Search term is required"})
        return
    }

    filter := bson.M{
        "$or": []bson.M{
            {"name": bson.M{"$regex": searchTerm, "$options": "i"}},
            {"email": bson.M{"$regex": searchTerm, "$options": "i"}},
        },
    }

    var users []models.User
    cursor, err := getUserCollection().Find(ctx, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            continue
        }
        users = append(users, user)
    }

    if users == nil {
        users = []models.User{}
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Search completed",
        "count":   len(users),
        "data":    users,
    })
}
