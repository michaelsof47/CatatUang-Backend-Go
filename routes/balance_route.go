package routes

import (
	"catatuangbackend/handlers"
	"catatuangbackend/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetBalanceRoutes(r *gin.Engine, db *sql.DB, redis *redis.Client) {
	h := handlers.BalancesHandlerFunc(db)

	protected := r.Group("/balance")
	protected.Use(middleware.AuthMiddleware(redis))
	{
		protected.PUT("/create_or_update_amount", h.CreateOrUpdateBalancesAmount)
	}
}
