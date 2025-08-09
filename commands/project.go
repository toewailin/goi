package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"goi/utils"

	"github.com/spf13/cobra"
)

var CLI_VERSION = "1.0.2"
const GO_PROJECT_TEMPLATE_URL = "https://github.com/toewailin/go-project.git"

// Command to create a new Go project
var CreateProjectCmd = &cobra.Command{
	Use:   "new <project_name>",
	Short: "Create a new Go project from a template",
	Long:  `The 'new' command initializes a fresh Go project by cloning a template from GitHub and setting up the Go module.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		// Check if Git and Go are installed
		if !utils.CheckGitInstalled() {
			return fmt.Errorf("git is not installed, please install git to clone the project")
		}

		if !utils.CheckGoInstalled() {
			return fmt.Errorf("go is not installed, please install Go to manage dependencies")
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
		fmt.Println("  cd", projectName, "&& go run cmd/api/main.go")
		return nil
	},
}

// Command to list all Go projects
var ListProjectCmd = &cobra.Command{
	Use:   "list",
	Short: "List all created Go projects",
	Long:  `The 'list' command lists all created Go projects with basic information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Assuming projects are in the current directory
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Look for directories that represent Go projects (directories with 'go.mod' inside)
		files, err := os.ReadDir(currentDir)
		if err != nil {
			return fmt.Errorf("failed to read the current directory: %w", err)
		}

		var projectCount int
		fmt.Println("Listing Go Projects:")
		for _, file := range files {
			if file.IsDir() {
				goModPath := filepath.Join(currentDir, file.Name(), "go.mod")
				if _, err := os.Stat(goModPath); err == nil {
					// Found a Go project directory
					projectCount++
					fmt.Println(" -", file.Name())
				}
			}
		}

		if projectCount == 0 {
			fmt.Println("No Go projects found.")
		}

		return nil
	},
}

// Command to remove a Go project
var RemoveProjectCmd = &cobra.Command{
	Use:   "remove <project_name>",
	Short: "Remove a Go project and its dependencies",
	Long:  `The 'remove' command removes a Go project from the local file system.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		targetDir, err := filepath.Abs(projectName)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for directory '%s': %w", projectName, err)
		}

		// Check if the directory exists
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			return fmt.Errorf("project directory '%s' does not exist", targetDir)
		}

		// Remove the directory and its contents
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to remove project directory '%s': %w", targetDir, err)
		}

		fmt.Println("Project", projectName, "removed successfully!")
		return nil
	},
}

// Version command to show the current version of GoI
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of the goi CLI",
	Long:  `The 'version' command displays the current version of the goi command-line interface tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("goi CLI version", CLI_VERSION)
	},
}

// Command to install dependencies (go mod tidy)
var InstallDepsCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dependencies using 'go mod tidy'",
	Long:  `The 'install' command runs 'go mod tidy' to install dependencies and clean up the Go module.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Run `go mod tidy` to clean up dependencies
		goModTidyCmd := exec.Command("go", "mod", "tidy")
		goModTidyCmd.Dir = projectDir
		if err := goModTidyCmd.Run(); err != nil {
			return fmt.Errorf("failed to run 'go mod tidy': %w", err)
		}

		fmt.Println("Dependencies installed successfully!")
		return nil
	},
}

// Command to add a dependency to the project
var AddDepCmd = &cobra.Command{
	Use:   "add <package_name>",
	Short: "Install a Go package",
	Long:  `The 'add' command adds a Go package using 'go get <package_name>'`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]

		// Run `go get` to install the package
		goGetCmd := exec.Command("go", "get", packageName)
		if err := goGetCmd.Run(); err != nil {
			return fmt.Errorf("failed to add package '%s': %w", packageName, err)
		}

		fmt.Printf("Package '%s' installed successfully!\n", packageName)
		return nil
	},
}

// Command to update Go dependencies
var UpdateDepsCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Go dependencies",
	Long:  `The 'update' command updates the dependencies in the go.mod file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Run `go get -u` to update the dependencies
		updateCmd := exec.Command("go", "get", "-u")
		updateCmd.Dir = projectDir
		if err := updateCmd.Run(); err != nil {
			return fmt.Errorf("failed to update dependencies: %w", err)
		}

		fmt.Println("Dependencies updated successfully!")
		return nil
	},
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
