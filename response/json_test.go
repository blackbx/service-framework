package response_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BlackBX/service-framework/response"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

type failingEncoder struct{}

func (f failingEncoder) Encode(v interface{}) error {
	return errors.New("failure")
}

func TestJSONResponderEncoderFails(t *testing.T) {
	tests := []struct {
		name              string
		expectedLogString string
		testCase          func(responder response.Responder)
	}{
		{
			name:              "problem",
			expectedLogString: "Could not respond with problem",
			testCase: func(responder response.Responder) {
				responder.RespondWithProblem(http.StatusBadRequest, "Hello, World!")
			},
		},
		{
			name:              "respond",
			expectedLogString: "Could not respond with value",
			testCase: func(responder response.Responder) {
				body := response.NewHTTPProblem(http.StatusBadRequest, "Hello, World!")
				responder.Respond(http.StatusBadRequest, body)
			},
		},
		{
			name:              "respond stream",
			expectedLogString: "Could not respond with value stream",
			testCase: func(responder response.Responder) {
				stream := make(chan interface{}, 1)
				body := response.NewHTTPProblem(http.StatusBadRequest, "Hello, World!")
				stream <- body
				close(stream)
				responder.RespondStream(http.StatusBadRequest, stream)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hookfunc, callCount := newResponseEncoderHookFunc(t, test.expectedLogString)
			options := zaptest.WrapOptions(zap.Hooks(hookfunc))
			logger := zaptest.NewLogger(t, options)
			responder := response.NewJSONResponder(
				logger,
				httptest.NewRecorder(),
				httptest.NewRequest(http.MethodGet, "https://example.com", strings.NewReader("")),
			)
			jsonResponder := responder.(response.JSONResponder)
			jsonResponder.Encoder = failingEncoder{}
			test.testCase(jsonResponder)
			if *callCount != 1 {
				t.Fatalf("expected log to be called 1 time, it was called (%d) time(s)", *callCount)
			}
		})
	}
}

type HookFunc func(entry zapcore.Entry) error

func newResponseEncoderHookFunc(t *testing.T, expectedLog string) (hookFunc HookFunc, callCount *int) {
	count := 0
	return func(entry zapcore.Entry) error {
		if entry.Message != expectedLog {
			t.Fatalf("expected (%s), got (%s)", expectedLog, entry.Message)
		}
		count++
		return nil
	}, &count
}
