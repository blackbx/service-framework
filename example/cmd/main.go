package main

import (
	"log"

	"github.com/BlackBX/service-framework/example/cmd/serve"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "example",
		Short: "Example application",
		Long:  "An example application to show how to use the service-framework",
	}
	root.AddCommand(serve.NewCommand())
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
