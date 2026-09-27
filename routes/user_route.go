package routes

import (
	"catatuangbackend/handlers"
	"catatuangbackend/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine, db *sql.DB) {
	h := handlers.UserHandlerFunc(db)

	r.POST("/create_account", h.CreateUserNew)
	r.POST("/login_user", h.LoginUserAccount)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/get_user/:id", h.GetUserById)
	}
}
