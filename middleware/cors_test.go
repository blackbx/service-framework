package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/BlackBX/service-framework/config"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/middleware"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type (
	CORSParams struct {
		fx.In
		Middleware []mux.MiddlewareFunc `group:"middleware"`
	}
)

func TestNewCORSGet(t *testing.T) {
	cmd := &cobra.Command{}
	wg := new(sync.WaitGroup)
	wg.Add(1)

	expectedHeaders := http.Header{
		"Access-Control-Allow-Credentials": []string{"true"},
		"Access-Control-Allow-Origin":      []string{"localhost"},
	}

	builder := dependency.
		NewBuilder(cmd).
		WithService(config.Service).
		WithService(middleware.CORSService).
		WithInvoke(func(params CORSParams) {
			if len(params.Middleware) != 1 {
				t.Error("Expected 1 and got 0")
			}
			request := httptest.NewRequest(http.MethodGet, "http://stampede.ai", http.NoBody)
			request.Header.Add("Origin", "localhost")
			response := httptest.NewRecorder()
			handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write([]byte("Hello"))
			})
			params.Middleware[0](handler).ServeHTTP(response, request)
			if !reflect.DeepEqual(expectedHeaders, response.Header()) {
				t.Errorf("Expected %+v got %+v", expectedHeaders, response.Header())
			}
			wg.Done()
		})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	app := builder.BuildTest(t)
	go app.Run()
	wg.Wait()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
	if err := app.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	cancel()
}

func TestNewCORSPreflight(t *testing.T) {
	cmd := &cobra.Command{}
	wg := new(sync.WaitGroup)
	wg.Add(1)

	expectedHeaders := http.Header{
		"Access-Control-Allow-Credentials": []string{"true"},
		"Access-Control-Allow-Headers":     []string{"X-Requested-With"},
		"Access-Control-Allow-Origin":      []string{"localhost"},
	}

	builder := dependency.
		NewBuilder(cmd).
		WithService(config.Service).
		WithService(middleware.CORSService).
		WithInvoke(func(params CORSParams) {
			if len(params.Middleware) != 1 {
				t.Error("Expected 1 and got 0")
			}
			request := httptest.NewRequest(http.MethodOptions, "http://localhost", http.NoBody)
			request.Header.Add("Access-Control-Request-Method", "GET")
			request.Header.Add("Access-Control-Request-Headers", "origin, x-requested-with")

			request.Header.Add("Origin", "localhost")
			response := httptest.NewRecorder()
			handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write([]byte("Hello"))
				t.Error("Should not need to send anything back")
			})
			params.Middleware[0](handler).ServeHTTP(response, request)
			if !reflect.DeepEqual(expectedHeaders, response.Header()) {
				t.Errorf("Expected %+v got %+v", expectedHeaders, response.Header())
			}
			wg.Done()
		})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	app := builder.BuildTest(t)
	go app.Run()
	wg.Wait()
	if err := app.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
	if err := app.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
	cancel()
}
