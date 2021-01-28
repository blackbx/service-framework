package test

import (
	"net/http"

	"github.com/BlackBX/service-framework/router"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// Module allows the routes from this module to be registered to the app
var Module = fx.Provide(
	NewRedisHandler,
	NewPGHandler,
	NewHTTPHandler,
	fx.Annotated{
		Group:  "server",
		Target: RegisterHandler,
	},
)

// HandlerParams is the type that defines the parameters that are
// required to register the handlers to the router
type HandlerParams struct {
	fx.In

	RedisHandler RedisHandler
	PGHandler    PGHandler
	HTTPHandler  HTTPHandler
}

// RegisterHandler registers the handlers to the router
func RegisterHandler(params HandlerParams) router.Module {
	return router.Module{
		Path: "test",
		Router: func(router *mux.Router) {
			router.Handle("/redis", handlers.MethodHandler{
				http.MethodGet: http.HandlerFunc(params.RedisHandler.Get),
			})
			router.Handle("/pg", handlers.MethodHandler{
				http.MethodGet: http.HandlerFunc(params.PGHandler.Get),
			})
			router.Handle("/http/{id}", handlers.MethodHandler{
				http.MethodGet: http.HandlerFunc(params.HTTPHandler.Get),
			})
		},
	}
}
