package dependency_test

import (
	"net/http"
	"testing"

	"github.com/BlackBX/service-framework/dependency"
	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func TestBuilder_WithConstructor(t *testing.T) {
	builder := dependency.NewBuilder(&cobra.Command{})
	builder = builder.WithConstructor(func() http.Handler {
		return handlers.MethodHandler{}
	})
	if len(builder.Provide) != 2 {
		t.Fatalf("expected there to be 2 things in the provide list, there is (%d)", len(builder.Provide))
	}
}

func TestBuilder_WithInvoke(t *testing.T) {
	builder := dependency.NewBuilder(&cobra.Command{})
	builder = builder.WithInvoke(func(command *cobra.Command) {

	})
	if len(builder.Invoke) != 1 {
		t.Fatalf("expected there be one thing in the invoke list, there is (%d)", len(builder.Invoke))
	}
}

func TestBuilder_WithModule(t *testing.T) {
	builder := dependency.NewBuilder(&cobra.Command{})
	builder = builder.WithModule(fx.Provide())
	if len(builder.Options) != 1 {
		t.Fatalf("expected there to be one Option, got (%d)", len(builder.Options))
	}
}

func TestBuilder_WithService(t *testing.T) {
	tests := []struct {
		name            string
		expectedOptions int
		expectedInvokes int
		service         dependency.Service
	}{
		{
			name:            "no invoke",
			expectedInvokes: 0,
			expectedOptions: 1,
			service: dependency.Service{
				Dependencies: fx.Options(),
				Constructor: func() string {
					const foo = "foo"
					return foo
				},
			},
		},
		{
			name:            "invoke",
			expectedInvokes: 1,
			expectedOptions: 1,
			service: dependency.Service{
				Dependencies: fx.Options(),
				InvokeFunc:   func(command *cobra.Command) {},
				Constructor: func() string {
					return "foo"
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := dependency.NewBuilder(&cobra.Command{})
			builder = builder.WithService(test.service)
			if len(builder.Options) != test.expectedOptions {
				t.Errorf("expected (%d), options, got (%d) options", test.expectedOptions, len(builder.Options))
			}
			if len(builder.Invoke) != test.expectedInvokes {
				t.Errorf("expected (%d) invokes, got (%d) invokes", test.expectedInvokes, len(builder.Invoke))
			}
			_ = builder.Build()
			_ = builder.BuildTest(t)
		})
	}
}

func TestBuilder_WithServiceCallsConfigFunc(t *testing.T) {
	timesCalled := 0
	service := dependency.Service{
		ConfigFunc: func(set dependency.FlagSet) {
			timesCalled++
		},
		Dependencies: fx.Provide(),
	}
	builder := dependency.NewBuilder(&cobra.Command{})
	builder = builder.WithService(service)
	if timesCalled != 1 {
		t.Fatalf("expected to be called 1 time, called (%d) times", timesCalled)
	}
}
