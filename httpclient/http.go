package httpclient

import (
	"net/http"

	"github.com/BlackBX/service-framework/dependency"
	"go.uber.org/fx"
)

// Service adds the ability to use the *http.Client type to the dependency injection container
var Service = dependency.Service{
	Name:        "httpclient",
	Constructor: New,
}

// Tripper is a type that can wrap an http.RoundTripper
type Tripper func(tripper http.RoundTripper) http.RoundTripper

// Params contains all of the required parameters to create an
// *http.Client, allowing you to wrap the transport
type Params struct {
	fx.In

	Trippers []Tripper `group:"trippers"`
}

// New creates a new instance of an *http.Client
func New(params Params) *http.Client {
	client := http.DefaultClient
	for _, tripper := range params.Trippers {
		client.Transport = tripper(client.Transport)
	}
	return client
}
