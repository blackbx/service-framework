package response

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// JSONEncoder is an interface that abstracts the encoding of JSON
type JSONEncoder interface {
	Encode(v interface{}) error
}

// NewJSONResponder creates a new instance of the JSONResponder type
// for the given request
func NewJSONResponder(logger *zap.Logger, rw http.ResponseWriter, r *http.Request) Responder {
	encoder := json.NewEncoder(rw)
	encoder.SetEscapeHTML(false)
	return JSONResponder{
		logger:         logger,
		responseWriter: rw,
		request:        r,
		Encoder:        encoder,
	}
}

// JSONResponder is a responder that will respond with JSON responses
type JSONResponder struct {
	logger         *zap.Logger
	responseWriter http.ResponseWriter
	request        *http.Request
	Encoder        JSONEncoder
}

// RespondWithProblem will respond with the given status code and
// detail, with an API problem
func (r JSONResponder) RespondWithProblem(statusCode int, detail string) {
	problem := NewHTTPProblem(statusCode, detail)
	r.responseWriter.WriteHeader(statusCode)
	if err := r.Encoder.Encode(problem); err != nil {
		r.logger.Error("Could not respond with problem", zap.Any("value", problem))
	}
}

// Respond will take a given struct and respond with it as the body
func (r JSONResponder) Respond(statusCode int, value interface{}) {
	r.responseWriter.WriteHeader(statusCode)
	if err := r.Encoder.Encode(value); err != nil {
		r.logger.Error("Could not respond with value", zap.Any("value", value))
	}
}

// RespondStream will stream a response of JSON values to the client
func (r JSONResponder) RespondStream(statusCode int, valueStream <-chan interface{}) {
	r.responseWriter.WriteHeader(statusCode)
	for value := range valueStream {
		if err := r.Encoder.Encode(value); err != nil {
			r.logger.Error("Could not respond with value stream", zap.Any("value", value))
		}
	}
}
