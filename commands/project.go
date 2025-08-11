package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"goi/config"
	"goi/utils"

	"github.com/spf13/cobra"
)

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

		// Update the import paths in all Go files
		if err := updateImportPaths(projectName); err != nil {
			return err
		}

		// Final Success Message
		utils.PrintSuccess("🎉 Your Go project has been created successfully!")

		utils.PrintInfo("🚀 To get started, follow these simple steps:")

		utils.PrintInfo("  1. Navigate to your project folder:")
		utils.PrintInfo("     cd " + projectName)

		utils.PrintInfo("  2. Install and update dependencies.")
		utils.PrintInfo("     goi sync  # Runs 'go mod tidy'")

		utils.PrintInfo("  3. Start your Go project:")
		utils.PrintInfo("     goi serve    # This will run 'go run cmd/api/main.go'")

		return nil
	},
}

// cloneRepo clones the Go project template repository and shows output in the terminal.
func cloneRepo(projectName string) error {
	// Create the git clone command
	cmd := exec.Command("git", "clone", "--depth", "1", config.GO_PROJECT_TEMPLATE_URL, projectName)

	// Show git clone output and errors in the terminal
	cmd.Stdout = os.Stdout  // Display standard output (git clone details)
	cmd.Stderr = os.Stderr  // Display standard error output (git clone errors)

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		utils.PrintError(fmt.Sprintf("Failed to clone repository: %v", err))
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Remove the .git directory to make it a fresh project
	gitDir := filepath.Join(projectName, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		if err := os.RemoveAll(gitDir); err != nil {
			// Print warning if unable to remove the .git directory
			utils.PrintWarning(fmt.Sprintf("WARNING: Failed to remove .git directory: %v", err))
		} else {
			utils.PrintSuccess("Removed .git directory.")
		}
	} else {
		utils.PrintInfo("No .git directory found.")
	}

	return nil
}

// Update the module path in the go.mod file (without extra output)
func updateGoMod(projectName string) error {
	goModPath := filepath.Join(projectName, "go.mod")

	// Read the go.mod file
	data, err := os.ReadFile(goModPath)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to read go.mod file: %v", err))
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

	// Write the updated go.mod file
	err = os.WriteFile(goModPath, []byte(updatedData), 0644)
	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to write updated go.mod file: %v", err))
		return fmt.Errorf("failed to write updated go.mod file: %w", err)
	}

	utils.PrintSuccess("go.mod updated successfully!")
	return nil
}

// Update import paths in Go files to reflect the project name
func updateImportPaths(projectName string) error {

	err := filepath.Walk(projectName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// Read the Go file
			data, err := os.ReadFile(path)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to read file %s: %v", path, err))
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			// Replace import paths
			updatedData := strings.ReplaceAll(string(data), "go-project", projectName)

			// Write the updated content back to the Go file
			err = os.WriteFile(path, []byte(updatedData), 0644)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to write file %s: %v", path, err))
				return fmt.Errorf("failed to write file %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		utils.PrintError(fmt.Sprintf("Failed to update import paths: %v", err))
	}

	utils.PrintSuccess("Import paths updated successfully!")
	return err
}


