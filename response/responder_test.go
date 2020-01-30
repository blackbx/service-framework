package response_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/BlackBX/service-framework/config"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/logging"
	"github.com/BlackBX/service-framework/response"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap/zaptest"
)

type PrinterFunc func(string, ...interface{})

func (p PrinterFunc) Printf(string, ...interface{}) {}

func TestNewResponderFactory(t *testing.T) {
	fxlogger := fx.Logger(PrinterFunc(func(string, ...interface{}) {}))
	command := &cobra.Command{}
	command.SetOut(ioutil.Discard)
	timesCalled := 0
	dependencyBuilder := dependency.NewBuilder(command).
		WithService(config.Service).
		WithService(logging.Service).
		WithService(response.Service).
		WithModule(fxlogger)

	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	dependencyBuilder = dependencyBuilder.
		WithInvoke(func(factory response.ResponderProvider) {
			timesCalled++
		})
	app := dependencyBuilder.BuildTest(t)
	go app.Run()
	if err := app.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if timesCalled != 1 {
		t.Fatalf("expected times called to be 1, got (%d)", timesCalled)
	}
}

func TestResponderFactory_Responder(t *testing.T) {
	logger := zaptest.NewLogger(t)
	factory := response.NewFactory(logger, response.NewJSONResponder)
	request := httptest.NewRequest(http.MethodGet,
		"https://example.com",
		strings.NewReader("Hello, World!"))
	resp := httptest.NewRecorder()

	factory.
		Responder(resp, request).
		RespondWithProblem(http.StatusBadRequest, "Hello, World!")

	expectedProblem := response.NewHTTPProblem(http.StatusBadRequest, "Hello, World!")
	gotProblem := response.Problem{}
	if err := json.NewDecoder(resp.Body).Decode(&gotProblem); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if !reflect.DeepEqual(expectedProblem, gotProblem) {
		t.Fatalf("expected (%+v), got (%+v)", expectedProblem, gotProblem)
	}
}
