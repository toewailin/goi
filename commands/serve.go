package commands

import (
	"fmt"
	"goi/utils"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// Command to Start the Go project
var ServeProjectCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Go project",
	Long:  `The 'serve' command runs 'go run cmd/api/main.go' to start the Go project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current directory
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Build the 'go run cmd/api/main.go' command
		runCmd := exec.Command("go", "run", "cmd/api/main.go")
		runCmd.Dir = projectDir

		// Run the command and show output in the terminal
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr

		// Run the command and check for errors
		if err := runCmd.Run(); err != nil {
			return fmt.Errorf("failed to run the project: %w", err)
		}

		// Print success message after running the command
		utils.PrintSuccess("Go project is running successfully!")

		return nil
	},
}