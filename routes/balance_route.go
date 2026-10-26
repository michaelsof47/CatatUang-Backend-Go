package routes

import (
	"catatuangbackend/handlers"
	"catatuangbackend/middleware"
	"catatuangbackend/repo"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetBalanceRoutes(r *gin.Engine, db *sql.DB, redis *redis.Client) {
	_repo := repo.BalanceRepository(db)
	h := handlers.BalancesHandlerFunc(_repo)

	protected := r.Group("/balance")
	protected.Use(middleware.AuthMiddleware(redis))
	{
		protected.PUT("/create_or_update_amount", h.CreateOrUpdateBalancesAmount)
		protected.PATCH("/add_balance_amount", h.AddBalanceAmount)
		protected.GET("/get_balance", h.GetBalanceAmount)
	}
}
