package main

import (
	"catatuangbackend/config"
	"catatuangbackend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDatabase()
	defer db.Close()

	redisdb := config.InitRedis()
	defer redisdb.Close()

	// for testing otp code
	os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "127.0.0.1:9099")

	firebaseApp, err := config.InitFirebase()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	routes.SetRoutes(r, db, redisdb, firebaseApp)

	r.Run(":8080")
}
