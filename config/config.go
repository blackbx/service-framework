package config

import (
	"fmt"
	"strings"

	"github.com/BlackBX/service-framework/dependency"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Service provides the config framework as a *viper.Viper and a
// framework.ConfigGetter
var Service = dependency.Service{
	Dependencies: fx.Provide(
		NewFactory().Configure,
	),
	Constructor: func(config *viper.Viper) dependency.ConfigGetter {
		return config
	},
}

// Viper is an interface that the *viper.Viper type adheres to, this is
// to enable the package to be thoroughly test
type Viper interface {
	AutomaticEnv()
	SetEnvKeyReplacer(r *strings.Replacer)
	BindPFlags(flags *pflag.FlagSet) error
}

// NewFactory gives a new instance of the Factory type with a string.Replacer
// that will replace - or . with _ when in an environment variable, it also
// sets the ConfigFunc to be ConfigureViper function
func NewFactory() Factory {
	return Factory{
		Replacer: strings.NewReplacer("-", "_",
			".", "_"),
		ConfigFunc: ConfigureViper,
	}
}

// Factory is a type that can produce new instances of the *viper.Viper type
type Factory struct {
	Replacer   *strings.Replacer
	ConfigFunc func(config Viper, cmd *cobra.Command, replacer *strings.Replacer) error
}

// Configure will produce a *viper.Viper type configured by the ConfigFunc
// held by the Factory
func (f Factory) Configure(cmd *cobra.Command) (*viper.Viper, error) {
	config := viper.New()
	if err := f.ConfigFunc(config, cmd, f.Replacer); err != nil {
		return nil, err
	}
	return config, nil
}

// ConfigureViper is a function that will configure viper
func ConfigureViper(config Viper, cmd *cobra.Command, replacer *strings.Replacer) error {
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(replacer)
	if err := config.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("failed to bind command flags with error (%w)", err)
	}
	return nil
}
