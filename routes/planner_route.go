package routes

import (
	"catatuangbackend/handlers"
	"catatuangbackend/middleware"
	"catatuangbackend/repo"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetPlannerRoutes(r *gin.Engine, db *sql.DB, redis *redis.Client) {
	_repo := repo.PlannerRepository(db)
	h := handlers.PlannerHandlerFunc(_repo)

	protected := r.Group("/planner")
	protected.Use(middleware.AuthMiddleware(redis))
	{
		protected.POST("/create_planner_book", h.CreatePlannerBook)
		protected.GET("/get_planner_book", h.GetPlannerBookList)
		protected.POST("/remove_planner_book", h.RemoveBookItem)
		protected.PUT("/update_book", h.UpdatePlannerBook)
		protected.POST("/create_item_book", h.CreateItemPlanner)
		protected.GET("/get_book/:id", h.GetDetailItemList)
	}
}
