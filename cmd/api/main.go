package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pick-cee/events-api/internal/config"
	"github.com/pick-cee/events-api/internal/database"
)

func main() {
	// connect to database
	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Println(err)
		panic(err)
	}

	defer database.Disconnect()

	if err := database.Migrate(); err != nil {
		log.Println(err)
		panic(err)
	}

	// start server
	router := gin.Default()

	// define routes

	router.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "Events API is up and running",
		})
	})

	router.Run(":" + cfg.Port)
}