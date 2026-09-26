package routes

import (
	"catatuangbackend/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine, db *sql.DB) {
	h := handlers.UserHandlerFunc(db)

	r.POST("/create_account", h.CreateUserNew)
	r.GET("/get_user/:id", h.GetUserById)
}
