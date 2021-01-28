package serve

import (
	framework "github.com/BlackBX/service-framework"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/example/server/internal/test"
	"github.com/BlackBX/service-framework/middleware"
	"github.com/spf13/cobra"
)

// NewCommand creates an instance of the Serve command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the web-server",
		Long:  "Start the example web-sever",
	}
	cmd.Run = Serve(framework.NewWebApplicationBuilder(cmd))
	return cmd
}

// Serve produces the function that is called when the
// command is called
func Serve(builder dependency.Builder) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		builder.
			WithModule(middleware.Module).
			WithModule(test.Module).
			Build().
			Run()
	}
}
