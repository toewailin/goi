package main

import (
	"fmt"
	"os"

	"goi/commands"

	"github.com/spf13/cobra"
)

const CLI_VERSION = "1.0.2"

// Main entry point of the application
func main() {
	var rootCmd = &cobra.Command{
		Use:     "goi",
		Short:   "goi is a CLI tool to manage Go projects",
		Long:    `goi is a command-line interface tool designed to streamline the creation, management, and deletion of Go projects.`,
		Version: CLI_VERSION,
	}

	// Add commands to the root command
	rootCmd.AddCommand(commands.CreateProjectCmd)
	rootCmd.AddCommand(commands.ListProjectCmd)
	rootCmd.AddCommand(commands.RemoveProjectCmd)
	rootCmd.AddCommand(commands.InstallDepsCmd)
	rootCmd.AddCommand(commands.AddDepCmd)
	rootCmd.AddCommand(commands.UpdateDepsCmd)
	rootCmd.AddCommand(commands.VersionCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", fmt.Sprintf("Error executing command: %v", err))
		os.Exit(1)
	}
}
