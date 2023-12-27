package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

func ConfigureLogger(cfg LogConfig) error {
	if err := setLogLevel(cfg.Level); err != nil {
		return fmt.Errorf("failed to set log-level: %v", err)
	}

	if err := setLogFormat(cfg.Formatter); err != nil {
		return fmt.Errorf("failed to set log format: %v", err)
	}

	return nil
}

func setLogFormat(format string) error {
	switch strings.ToLower(format) {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return fmt.Errorf("unknown log format %q, valid values are 'text' or 'json'", format)
	}

	return nil
}

func setLogLevel(levelName string) error {
	var logLevel log.Level

	switch strings.ToLower(levelName) {
	case "trace":
		logLevel = log.TraceLevel
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "warn":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	case "fatal":
		logLevel = log.FatalLevel
	case "panic":
		logLevel = log.PanicLevel
	default:
		return fmt.Errorf("unknown log level %q, valid values are 'trace', 'debug', 'info', 'warn', 'error', 'fatal' or 'panic'", levelName)
	}

	log.SetLevel(logLevel)
	return nil
}
