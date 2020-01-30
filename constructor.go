package framework

import (
	"github.com/BlackBX/service-framework/config"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/health"
	"github.com/BlackBX/service-framework/httpclient"
	"github.com/BlackBX/service-framework/logging"
	"github.com/BlackBX/service-framework/newrelic"
	"github.com/BlackBX/service-framework/postgres"
	"github.com/BlackBX/service-framework/redis"
	"github.com/BlackBX/service-framework/response"
	"github.com/BlackBX/service-framework/router"
	"github.com/BlackBX/service-framework/server"
	"github.com/spf13/cobra"
)

// NewWebApplicationBuilder will give you a builder that can
// create a new web application
func NewWebApplicationBuilder(command *cobra.Command) dependency.Builder {
	return dependency.
		NewBuilder(command).
		WithService(postgres.Service).
		WithService(newrelic.Service).
		WithService(config.Service).
		WithService(logging.Service).
		WithService(health.Service).
		WithService(router.Service).
		WithService(response.Service).
		WithService(redis.Service).
		WithService(httpclient.Service).
		WithService(server.Service)
}
