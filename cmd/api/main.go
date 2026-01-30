package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pick-cee/events-api/internal/config"
	"github.com/pick-cee/events-api/internal/database"
	"github.com/pick-cee/events-api/internal/middleware"
)

func main() {
	// connect to database
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := database.Connect(cfg); err != nil {
		log.Println(err)
		panic(err)
	}

	if err := database.Migrate(); err != nil {
		log.Println(err)
		panic(err)
	}

	// start server
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())

	// define routes

	router.GET("/health", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "Events API is up and running",
		})
	})

	srv := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("🚀 Server running on port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := database.Disconnect(); err != nil {
		log.Println("DB disconnect error:", err)
	}

	log.Println("✅ Server exited properly")
}