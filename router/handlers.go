package router

import (
	"net/http"

	"github.com/BlackBX/service-framework/response"
)

// New404Handler produces an http.Handler that is used as the default 404 handler
func New404Handler(provider response.ResponderProvider) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		provider.
			Responder(rw, r).
			RespondWithProblem(http.StatusNotFound, "ROUTE_NOT_FOUND")
	})
}

// New405Handler produces an http.Handler that is used as the default 405 handler
func New405Handler(provider response.ResponderProvider) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		provider.
			Responder(rw, r).
			RespondWithProblem(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
	})
}
