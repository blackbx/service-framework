package postgres_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/BlackBX/service-framework/config"
	"github.com/BlackBX/service-framework/postgres"
	"github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"
)

func TestConfig_String(t *testing.T) {
	tests := []struct {
		name           string
		config         postgres.Config
		expectedString string
	}{
		{
			name: "base",
			config: postgres.Config{
				DBName:                  "postgres",
				User:                    "root",
				Password:                "hunter2",
				Host:                    "localhost",
				Port:                    5432,
				SSLMode:                 "disable",
				FallbackApplicationName: "an-application",
			},
			expectedString: "dbname=postgres user=root password=hunter2 host=localhost port=5432 fallback_application_name=an-application",
		},
		{
			name: "SSL",
			config: postgres.Config{
				DBName:                  "postgres",
				User:                    "root",
				Password:                "hunter2",
				Host:                    "localhost",
				Port:                    5432,
				SSLMode:                 "require",
				FallbackApplicationName: "an-application",
				ConnectTimeout:          20 * time.Second,
				SSLCert:                 "foo.pem",
				SSLKey:                  "bar.pem",
				SSLRootCert:             "baz.pem",
			},
			// nolint: lll
			expectedString: "dbname=postgres user=root password=hunter2 host=localhost port=5432 fallback_application_name=an-application connect_timeout=20 sslmode=require sslcert=foo.pem sslkey=bar.pem parameters=baz.pem",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotString := test.config.String()
			if test.expectedString != gotString {
				t.Logf("expected string to be (%s), got (%s)", test.expectedString, gotString)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	cmd := &cobra.Command{}
	postgres.Service.ConfigFunc(cmd.PersistentFlags())
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.NewFactory().Configure(cmd)
	if err != nil {
		t.Fatal(err)
	}
	expectedString := "dbname=postgres user=root host=localhost port=5432 sslmode=disable"
	pgCfg := postgres.NewConfig(cfg)
	if pgCfg.String() != expectedString {
		t.Fatalf("expected string to be (%s), got (%s)", expectedString, pgCfg)
	}
}

func TestNew(t *testing.T) {
	cmd := &cobra.Command{}
	postgres.Service.ConfigFunc(cmd.PersistentFlags())
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	cfg, err := config.NewFactory().Configure(cmd)
	if err != nil {
		t.Fatal(err)
	}
	_, err = postgres.NewFactory().DB(postgres.NewConfig(cfg), healthcheck.NewHandler())
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewFails(t *testing.T) {
	cfg := postgres.Config{
		DBName:                  "postgres",
		User:                    "root",
		Host:                    "localhost",
		Port:                    5432,
		SSLMode:                 "require",
		FallbackApplicationName: "an-application",
		ConnectTimeout:          20 * time.Second,
		SSLCert:                 "testdata/foo.pem",
		SSLKey:                  "testdata/bar.pem",
		SSLRootCert:             "testdata/baz.pem",
	}
	factory := postgres.NewFactory()
	factory.Opener = func(driverName, dataSourceName string) (db *sql.DB, err error) {
		return nil, errors.New("an error")
	}
	_, err := factory.DB(cfg, healthcheck.NewHandler())
	if err == nil {
		t.Fatal("expected error to not be nil, got no error")
	}
}
