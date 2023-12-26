package main

import (
	"github.com/carloscasalar/gin-starter/internal/infrastructure/app"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	cfg, err := app.ReadConfig()
	if err != nil {
		log.Errorf("unable to start the app: %v", err)
		os.Exit(1)
	}
	log.Infof("service configuration %v", cfg)
}
