package logging_test

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/BlackBX/service-framework/config"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PrinterFunc func(string, ...interface{})

func (p PrinterFunc) Printf(string, ...interface{}) {}

func TestLoggerFactory_LoggerSucceeds(t *testing.T) {
	fxLogger := fx.Logger(PrinterFunc(func(string, ...interface{}) {}))

	cmd := &cobra.Command{}
	cmd.SetOut(ioutil.Discard)

	builder := dependency.NewBuilder(cmd).
		WithService(config.Service).
		WithService(logging.Service).
		WithModule(fxLogger)
	cmd.Flags().String("logger", "nop", "foo bar baz")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error from executing the command, got (%s)", err)
	}

	timesCalled := 0
	invokeFunc := func(logger *zap.Logger) {
		timesCalled++
	}

	app := builder.WithInvoke(invokeFunc).BuildTest(t)
	go app.Run()
	if timesCalled != 1 {
		t.Fatalf("expected logger to be called 1 time, it was called (%d) time(s)", timesCalled)
	}
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("expected no error when stopping app, got (%s)", err)
	}
}

func TestLoggerFactory_LoggerInvalidLogger(t *testing.T) {
	fxLogger := fx.Logger(PrinterFunc(func(string, ...interface{}) {}))

	cmd := &cobra.Command{}
	cmd.SetOut(ioutil.Discard)

	logging.Service.ConfigFunc(cmd.PersistentFlags())
	builder := dependency.NewBuilder(cmd).
		WithService(config.Service).
		WithModule(fxLogger)
	cmd.
		Flags().
		String("logger", "notalogger", "foo bar baz")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error from executing the command, got (%s)", err)
	}

	invokeFunc := func(cfg *viper.Viper) {
		_, err := logging.NewLoggerFactory().Logger(cfg)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	}

	app := builder.WithInvoke(invokeFunc).BuildTest(t)
	app.RequireStop()
	go app.Run()
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("expected no error when stopping app, got (%s)", err)
	}
}

func TestLoggerFactory_LoggerFailingConstructor(t *testing.T) {
	fxLogger := fx.Logger(PrinterFunc(func(string, ...interface{}) {}))

	cmd := &cobra.Command{}
	cmd.SetOut(ioutil.Discard)

	logging.Service.ConfigFunc(cmd.PersistentFlags())
	builder := dependency.NewBuilder(cmd).
		WithService(config.Service).
		WithModule(fxLogger)
	cmd.
		Flags().
		String("logger", "failure", "foo bar baz")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error from executing the command, got (%s)", err)
	}

	invokeFunc := func(cfg *viper.Viper) {
		factory := logging.NewLoggerFactory()
		factory.LoggerConstructors["failure"] = func(options ...zap.Option) (logger *zap.Logger, err error) {
			return nil, errors.New("failure")
		}
		_, err := factory.Logger(cfg)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	}

	app := builder.WithInvoke(invokeFunc).BuildTest(t)
	app.RequireStop()
	go app.Run()
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("expected no error when stopping app, got (%s)", err)
	}
}
