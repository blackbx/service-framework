package response

import (
	"net/http"

	"github.com/BlackBX/service-framework/dependency"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Service is the definition of the dependency
var Service = dependency.Service{
	Dependencies: fx.Provide(
		func() ResponderConstructor {
			return NewJSONResponder
		},
		NewFactory,
	),
	Constructor: func(factory ResponderFactory) ResponderProvider {
		return factory
	},
}

// ResponderConstructor is a function that can create a new instance of
// a responder
type ResponderConstructor func(logger *zap.Logger, rw http.ResponseWriter, r *http.Request) Responder

// NewFactory creates a new instance of the ResponderFactory
func NewFactory(logger *zap.Logger, defaultResponder ResponderConstructor) ResponderFactory {
	return ResponderFactory{
		Logger:           logger,
		DefaultResponder: defaultResponder,
	}
}

// ResponderProvider is an interface that abstracts the providing of Responders
type ResponderProvider interface {
	Responder(rw http.ResponseWriter, r *http.Request) Responder
}

// ResponderFactory is a factory that can create new Responders, it allows
// for a responder to be created in a handler and subsequently called.
type ResponderFactory struct {
	Logger           *zap.Logger
	DefaultResponder ResponderConstructor
}

// Responder creates a new instance of a responder
func (rf ResponderFactory) Responder(rw http.ResponseWriter, r *http.Request) Responder {
	return rf.DefaultResponder(rf.Logger, rw, r)
}

// Responder is an interface that abstracts the production of the
// response away from handlers
type Responder interface {
	RespondWithProblem(statusCode int, detail string)
	Respond(statusCode int, value interface{})
	RespondStream(statusCode int, valueStream <-chan interface{})
}
