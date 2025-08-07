package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// cloneRepo clones the Go project template repository
func cloneRepo(projectName string) error {
	// Specify the repository URL
	repoURL := "https://github.com/toewailin/go-project" // Your repo URL
	cmd := exec.Command("git", "clone", repoURL, projectName)
	cmd.Dir = "./" // Set the directory to clone into
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// createProjectCmd is the 'new' command to create a new Go project
var createProjectCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Go project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please specify a project name")
		}

		projectName := args[0]
		// Clone the project template into the specified directory
		if err := cloneRepo(projectName); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}

		// Run any initialization commands, such as `go mod tidy` (optional)
		projectDir := filepath.Join("./", projectName)
		cmdInit := exec.Command("go", "mod", "tidy")
		cmdInit.Dir = projectDir
		if err := cmdInit.Run(); err != nil {
			return fmt.Errorf("failed to run 'go mod tidy': %w", err)
		}

		fmt.Println("Project created successfully!")
		fmt.Printf("Go to the project folder and run: cd %s\n", projectName)
		return nil
	},
}

func main() {
	var rootCmd = &cobra.Command{Use: "gobase"}

	// Add the 'new' command to create a project
	rootCmd.AddCommand(createProjectCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
