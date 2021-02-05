package response_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/BlackBX/service-framework/response"
)

func TestNewHTTPProblem(t *testing.T) {
	expectedProblem := &response.Problem{
		Status: http.StatusBadRequest,
		Type:   "https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
		Title:  http.StatusText(http.StatusBadRequest),
		Detail: "Hello, World!",
	}

	gotProblem := response.NewHTTPProblem(http.StatusBadRequest, "Hello, World!")
	if !reflect.DeepEqual(expectedProblem, gotProblem) {
		t.Fatalf("expected (%+v), got (%+v)", expectedProblem, gotProblem)
	}
}
