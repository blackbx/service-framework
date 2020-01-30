package newrelic

import (
	"os"
	"path/filepath"

	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/httpclient"
	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Service allows newrelic to be added to an application, it adds the middleware aswell too
var Service = dependency.Service{
	Dependencies: fx.Provide(
		fx.Annotated{
			Group: "trippers",
			Target: func() httpclient.Tripper {
				return newrelic.NewRoundTripper
			},
		},
	),
	ConfigFunc: func(set dependency.FlagSet) {
		set.String("newrelic-app-name", filepath.Base(os.Args[0]), "The name of the application")
		set.String("newrelic-license-key", "", "Newrelic license key")
		set.Bool("newrelic-distributed-tracer-enabled", true, "Whether to add and read distributed tracing headers")
	},
	Constructor: NewApp,
}

// NewApp will create a new instance of the *newrelic.Application
func NewApp(config dependency.ConfigGetter, logger *zap.Logger) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(config.GetString("newrelic-app-name")),
		newrelic.ConfigLicense(config.GetString("newrelic-license-key")),
		newrelic.ConfigDistributedTracerEnabled(config.GetBool("newrelic-distributed-tracer-enabled")),
		nrzap.ConfigLogger(logger),
	)
}
