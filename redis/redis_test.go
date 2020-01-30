package redis_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/BlackBX/service-framework/config"
	redisconfig "github.com/BlackBX/service-framework/redis"
	"github.com/go-redis/redis/v7"
	"github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"
)

func parseDuration(t *testing.T, durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		t.Fatal(err)
	}
	return duration
}

func TestNewOptions(t *testing.T) {
	cmd := &cobra.Command{}
	redisconfig.Service.ConfigFunc(cmd.PersistentFlags())
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.NewFactory().Configure(cmd)
	if err != nil {
		t.Fatal(err)
	}
	expectedOptions := &redis.Options{
		Network:            "tcp",
		Addr:               "localhost:6379",
		Password:           "",
		DB:                 0,
		MaxRetries:         0,
		MinRetryBackoff:    parseDuration(t, "8ms"),
		MaxRetryBackoff:    parseDuration(t, "512ms"),
		DialTimeout:        parseDuration(t, "5s"),
		ReadTimeout:        parseDuration(t, "3s"),
		WriteTimeout:       parseDuration(t, "3s"),
		PoolSize:           10 * runtime.NumCPU(),
		MinIdleConns:       0,
		MaxConnAge:         parseDuration(t, "0s"),
		PoolTimeout:        parseDuration(t, "4s"),
		IdleTimeout:        parseDuration(t, "0s"),
		IdleCheckFrequency: parseDuration(t, "1m0s"),
	}

	gotOptions := redisconfig.NewOptions(cfg)
	if !reflect.DeepEqual(expectedOptions, gotOptions) {
		t.Fatalf("expected options (%+v), got (%+v)", expectedOptions, gotOptions)
	}
}

type stubHealthchecker struct {
	addLivenessCheck  func(name string, check healthcheck.Check)
	addReadinessCheck func(name string, check healthcheck.Check)
}

func (s stubHealthchecker) ServeHTTP(http.ResponseWriter, *http.Request) {}

func (s stubHealthchecker) AddLivenessCheck(name string, check healthcheck.Check) {
	s.addLivenessCheck(name, check)
}

func (s stubHealthchecker) AddReadinessCheck(name string, check healthcheck.Check) {
	s.addReadinessCheck(name, check)
}

func (s stubHealthchecker) LiveEndpoint(http.ResponseWriter, *http.Request) {}

func (s stubHealthchecker) ReadyEndpoint(http.ResponseWriter, *http.Request) {}

type stubCmdable struct {
	redis.Cmdable
	cmdFunc func() *redis.StatusCmd
}

func (s stubCmdable) Ping() *redis.StatusCmd {
	return s.cmdFunc()
}

func TestHealthcheckFails(t *testing.T) {
	cmd := stubCmdable{
		cmdFunc: func() *redis.StatusCmd {
			return redis.NewStatusResult("", errors.New("an error"))
		},
	}
	constuctor := func() redisconfig.Cmdable {
		return cmd
	}
	healthChecker := stubHealthchecker{
		addReadinessCheck: func(name string, check healthcheck.Check) {
			expectedName := ""
			if name != expectedName {
				t.Fatalf("expected name to be (%s), got (%s)", expectedName, name)
			}
			err := check()
			if err == nil {
				t.Fatal("expected error to not be nil, got nil")
			}
		},
	}
	_ = redisconfig.New(constuctor, &redis.Options{}, healthChecker)
}

func TestHealthcheckSucceeds(t *testing.T) {
	cmd := stubCmdable{
		cmdFunc: func() *redis.StatusCmd {
			return redis.NewStatusResult("PONG", nil)
		},
	}
	constuctor := func() redisconfig.Cmdable {
		return cmd
	}
	healthChecker := stubHealthchecker{
		addReadinessCheck: func(name string, check healthcheck.Check) {
			expectedName := ""
			if name != expectedName {
				t.Fatalf("expected name to be (%s), got (%s)", expectedName, name)
			}
			err := check()
			if err != nil {
				t.Fatalf("expected error to be nil, got (%s)", err)
			}
		},
	}
	_ = redisconfig.New(constuctor, &redis.Options{}, healthChecker)
}

func (s stubCmdable) WithContext(ctx context.Context) *redis.Client {
	return new(redis.Client)
}

func TestNewClientConstructor(t *testing.T) {
	constructor := redisconfig.NewClientConstructor(&redis.Options{})
	client := constructor()
	_, ok := client.(*redis.Client)
	if !ok {
		t.Fatal("expected constructor to be a *redis.Client")
	}
}
