package commands

import (
	"fmt"
	"goi/utils" // Assuming goi/utils provides PrintSuccess etc.
	"os"
	"os/exec"
	"path/filepath" // For joining paths safely

	"github.com/spf13/cobra"
)

// mainPath stores the path to the main Go file, configurable by a flag.
var mainPath string

// ServeProjectCmd is the 'serve' command to start the Go project.
var ServeProjectCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Go project",
	Long: `The 'serve' command runs 'go run' to start the Go project.

By default, it tries to run 'cmd/api/main.go'. You can specify a different
path to your main executable file using the --path or -p flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current directory where the command is executed
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Determine the actual path to run
		targetMainFile := mainPath // Use the value from the flag

		// If the flag was not provided, try to find a default
		if targetMainFile == "" {
			// Check for cmd/api/main.go
			if _, err := os.Stat(filepath.Join(projectDir, "cmd", "api", "main.go")); err == nil {
				targetMainFile = filepath.Join("cmd", "api", "main.go")
			} else if _, err := os.Stat(filepath.Join(projectDir, "internal", "server", "main.go")); err == nil {
				targetMainFile = filepath.Join("internal", "server", "main.go")
			} else if _, err := os.Stat(filepath.Join(projectDir, "main.go")); err == nil {
				// Fallback to main.go in the root if cmd/api/main.go not found
				targetMainFile = "main.go"
			} else {
				// If neither common path exists, instruct the user to specify
				return fmt.Errorf("could not find 'main.go' in default paths (cmd/api/main.go or ./main.go).\n" +
					"Please specify the path to your main file using 'goi serve --path <your-main-file-path>'")
			}
		} else {
			// If a path was provided, make sure it exists
			fullPath := filepath.Join(projectDir, targetMainFile)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return fmt.Errorf("specified main file does not exist: %s", fullPath)
			}
		}

		utils.PrintInfo(fmt.Sprintf("Starting Go project from: %s", targetMainFile))

		// Build the 'go run' command
		runCmd := exec.Command("go", "run", targetMainFile)
		runCmd.Dir = projectDir // Ensure command runs from the project root

		// Pipe command output to current terminal
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Stdin = os.Stdin // Allow input (e.g., if the Go app needs it)

		// Run the command
		if err := runCmd.Run(); err != nil {
			return fmt.Errorf("failed to run the project: %w", err)
		}

		utils.PrintSuccess("Go project is running successfully!")
		return nil
	},
}

func init() {
	// Add the --path flag to the serve command
	ServeProjectCmd.Flags().StringVarP(&mainPath, "path", "p", "", "Path to the main Go executable file (e.g., cmd/api/main.go or main.go)")
}