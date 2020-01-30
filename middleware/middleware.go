package middleware

import (
	"github.com/BlackBX/service-framework/logging"
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"go.uber.org/fx"
)

// Module allows the default middlewares to be registered to an app
var Module = fx.Provide(
	fx.Annotated{
		Group: "middleware",
		Target: func(logger logging.PrintLogger) mux.MiddlewareFunc {
			return handlers.RecoveryHandler(
				handlers.RecoveryLogger(logger),
			)
		},
	},
	fx.Annotated{
		Group: "middleware",
		Target: func() mux.MiddlewareFunc {
			return gziphandler.GzipHandler
		},
	},
	fx.Annotated{
		Group:  "middleware",
		Target: nrgorilla.Middleware,
	},
)
