package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hieu9721/media-store-backend/api"
	"github.com/hieu9721/media-store-backend/middleware"
)

func SetupRoutes() *gin.Engine {
    router := gin.Default()

    // Middleware
    router.Use(middleware.CORS())
    router.Use(middleware.ErrorHandler())

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "OK",
            "message": "Server is running",
        })
    })

    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // User routes
        users := v1.Group("/users")
        {
            users.POST("", api.CreateUser)
            users.GET("", api.GetUsers)
            users.GET("/search", api.SearchUsers)
            users.GET("/:id", api.GetUser)
            users.PUT("/:id", api.UpdateUser)
            users.DELETE("/:id", api.DeleteUser)
        }
    }

    return router
}
