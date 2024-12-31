package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"billyb/internal/config"
	"billyb/internal/handlers"
	"billyb/internal/middleware"
	"billyb/pkg/logger"
	"billyb/pkg/tracer"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.NewLogger(cfg.Logging.Level, cfg.Logging.Encoding)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	if cfg.Tracing.Enabled {
		tp, err := tracer.InitTracer("billyb", cfg.Tracing.Endpoint)
		if err != nil {
			log.Fatal("Failed to initialize tracer", zap.Error(err))
		}
		defer tp.Shutdown(context.Background())
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.CORSMiddleware())

	itemHandler := handlers.NewItemHandler(log, &cfg.Server)

	// Administration Routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Service Routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/items", itemHandler.GetItems)
		v1.POST("/items", itemHandler.AddItem)
		v1.DELETE("/items", itemHandler.DeleteAllItems)
		router.DELETE("/items/:id", itemHandler.DeleteItem)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
	}

	go func() {
		log.Info("Starting billyb server.", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start billyb server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal then gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Billyb server was forced to shutdown", zap.Error(err))
	}

	log.Info("Billyb Server exited succesfully.")
}
