package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "orca-temporal",
	Short: "Bridge ORCA IDE with Temporal.io for durable multi-agent orchestration",
	Long: `ORCA-Temporal Bridge

Connects ORCA IDE (local agent orchestration) with Temporal.io (durable workflow execution).
Enables coordinated multi-agent systems with full audit trails, failure recovery, and replay.`,
	Version: "0.1.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(workflowCmd)
	rootCmd.AddCommand(replayCmd)
}
