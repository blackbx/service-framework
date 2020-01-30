package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/response"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// Service is how the dependency is provided to the dependency builder
var Service = dependency.Service{
	Dependencies: fx.Provide(
		func(router *mux.Router) http.Handler {
			return router
		},
	),
	Constructor: New,
}

// ApplierFunc is a function type that allows routes to be applied to
// the main router
type ApplierFunc func(router *mux.Router)

// Module is a group of routes to route to based on a path
type Module struct {
	Path   string
	Router ApplierFunc
}

// PathPrefix returns the path with a slash at the start
func (m Module) PathPrefix() string {
	if strings.HasPrefix(m.Path, "/") {
		return m.Path
	}
	return fmt.Sprintf("/%s", m.Path)
}

// Params are the parameters required to build the router
type Params struct {
	fx.In

	ResponseProvider response.ResponderProvider
	Modules          []Module             `group:"server"`
	Middlewares      []mux.MiddlewareFunc `group:"middleware"`
}

// New creates a new instance of a *mux.Router with all of the modules added
func New(params Params) *mux.Router {
	router := mux.NewRouter()
	for _, module := range params.Modules {
		subRouter := router.PathPrefix(module.PathPrefix()).Subrouter()
		module.Router(subRouter)
	}
	router.Use(params.Middlewares...)
	router.NotFoundHandler = New404Handler(params.ResponseProvider)
	router.MethodNotAllowedHandler = New405Handler(params.ResponseProvider)
	return router
}
