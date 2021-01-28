package read

import (
	framework "github.com/BlackBX/service-framework"
	"github.com/BlackBX/service-framework/dependency"
	"github.com/BlackBX/service-framework/example/queue/cmd/internal/reader"
	"github.com/spf13/cobra"
)

// NewCommand creates a new instance of the *cobra.Command
// that will start the queue reader
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Start the queue reader",
		Long:  "Start the queue reader reading from the SQS Queue",
	}
	cmd.Run = Read(framework.NewQueueApplicationBuilder(cmd))
	return cmd
}

// Read adds the module to the application, along with its invoke
// function to start the server
func Read(builder dependency.Builder) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		builder.
			WithModule(reader.Module).
			WithInvoke(reader.Run).
			Build().
			Run()
	}
}
