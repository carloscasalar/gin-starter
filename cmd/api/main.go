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
	cfg, err := app.ReadConfig()
	if err != nil {
		log.Errorf("Unable to start the API: %v", err)
		os.Exit(1)
	}

	api := app.New(cfg)
	go func() {
		ctx := context.Background()
		if err := api.Start(ctx); err != nil {
			log.WithContext(ctx).Errorf("Unable to start the API: %v", err)
			os.Exit(1)
		}
	}()

	gracefulShutdownOnSigIntOrTerm(api.Stop)

	log.Info("Server exiting")
}

func gracefulShutdownOnSigIntOrTerm(onStopFn func(ctx context.Context) error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := onStopFn(ctx); err != nil {
		log.Fatalf("Unable to shut down gracefully: %v", err)
	}
	os.Exit(1)
}
