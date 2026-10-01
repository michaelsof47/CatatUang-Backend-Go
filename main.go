package main

import (
	"catatuangbackend/config"
	"catatuangbackend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDatabase()
	defer db.Close()

	redisdb := config.InitRedis()
	defer redisdb.Close()

	r := gin.Default()
	routes.SetRoutes(r, db, redisdb)

	r.Run(":8080")
}
