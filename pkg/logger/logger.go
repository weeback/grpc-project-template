package logger

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ContextLoggerType uint16

const ContextLogger ContextLoggerType = 1

var (
	logOnce sync.Once
	logger  *zap.Logger
)

func SetLogEntry(zlg *zap.Logger) {
	if zlg == nil {
		fmt.Println("entry logger set is unavailable, skipping logger setup")
		return
	}

	// cache the logger
	buffered := logger
	// Set the logger by the given zap.Logger
	logger = zlg
	// Flush any buffered entry entries
	if buffered != nil {
		if err := buffered.Sync(); err != nil {
			fmt.Printf("failed to sync logger: %v\n", err)
		}
	}
}

func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ContextLogger).(*zap.Logger); ok {
		return logger
	}
	return NewEntry()
}

func SetLoggerToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextLogger, logger)
}

func NewEntry() *zap.Logger {
	logOnce.Do(func() {
		if logger != nil {
			// If logger is already initialized, skip re-initialization
			return
		}
		// Renew the logger default
		val := os.Getenv("DEPLOYMENT_ENVIRONMENT")
		switch val {
		case "development", "dev":
			// Initialize logger
			config := zap.NewDevelopmentConfig()
			config.Encoding = "json"
			dev, err := config.Build()
			if err != nil {
				logger = zap.NewExample()
				logger.Debug("failed to initialize production logger", zap.Error(err))
			} else {

				logger = dev
			}
		default:
			// Consider production environment
			// Initialize logger
			prod, err := zap.NewProduction()
			if err != nil {
				logger = zap.NewExample()
				logger.Debug("failed to initialize production logger", zap.Error(err))
			} else {
				logger = prod
			}
		}
	})

	return logger.With(
		zap.Time("timestamp", time.Now()),
		zap.String("environment", os.Getenv("DEPLOYMENT_ENVIRONMENT")))
}
