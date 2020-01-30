package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BlackBX/service-framework/dependency"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Service is the exported variable that can be used by the framework package
var Service = dependency.Service{
	Dependencies: fx.Provide(
		NewPrintLogger,
		fx.Annotated{
			Group:  "middleware",
			Target: NewMidlleware,
		},
	),
	ConfigFunc: func(set dependency.FlagSet) {
		set.String("app-name", filepath.Base(os.Args[0]), "The name of the application being configured")
		set.String("app-version", "dev", "The version of the application being configured")
		set.String("environment", "test", "The environment that the application is deployed in")
		set.String("logger", "development", "Whether to log in development mode.")
	},
	Constructor: NewLoggerFactory().Logger,
}

// LoggerConstructor is a type that can give you an instance of a logger
type LoggerConstructor func(options ...zap.Option) (*zap.Logger, error)

// NewLoggerFactory will create a new instance of a logger factory
func NewLoggerFactory() LoggerFactory {
	return LoggerFactory{
		LoggerConstructors: map[string]LoggerConstructor{
			"production":  zap.NewProduction,
			"development": zap.NewDevelopment,
			"nop": func(options ...zap.Option) (logger *zap.Logger, err error) {
				return zap.NewNop(), nil
			},
		},
	}
}

// LoggerFactory is a type that can create instances of loggers
type LoggerFactory struct {
	LoggerConstructors map[string]LoggerConstructor
}

// Logger creates a new instance of a *zap.Logger
func (f LoggerFactory) Logger(settings dependency.ConfigGetter) (*zap.Logger, error) {
	options := []zap.Option{
		zap.Fields(
			zap.String("app-name", settings.GetString("app-name")),
			zap.String("app-version", settings.GetString("app-version")),
			zap.String("environment", settings.GetString("environment")),
		),
	}

	loggerType := settings.GetString("logger")
	loggerConstructor, ok := f.LoggerConstructors[loggerType]
	if !ok {
		return nil, fmt.Errorf("the logger type (%s), is not a valid logger", loggerType)
	}
	logger, err := loggerConstructor(options...)
	if err != nil {
		return nil, fmt.Errorf("could not create instance of logger, got error (%w)", err)
	}
	return logger, nil
}

// NewPrintLogger creates a new instance of the PrintLogger
func NewPrintLogger(logger *zap.Logger) PrintLogger {
	return PrintLogger{
		Logger: logger,
	}
}

// PrintLogger implements a generalised logging interface
type PrintLogger struct {
	Logger *zap.Logger
}

// Println prints the arguments to the zap logger
func (p PrintLogger) Println(arguments ...interface{}) {
	p.Logger.Info(fmt.Sprintln(arguments...))
}
