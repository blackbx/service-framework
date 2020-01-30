package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/BlackBX/service-framework/dependency"
	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq" // We need the database
	"go.uber.org/fx"
)

// Service defines the configuration and constructors required to get a postgres *sql.DB
var Service = dependency.Service{
	ConfigFunc: func(set dependency.FlagSet) {
		set.String("postgres-dbname", "postgres", "The name of the database to connect to")
		set.String("postgres-user", "postgres", "The name of the user to connect to the postgres db with")
		set.String("postgres-password", "", "The password to connect to the postgres database")
		set.String("postgres-host", "localhost", "The host to connect to the postgres database on")
		set.Int("postgres-port", 5432, "The port to connect to the postgres database on")
		set.String("postgres-sslmode", "disable", "What sslmode to use with the postgres database")
		set.String("postgres-fallback-application-name", "", "An application_name for postgres to fall back to if one isn't provided.")
		set.Duration("postgres-connect-timeout", 0, "Maximum wait for connection, 0 means wait indefinitely")
		set.String("postgres-sslcert", "", "Cert file location. The file must contain PEM encoded data.")
		set.String("postgres-sslkey", "", "Key file location. The file must contain PEM encoded data.")
		set.String("postgres-sslrootkey", "", "The location of the root certificate file. The file must contain PEM encoded data.")
	},
	Dependencies: fx.Provide(
		NewConfig, NewSQLX,
	),
	Constructor: NewFactory().DB,
}

// Config defines the configuration required to create a postgres database connection
type Config struct {
	DBName                  string
	User                    string
	Password                string
	Host                    string
	Port                    int
	SSLMode                 string
	FallbackApplicationName string
	ConnectTimeout          time.Duration
	SSLCert                 string
	SSLKey                  string
	SSLRootCert             string
}

func (c Config) sslStringParts() []string {
	parameters := make([]string, 0, 4)
	parameters = append(parameters, fmt.Sprintf("sslmode=%s", c.SSLMode))
	if c.SSLCert != "" {
		parameters = append(parameters, fmt.Sprintf("sslcert=%s", c.SSLCert))
	}
	if c.SSLKey != "" {
		parameters = append(parameters, fmt.Sprintf("sslkey=%s", c.SSLKey))
	}
	if c.SSLRootCert != "" {
		parameters = append(parameters, fmt.Sprintf("parameters=%s", c.SSLRootCert))
	}
	return parameters
}

// String implements fmt.Stringer, and returns the required configuration as a
// libpq connection string
func (c Config) String() string {
	parameters := make([]string, 0, 11)
	if c.DBName != "" {
		parameters = append(parameters, fmt.Sprintf("dbname=%s", c.DBName))
	}
	if c.User != "" {
		parameters = append(parameters, fmt.Sprintf("user=%s", c.User))
	}
	if c.Password != "" {
		parameters = append(parameters, fmt.Sprintf("password=%s", c.Password))
	}
	if c.Host != "" {
		parameters = append(parameters, fmt.Sprintf("host=%s", c.Host))
	}
	if c.Port != 0 {
		parameters = append(parameters, fmt.Sprintf("port=%d", c.Port))
	}
	if c.FallbackApplicationName != "" {
		parameters = append(parameters, fmt.Sprintf("fallback_application_name=%s", c.FallbackApplicationName))
	}
	if c.ConnectTimeout != 0 {
		parameters = append(parameters, fmt.Sprintf("connect_timeout=%d", c.ConnectTimeout/time.Second))
	}
	switch c.SSLMode {
	case "require", "verify-ca", "verify-full":
		parameters = append(parameters, c.sslStringParts()...)
	case "disable":
		parameters = append(parameters, "sslmode=disable")
	}
	return strings.Join(parameters, " ")
}

// NewConfig creates a new instance of the configuration from app configuration
func NewConfig(config dependency.ConfigGetter) Config {
	return Config{
		DBName:                  config.GetString("postgres-dbname"),
		User:                    config.GetString("postgres-user"),
		Password:                config.GetString("postgres-password"),
		Host:                    config.GetString("postgres-host"),
		Port:                    config.GetInt("postgres-port"),
		SSLMode:                 config.GetString("postgres-sslmode"),
		FallbackApplicationName: config.GetString("postgres-fallback-application-name"),
		ConnectTimeout:          config.GetDuration("postgres-connect-timeout"),
		SSLCert:                 config.GetString("postgres-sslcert"),
		SSLKey:                  config.GetString("postgres-sslkey"),
		SSLRootCert:             config.GetString("postgres-sslrootkey"),
	}
}

// NewFactory creates a new instance of a Factory that can create postgres *sql.DBs
func NewFactory() Factory {
	return Factory{Opener: sql.Open}
}

// Factory is a type that can create
type Factory struct {
	Opener func(driverName, dataSourceName string) (*sql.DB, error)
}

// DB creates a new instance of a postgres *sql.DB and registers it to the health check
func (f Factory) DB(config Config, check healthcheck.Handler) (*sql.DB, error) {
	connString := config.String()
	db, err := f.Opener("nrpostgres", connString)
	if err != nil {
		return nil, fmt.Errorf("could not create database, got error (%w)", err)
	}
	// nolint: gomnd
	check.AddReadinessCheck(config.Host, healthcheck.DatabasePingCheck(db, 10*time.Second))
	return db, nil
}

func NewSQLX(db *sql.DB) *sqlx.DB {
	return sqlx.NewDb(db, "postgres")
}
