package health_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/BlackBX/service-framework/health"
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
)

func TestRegisterHealthcheck(t *testing.T) {
	router := mux.NewRouter()
	health.RegisterHealthcheck(healthcheck.NewHandler()).Router(router)
	expectedRoutes := []string{"/live", "/ready"}
	gotRotues := make([]string, 0, len(expectedRoutes))
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			t.Fatal(err)
		}
		gotRotues = append(gotRotues, path)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(expectedRoutes)
	sort.Strings(gotRotues)
	if !reflect.DeepEqual(expectedRoutes, gotRotues) {
		t.Fatalf("expected (%+v), got (%+v)", expectedRoutes, gotRotues)
	}
}
