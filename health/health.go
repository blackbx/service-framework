package health

import (
	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/router"
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"go.uber.org/fx"
)

// Service allows the Health service to be registered with an application
var Service = dependency.Service{
	Dependencies: fx.Provide(
		healthcheck.NewHandler,
	),
	Constructor: fx.Annotated{
		Group:  "server",
		Target: RegisterHealthcheck,
	},
}

// RegisterHealthcheck registers the Healthcheck module with the router
func RegisterHealthcheck(check healthcheck.Handler) router.Module {
	return router.Module{
		Path: "health",
		Router: func(router *mux.Router) {
			router.HandleFunc("/live", check.LiveEndpoint)
			router.HandleFunc("/ready", check.ReadyEndpoint)
		},
	}
}
