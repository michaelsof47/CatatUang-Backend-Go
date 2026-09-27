package main

import (
	"catatuangbackend/config"
	"catatuangbackend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDatabase()
	defer db.Close()

	r := gin.Default()
	routes.SetRoutes(r, db)

	r.Run(":8080")
}
