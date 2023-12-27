package app

import (
	"context"
	"fmt"
	"github.com/carloscasalar/gin-starter/internal/infrastructure/middleware"
	"github.com/carloscasalar/gin-starter/internal/infrastructure/status"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Instance struct {
	config *Config
	srv    *http.Server
}

func New(config *Config) *Instance {
	return &Instance{
		config: config,
	}
}

func (i *Instance) Start(ctx context.Context) error {
	if err := ConfigureLogger(i.config.Log); err != nil {
		return fmt.Errorf("unable to configure logger: %w", err)
	}
	log.WithContext(ctx).Debugf("api configuration %v", i.config)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		middleware.StructuredLogger(),
		gin.Recovery(),
	)

	v1 := router.Group("/v1")
	v1.GET("/status", status.Handler)

	port := fmt.Sprintf(":%v", i.config.Port)
	i.srv = &http.Server{
		Addr:    port,
		Handler: router,
	}

	if err := i.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (i *Instance) Stop(ctx context.Context) error {
	if err := i.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
