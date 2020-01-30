package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BlackBX/service-framework/response"
	"github.com/BlackBX/service-framework/router"
)

type stubResponderProvider struct {
	responder stubResponder
}

func (s stubResponderProvider) Responder(rw http.ResponseWriter, r *http.Request) response.Responder {
	return s.responder
}

type stubResponder struct {
	problem       func(statusCode int, detail string)
	respond       func(statusCode int, value interface{})
	respondStream func(statusCode int, valueStream <-chan interface{})
}

func (s stubResponder) RespondWithProblem(statusCode int, detail string) {
	s.problem(statusCode, detail)
}

func (s stubResponder) Respond(statusCode int, value interface{}) {
	s.respond(statusCode, value)
}

func (s stubResponder) RespondStream(statusCode int, valueStream <-chan interface{}) {
	s.respondStream(statusCode, valueStream)
}

// nolint: dupl
func TestNew404Handler(t *testing.T) {
	timesCalled := 0
	responderProvider := stubResponderProvider{
		responder: stubResponder{
			problem: func(statusCode int, detail string) {
				timesCalled++
				if statusCode != http.StatusNotFound {
					t.Fatalf("expected status code to be (%d), got (%d)", http.StatusNotFound, statusCode)
				}
				expectedDetail := "ROUTE_NOT_FOUND"
				if detail != expectedDetail {
					t.Fatalf("expected detail to be (%s), got (%s)", expectedDetail, detail)
				}
			},
		},
	}
	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com",
		http.NoBody,
	)
	router.
		New404Handler(responderProvider).
		ServeHTTP(httptest.NewRecorder(), req)
	if timesCalled != 1 {
		t.Fatalf("expected to be called 1 time, called (%d) time(s)", timesCalled)
	}
}

// nolint: dupl
func TestNew405Handler(t *testing.T) {
	timesCalled := 0
	responderProvider := stubResponderProvider{
		responder: stubResponder{
			problem: func(statusCode int, detail string) {
				timesCalled++
				if statusCode != http.StatusMethodNotAllowed {
					t.Fatalf("expected status code to be (%d), got (%d)", http.StatusNotFound, statusCode)
				}
				expectedDetail := "METHOD_NOT_ALLOWED"
				if detail != expectedDetail {
					t.Fatalf("expected detail to be (%s), got (%s)", expectedDetail, detail)
				}
			},
		},
	}
	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com",
		http.NoBody,
	)
	router.
		New405Handler(responderProvider).
		ServeHTTP(httptest.NewRecorder(), req)
	if timesCalled != 1 {
		t.Fatalf("expected to be called 1 time, called (%d) time(s)", timesCalled)
	}
}
