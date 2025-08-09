package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var version = "v1.0.0" // Define the version of your CLI

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

// updateGoMod updates the module name in the go.mod file to match the project name
func updateGoMod(projectName string) error {
	// Path to the go.mod file in the newly cloned project
	goModPath := filepath.Join(projectName, "go.mod")

	// Read the go.mod file
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod file: %w", err)
	}

	// Update the module name to match the project name
	updatedData := strings.Replace(string(data), "module go-project", fmt.Sprintf("module %s", projectName), 1)

	// Write the updated content back to the go.mod file
	err = os.WriteFile(goModPath, []byte(updatedData), 0644)
	if err != nil {
		return fmt.Errorf("failed to write go.mod file: %w", err)
	}

	return nil
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

		// Update the go.mod file with the correct module name
		if err := updateGoMod(projectName); err != nil {
			return fmt.Errorf("failed to update go.mod: %w", err)
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

// versionCmd is the 'version' command to show the current version of the CLI
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoBase CLI version %s\n", version)
	},
}

func main() {
	var rootCmd = &cobra.Command{Use: "gobase"}

	// Add the 'new' command to create a project
	rootCmd.AddCommand(createProjectCmd)
	// Add the 'version' command to show the current version
	rootCmd.AddCommand(versionCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
