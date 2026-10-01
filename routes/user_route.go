package routes

import (
	"catatuangbackend/handlers"
	"catatuangbackend/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetRoutes(r *gin.Engine, db *sql.DB, redis *redis.Client) {
	h := handlers.UserHandlerFunc(db, redis)

	r.POST("/create_account", h.CreateUserNew)
	r.POST("/login_user", h.LoginUserAccount)
	r.POST("/logout_user", h.LogoutUserAccount)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(redis))
	{
		protected.GET("/get_user/:id", h.GetUserById)
	}
}
