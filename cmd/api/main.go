package main

import (
	"context"
	"github.com/carloscasalar/gin-starter/internal/infrastructure/app"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("unable to start the API: %v", err)
		os.Exit(1)
	}

	api := app.New(cfg)
	go func() {
		if err := api.Start(ctx); err != nil {
			log.WithContext(ctx).Errorf("unable to start the API: %v", err)
			os.Exit(1)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := api.Stop(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Info("Server exiting")
}
