package redis

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/BlackBX/service-framework/dependency"
	"github.com/go-redis/redis/v7"
	"github.com/heptiolabs/healthcheck"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v7"
	"go.uber.org/fx"
)

// Cmdable is an interface to abstract redis
type Cmdable interface {
	redis.Cmdable
	WithContext(ctx context.Context) *redis.Client
}

// Service is the service to be used by the framework.Builder
var Service = dependency.Service{
	// nolint: gomnd
	ConfigFunc: func(set dependency.FlagSet) {
		set.String(
			"redis-network",
			"tcp",
			"Network type to connect to redis with (unix/tcp)",
		)
		set.String(
			"redis-host",
			"localhost",
			"The host to connect to redis on",
		)
		set.Int(
			"redis-port",
			6379,
			"The port to connect to redis on",
		)
		set.String(
			"redis-password",
			"",
			"The password to connect to redis with",
		)
		set.Int(
			"redis-db",
			0,
			"The database to use on the redis connection",
		)
		set.Int(
			"redis-max-retries",
			0,
			"Maximum number of times to retry connecting to redis",
		)
		set.Duration(
			"redis-min-retry-backoff",
			8*time.Millisecond,
			"The minimum backoff time from a retry",
		)
		set.Duration(
			"redis-max-retry-backoff",
			512*time.Millisecond,
			"The maximum backoff time from a retry",
		)
		set.Duration(
			"redis-dial-timeout",
			5*time.Second,
			"The time to wait for a connection",
		)
		set.Duration(
			"redis-read-timeout",
			3*time.Second,
			"The time to wait for a read",
		)
		set.Duration(
			"redis-write-timeout",
			3*time.Second,
			"The time to wait on write",
		)
		set.Int(
			"redis-pool-size",
			10*runtime.NumCPU(),
			"Size of the connection pool",
		)
		set.Int(
			"redis-min-idle-conns",
			0,
			"The minimum number of idle connections",
		)
		set.Duration(
			"redis-max-conn-age",
			0,
			"The max age of a redis pool connection",
		)
		set.Duration(
			"redis-pool-timeout",
			4*time.Second,
			"Amount of time client waits for connection if all connections are busy before returning an error.",
		)
		set.Duration(
			"redis-idle-check-frequency",
			time.Minute,
			"The frequency that connections are checked to be dropped",
		)
	},
	Dependencies: fx.Provide(
		NewOptions,
		NewClientConstructor,
	),
	Constructor: New,
}

// NewClientConstructor is a type that can create an instance of a redis.Cmdable
type ClientConstructor func() Cmdable

// NewClientConstructor is a type that can create a new instance of a
// ClientConstructor from a *redis.Options
func NewClientConstructor(options *redis.Options) ClientConstructor {
	return func() Cmdable {
		client := redis.NewClient(options)
		client.AddHook(nrredis.NewHook(options))
		return client
	}
}

// New creates a new instance of a *redis.Client
func New(constructor ClientConstructor, options *redis.Options, checker healthcheck.Handler) Cmdable {
	client := constructor()
	checker.AddReadinessCheck(options.Addr, func() error {
		_, err := client.Ping().Result()
		if err != nil {
			return fmt.Errorf("could not ping redis, got (%w)", err)
		}
		return nil
	})
	return client
}

// NewOptions creates a new *redis.Options struct from configuration
func NewOptions(config dependency.ConfigGetter) *redis.Options {
	addr := fmt.Sprintf(
		"%s:%s",
		config.GetString("redis-host"),
		config.GetString("redis-port"),
	)
	return &redis.Options{
		Network:            config.GetString("redis-network"),
		Addr:               addr,
		Password:           config.GetString("redis-password"),
		DB:                 config.GetInt("redis-db"),
		MaxRetries:         config.GetInt("redis-max-retries"),
		MinRetryBackoff:    config.GetDuration("redis-min-retry-backoff"),
		MaxRetryBackoff:    config.GetDuration("redis-max-retry-backoff"),
		DialTimeout:        config.GetDuration("redis-dial-timeout"),
		ReadTimeout:        config.GetDuration("redis-read-timeout"),
		WriteTimeout:       config.GetDuration("redis-write-timeout"),
		PoolSize:           config.GetInt("redis-pool-size"),
		MinIdleConns:       config.GetInt("redis-min-idle-conns"),
		MaxConnAge:         config.GetDuration("redis-max-conn-age"),
		PoolTimeout:        config.GetDuration("redis-pool-timeout"),
		IdleTimeout:        config.GetDuration("redis-idle-timeout"),
		IdleCheckFrequency: config.GetDuration("redis-idle-check-frequency"),
	}
}
