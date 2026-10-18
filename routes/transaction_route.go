package routes

import (
	"catatuangbackend/handlers"
	"catatuangbackend/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetTransactionRoutes(r *gin.Engine, db *sql.DB, redis *redis.Client) {
	h := handlers.TransactionHandlerFunc(db)

	protected := r.Group("/transaction")
	protected.Use(middleware.AuthMiddleware(redis))
	{
		protected.POST("/create_category", h.CreateCategory)
	}
}
