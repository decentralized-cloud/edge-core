// Package cmd implements different commands that can be executed against EdgeCluster service
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/micro-business/go-core/pkg/util"
)

// NewRootCommand returns root CLI application command interface
func NewRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "edge-core",
		PreRun: func(cmd *cobra.Command, args []string) {
			printHeader()
		},
	}

	// Register all commands
	cmd.AddCommand(
		newStartCommand(),
		newVersionCommand(),
	)

	return cmd
}

func printHeader() {
	util.PrintInfo("Edge Cluster Serice")
}
