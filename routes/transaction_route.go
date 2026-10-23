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
		protected.GET("/get_categories", h.GetCategoriesList)
		protected.POST("/get_category_photo", h.GetCategoryPhotoItem)
		protected.POST("/create_transaction", h.CreateTransaction)
		protected.GET("/get_transactions", h.GetTransactionList)
		protected.DELETE("/remove_category", h.RemoveCategory)
	}
}
