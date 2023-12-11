package main

import (
	"context"
	"fmt"
	"main/configs"
	"main/internal/app/router"
	"main/internal/db"
	"main/internal/services/logger"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log := logger.Init(true)
	logger.Log.Info("Logger initialized")
	defer log.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Log.Panicw("Error loading .env file", "error", err)
	}
	configs.InitEnvConfig()

	_, err = db.Init()

	if err != nil {
		logger.Log.Panicw("Error initializing database", "error", err)
	}

	router := router.Init(gin.Default())

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", configs.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Panicw("Error initializing server", "error", err)
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	logger.Log.Info("Received interrupt signal. Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Panicw("Error shutting down the server", "error", err)
	} else {
		logger.Log.Info("Server gracefully stopped.")
	}
}
