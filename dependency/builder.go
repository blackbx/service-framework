package dependency

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// ConfigRegisterer is a function type for dependencies to register
// the configuration that they require
type ConfigRegisterer func(set FlagSet)

// Service is a dependency that is required by multiple modules of the
// application, such as logging.
type Service struct {
	Name         string
	ConfigFunc   ConfigRegisterer
	Dependencies fx.Option
	Constructor  interface{}
	InvokeFunc   interface{}
}

// NewBuilder creates a new instance of the Builder type, by default it will provide
// the cobra.Command as a dependency.
func NewBuilder(cmd *cobra.Command) Builder {
	return Builder{
		Cmd: cmd,
		Provide: []interface{}{
			func() *cobra.Command {
				return cmd
			},
		},
		Invoke:  []fx.Option{},
		Options: []fx.Option{},
	}
}

// Builder is a type that will build an *fx.App from dependencies, it stores the command
// that started the app, along with all of the dependencies required to start the app.
type Builder struct {
	Cmd     *cobra.Command
	Provide []interface{}
	Invoke  []fx.Option
	Options []fx.Option
}

// WithConstructor is analogous to the fx.Provide option, it allows you to provide multiple
// constructors
func (b Builder) WithConstructor(constructors ...interface{}) Builder {
	b.Provide = append(b.Provide, constructors...)
	return b
}

// WithService allows you to register a Service dependency, this should be done within
// an `invoke()` function, so that the flags from the
func (b Builder) WithService(service Service) Builder {
	b = b.WithConstructor(service.Constructor)
	if service.Dependencies != nil {
		b.Options = append(b.Options, service.Dependencies)
	}
	if service.ConfigFunc != nil {
		service.ConfigFunc(b.Cmd.PersistentFlags())
	}
	if service.InvokeFunc == nil {
		return b
	}
	return b.WithInvoke(service.InvokeFunc)
}

// WithModule adds an fx.Option into the list of dependencies to be built, used for adding
// application modules
func (b Builder) WithModule(module fx.Option) Builder {
	b.Options = append(b.Options, module)
	return b
}

// WithInvoke will add a function as a function to be invoked
func (b Builder) WithInvoke(funcs ...interface{}) Builder {
	b.Invoke = append(b.Invoke, fx.Invoke(funcs...))
	return b
}

// Build will produce a new instance of the *fx.App from the variables of the builder
func (b Builder) Build() *fx.App {
	return fx.New(
		fx.Provide(b.Provide...),
		fx.Options(b.Invoke...),
		fx.Options(b.Options...),
	)
}

// BuildTest will create a new instance of *fxtest.App from the contained dependencies
func (b Builder) BuildTest(tb fxtest.TB) *fxtest.App {
	return fxtest.New(
		tb,
		fx.Provide(b.Provide...),
		fx.Options(b.Invoke...),
		fx.Options(b.Options...),
	)
}
