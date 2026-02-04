package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pick-cee/events-api/internal/config"
	"github.com/pick-cee/events-api/internal/database"
	"github.com/pick-cee/events-api/internal/middleware"
)

func main() {
	// connect to database
	cfg := config.Load()
	log.Println("✅ Configuration loaded")

	_, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := database.Connect(cfg); err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("❌ Failed to run migrations:", err)
	}

	if err := database.ConnectRedis(cfg); err != nil {
		log.Fatal("❌ Failed to connect to Redis:", err)
	}

	// Setup routes
	router := setupRoutes(cfg)

	router.Use(middleware.CORSMiddleware())

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
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

	if err := database.DisconnectRedis(); err != nil {
		log.Println("Redis disconnect error:", err)
	}

	log.Println("✅ Server exited properly")
}
