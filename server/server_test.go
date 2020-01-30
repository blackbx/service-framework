package server_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/BlackBX/service-framework/config"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/server"
	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

type stubServer struct {
	listenAndServe func() error
	shutdown       func(ctx context.Context) error
}

func (s stubServer) ListenAndServe() error {
	return s.listenAndServe()
}

func (s stubServer) Shutdown(ctx context.Context) error {
	return s.shutdown(ctx)
}

type PrinterFunc func(string, ...interface{})

func (p PrinterFunc) Printf(string, ...interface{}) {}

func TestInvoke(t *testing.T) {
	logger := PrinterFunc(func(format string, args ...interface{}) {
		t.Logf(format, args...)
	})

	serverWg := new(sync.WaitGroup)
	testWg := new(sync.WaitGroup)
	testWg.Add(1) // add one to the test waitgroup so the test waits for the server to be shutdown
	serveCalled := atomic.NewInt32(0)
	shutdownCalled := atomic.NewInt32(0)
	srv := stubServer{
		listenAndServe: func() error {
			serveCalled.Add(1)
			serverWg.Add(1)
			testWg.Done()
			serverWg.Wait() // wait for stop to be called
			return nil
		},
		shutdown: func(ctx context.Context) error {
			shutdownCalled.Add(1)
			serverWg.Done()
			return nil
		},
	}
	app := fxtest.New(t,
		fx.Logger(logger),
		fx.Provide(
			zap.NewNop,
			func() server.Server {
				return srv
			},
		),
		fx.Invoke(server.Invoke),
	)
	go app.Run()
	testWg.Wait() // wait for the app to be started before stopping
	if err := app.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}

	if serveCalled.Load() != 1 {
		t.Fatalf("expected serve to be called once, it was called (%d) times", serveCalled)
	}
	if shutdownCalled.Load() != 1 {
		t.Fatalf("expected shutdown to be called once, called (%d) times", shutdownCalled)
	}
}

func TestNew(t *testing.T) {
	cmd := &cobra.Command{}
	handler := handlers.MethodHandler{}
	server.Service.ConfigFunc(cmd.PersistentFlags())
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	logger := PrinterFunc(func(format string, args ...interface{}) {
		t.Logf(format, args...)
	})
	expectedServer := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)

	app := fxtest.New(
		t,
		fx.Logger(logger),
		fx.Provide(
			func() *cobra.Command {
				return cmd
			},
			config.NewFactory().Configure,
			func(viper *viper.Viper) dependency.ConfigGetter {
				return viper
			},
			func() http.Handler {
				return handlers.MethodHandler{}
			},
			server.Service.Constructor,
		),
		server.Service.Dependencies,
		fx.Invoke(func(serverServer server.Server) {
			gotServer, ok := serverServer.(*http.Server)
			if !ok {
				t.Fatal("cannot convert server to http.Server")
			}
			if !reflect.DeepEqual(expectedServer, gotServer) {
				t.Fatalf("expected (%+v), got (%+v)", expectedServer, gotServer)
			}
			wg.Done()
		}),
	)
	go app.Run()
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
}

func TestStartServerFails(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	srv := stubServer{
		listenAndServe: func() error {
			wg.Done()
			return errors.New("an error")
		},
	}
	expectedMessages := []string{
		"Starting HTTP Server",
		"Could not start Server",
	}
	gotMessages := make([]string, 0, len(expectedMessages))
	wg.Add(len(expectedMessages))
	hookFunc := func(entry zapcore.Entry) error {
		gotMessages = append(gotMessages, entry.Message)
		wg.Done()
		return nil
	}
	hooks := zap.Hooks(hookFunc)
	logger := zaptest.NewLogger(
		t,
		zaptest.WrapOptions(hooks),
	)
	err := server.StartServer(srv, logger)(context.Background())
	if err != nil {
		t.Fatalf("expected error to be nil, got (%s)", err)
	}
	wg.Wait()
	if !reflect.DeepEqual(expectedMessages, gotMessages) {
		t.Fatalf("expected messages to be (%+v), got (%+v)", expectedMessages, gotMessages)
	}
}

func TestStopServerFails(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	srv := stubServer{
		shutdown: func(ctx context.Context) error {
			wg.Done()
			return errors.New("an error")
		},
	}
	expectedMessages := []string{
		"Stopping HTTP Server",
		"Error when shutting down Server",
	}
	gotMessages := make([]string, 0, len(expectedMessages))
	wg.Add(len(expectedMessages))
	hookFunc := func(entry zapcore.Entry) error {
		gotMessages = append(gotMessages, entry.Message)
		wg.Done()
		return nil
	}
	hooks := zap.Hooks(hookFunc)
	logger := zaptest.NewLogger(
		t,
		zaptest.WrapOptions(hooks),
	)
	err := server.StopServer(srv, logger)(context.Background())
	if err == nil {
		t.Fatal("expected error got, nil")
	}
	wg.Wait()
	if !reflect.DeepEqual(expectedMessages, gotMessages) {
		t.Fatalf("expected messages to be (%+v), got (%+v)", expectedMessages, gotMessages)
	}
}
