package main

import (
	"log"

	"github.com/BlackBX/service-framework/example/queue/cmd/read"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "queue",
		Short: "Example queue application",
		Long:  "An example application to show the running of a queue worker",
	}
	root.AddCommand(read.NewCommand())
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
