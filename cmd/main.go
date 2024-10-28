package main

import (
	"ads/config"
	"ads/internal/ads/delivery"
	"ads/internal/cache"
	"ads/internal/db"
	"ads/pkg/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// todo: write test files and jmeter
	cfg := config.LoadConfig()

	logger := logger.InitLogger(cfg.LogFilePath, cfg.LogFileName)
	logger.Info("Starting the ad service")

	dbPool := db.NewDBpool(cfg)
	defer dbPool.Close()

	rClient := cache.NewRedisCache("localhost:6379")
	rClient.SetAllADs(context.Background(), dbPool)

	r := gin.Default()
	delivery.RegisterAdsRoutes(r, dbPool, rClient)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("ListenAndServe error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Println("Shutting down server...")
	shutdownTimeout := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	<-ctx.Done()
	logger.Println("Server exiting")

}
