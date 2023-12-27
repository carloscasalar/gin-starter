package main

import (
	"context"
	"fmt"
	"github.com/carloscasalar/gin-starter/internal/infrastructure/app"
	"github.com/carloscasalar/gin-starter/internal/infrastructure/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := app.ReadConfig()
	if err != nil {
		log.Errorf("unable to start the app: %v", err)
		os.Exit(1)
	}
	if err := app.ConfigureLogger(cfg.Log); err != nil {
		log.Errorf("unable to configure logger: %v", err)
		os.Exit(1)
	}
	log.Debugf("api configuration %v", cfg)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		middleware.StructuredLogger(),
		gin.Recovery(),
	)
	v1 := router.Group("/v1")
	v1.GET("/status", statusHandler)

	port := fmt.Sprintf(":%v", cfg.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Info("Server exiting")
}

func statusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "server is ready and healthy"})
}
