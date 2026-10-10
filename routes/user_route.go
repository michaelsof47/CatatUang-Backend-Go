package routes

import (
	"catatuangbackend/handlers"
	"catatuangbackend/middleware"
	"database/sql"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetRoutes(r *gin.Engine, db *sql.DB, redis *redis.Client, firebaseApp *firebase.App) {
	h := handlers.UserHandlerFunc(db, redis, firebaseApp)

	r.POST("/create_account", h.CreateUserNew)
	r.POST("/login_user", h.LoginUserAccount)
	r.POST("/logout_user", h.LogoutUserAccount)
	r.POST("/verify_otp", h.VerifyFirebaseToken)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(redis))
	{
		protected.GET("/get_user/:id", h.GetUserById)
		protected.PUT("/update_password", h.UpdatePassword)
		protected.PUT("/update_user", h.UpdateUserAccount)
		protected.PUT("/update_profile_image", h.UpdateProfileImage)
	}
}
