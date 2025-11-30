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

    // Serve static files (uploaded images)
    router.Static("/uploads", "./uploads")

    // API v1 routes
    v1 := router.Group("/api/v1")
    {
        // Auth routes (public)
        auth := v1.Group("/auth")
        {
            auth.POST("/register", api.Register)
            auth.POST("/login", api.Login)
        }

        // Protected routes (require authentication)
        protected := v1.Group("")
        protected.Use(middleware.AuthRequired())
        {
            // Current user
            protected.GET("/me", api.GetCurrentUser)

            // Upload routes
            upload := protected.Group("/upload")
            {
                upload.POST("/avatar", api.UploadAvatar)           // Upload avatar
                upload.POST("/image", api.UploadUserImage)         // Upload image to user gallery
                upload.POST("/video", api.UploadVideo)             // Upload video
            }

            // User routes (protected)
            users := protected.Group("/users")
            {
                users.GET("", api.GetUsers)
                users.GET("/search", api.SearchUsers)
                users.GET("/:id", api.GetUser)
                users.PUT("/:id", api.UpdateUser)
            }

            // Admin only routes
            admin := users.Group("")
            admin.Use(middleware.AdminRequired())
            {
                admin.POST("", api.CreateUser)
                admin.DELETE("/:id", api.DeleteUser)
            }
        }
    }

    return router
}
