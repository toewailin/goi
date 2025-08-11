package main

import (
	"fmt"
	"os"
	"strings"

	"goi/commands"

	"github.com/spf13/cobra"
)

const (
	CLI_VERSION               = "1.0.2"
	GO_PROJECT_TEMPLATE_URL = "https://github.com/toewailin/go-project.git"
)

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
	rootCmd.AddCommand(commands.ServeProjectCmd)
	rootCmd.AddCommand(commands.SyncCmd)
	rootCmd.AddCommand(commands.BuildCmd)
	rootCmd.AddCommand(commands.DeployCmd)
	rootCmd.AddCommand(commands.CleanCmd)
	rootCmd.AddCommand(commands.MakeCmd)
	rootCmd.AddCommand(commands.InstallCmd)
	rootCmd.AddCommand(commands.UninstallCmd)
	rootCmd.AddCommand(commands.UpgradeCmd)
	rootCmd.AddCommand(commands.VersionCmd)
	rootCmd.AddCommand(commands.HistoryCmd)
	rootCmd.AddCommand(commands.TreeCmd)

// Hook into the 'Run' function of each command to save executed commands to history
	cobra.OnInitialize(func() {
		// Capture the current command name and arguments
		rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			// Combine the command name and its arguments
			executedCommand := fmt.Sprintf("%s %s", cmd.Use, strings.Join(args, " "))
			commands.SaveToHistory(executedCommand)
		}
	})

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", fmt.Sprintf("Error executing command: %v", err))
		os.Exit(1)
	}
}
