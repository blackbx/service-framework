package config_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/BlackBX/service-framework/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type failingViper struct{}

func (f failingViper) AutomaticEnv() {}

func (f failingViper) SetEnvKeyReplacer(_ *strings.Replacer) {}

func (f failingViper) BindPFlags(flags *pflag.FlagSet) error {
	return errors.New("an error")
}

func TestFactory_ConfigureSucceeds(t *testing.T) {
	_, err := config.NewFactory().Configure(&cobra.Command{})
	if err != nil {
		t.Fatalf("expected error to be nil, got (%s)", err)
	}
}

func TestFactory_ConfigureFails(t *testing.T) {
	_, err := config.Factory{
		Replacer: strings.NewReplacer("-", "_",
			".", "_"),
		ConfigFunc: func(config config.Viper, cmd *cobra.Command, replacer *strings.Replacer) error {
			return errors.New("an error")
		},
	}.Configure(&cobra.Command{})
	if err == nil {
		t.Fatal("expected an error, got none")
	}
}

func TestConfigureViperFails(t *testing.T) {
	if err := config.ConfigureViper(failingViper{}, &cobra.Command{}, strings.NewReplacer("-", "_",
		".", "_")); err == nil {
		t.Fatal("expected an error, got none")
	}
}
