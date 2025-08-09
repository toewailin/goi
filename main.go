package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Constants for the CLI tool version and GitHub project template URL.
const CLI_VERSION = "1.0.2"
const GO_PROJECT_TEMPLATE_URL = "https://github.com/toewailin/go-project.git"

// Check if Git is installed
func checkGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// Clone a Go project template from GitHub (Suppressing output)
func cloneRepo(projectName string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", GO_PROJECT_TEMPLATE_URL, projectName)
	cmd.Stdout = nil  // Suppress the standard output (git clone details)
	cmd.Stderr = nil  // Suppress the standard error output (git clone errors)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	return nil
}

// Update the module path in the go.mod file (without extra output)
func updateGoMod(projectName string) error {
	goModPath := filepath.Join(projectName, "go.mod")

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod file: %w", err)
	}

	// Update the module name in go.mod
	lines := strings.Split(string(data), "\n")
	updatedLines := make([]string, len(lines))
	replaced := false
	for i, line := range lines {
		if strings.HasPrefix(line, "module ") && !replaced {
			updatedLines[i] = fmt.Sprintf("module %s", strings.ToLower(strings.ReplaceAll(projectName, "-", "/")))
			replaced = true
		} else {
			updatedLines[i] = line
		}
	}
	updatedData := strings.Join(updatedLines, "\n")

	err = os.WriteFile(goModPath, []byte(updatedData), 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated go.mod file: %w", err)
	}

	return nil
}

// Run `go mod tidy` to clean up dependencies (without output)
func runGoModTidy(projectDir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	cmd.Stdout = nil  // Suppress the output
	cmd.Stderr = nil  // Suppress the error output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %w", err)
	}
	return nil
}

// Command to create a new Go project
var createProjectCmd = &cobra.Command{
	Use:   "new <project_name>",
	Short: "Create a new Go project from a template",
	Long:  `The 'new' command initializes a fresh Go project by cloning a template from GitHub and setting up the Go module.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		if !checkGitInstalled() {
			return fmt.Errorf("git is not installed, please install git to clone the project")
		}

		// Check if target directory exists and is not empty
		targetDir, err := filepath.Abs(projectName)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for directory '%s': %w", projectName, err)
		}

		if _, err := os.Stat(targetDir); err == nil {
			dir, err := os.Open(targetDir)
			if err != nil {
				return fmt.Errorf("could not open directory '%s': %w", targetDir, err)
			}
			defer dir.Close()

			_, err = dir.Readdirnames(1)
			if err == nil {
				return fmt.Errorf("target directory '%s' already exists and is not empty, please specify an empty or non-existent directory", targetDir)
			}
			fmt.Println("Directory", targetDir, "exists but is empty. Proceeding with clone.")
		} else if os.IsNotExist(err) {
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", targetDir, err)
			}
		} else {
			return fmt.Errorf("error checking directory '%s': %w", targetDir, err)
		}

		// Clone the project template into the specified directory
		if err := cloneRepo(projectName); err != nil {
			return err
		}

		// Update the go.mod file with the correct module name
		if err := updateGoMod(projectName); err != nil {
			return err
		}

		// Run `go mod tidy` in the new project directory
		if err := runGoModTidy(projectName); err != nil {
			return err
		}

		// Final Success Message
		fmt.Println("╔════════════════════════════════════╗")
		fmt.Println("║   To Start Your Go Project, Run:   ║")
		fmt.Println("╚════════════════════════════════════╝")
		fmt.Println("  \033[1;33mcd", projectName, "&& go run cmd/api/main.go\033[0m")
		return nil
	},
}

// Command to display the current version of goi
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of the goi CLI",
	Long:  `The 'version' command displays the current version of the goi command-line interface tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("goi CLI version", CLI_VERSION)
	},
}

func main() {
	var rootCmd = &cobra.Command{
		Use:     "goi",
		Short:   "goi is a CLI tool to manage Go projects",
		Long:    `goi is a command-line interface tool designed to streamline the creation, management, and deletion of Go projects.`,
	}

	// Add custom information in the Long or Example field
	rootCmd.Long += "\n\nAuthor: Toe Wai Lin\nGitHub: https://github.com/toewailin/go-project"

	// Add commands to the root command
	rootCmd.AddCommand(createProjectCmd)
	rootCmd.AddCommand(versionCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", fmt.Sprintf("Error executing command: %v", err))
		os.Exit(1)
	}
}


