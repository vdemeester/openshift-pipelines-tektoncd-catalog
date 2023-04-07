package main

import (
	"log"

	"github.com/spf13/cobra"
	// "golang.org/x/sync/errgroup"
)

func main() {
	if err := New().Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "paw",
		DisableAutoGenTag: true,
		SilenceUsage:      true,
	}

	cmd.AddCommand(Generate())
	return cmd
}
