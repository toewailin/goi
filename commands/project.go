package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

		// Update the import paths in all Go files
		if err := updateImportPaths(projectName); err != nil {
			return err
		}

		// Final Success Message
		utils.PrintSuccess("🎉 Your Go project has been created successfully!")
		utils.PrintInfo("🚀 To get started, follow these simple steps:")

		utils.PrintInfo("  1. Navigate to your project folder:")
		utils.PrintInfo("     cd " + projectName)

		utils.PrintInfo("  2. Install the dependencies and clean up the Go module:")
		utils.PrintInfo("     goi install  # Runs 'go mod tidy'")

		utils.PrintInfo("  3. Start your Go project:")
		utils.PrintInfo("     goi start    # This will run 'go run cmd/api/main.go'")

		utils.PrintInfo("\n🚀 Your API server will now be live at http://localhost:9090")

		return nil
	},
}

// cloneRepo clones the Go project template repository and shows output in the terminal.
func cloneRepo(projectName string) error {
	// Create the git clone command
	cmd := exec.Command("git", "clone", "--depth", "1", GO_PROJECT_TEMPLATE_URL, projectName)

	// Show git clone output and errors in the terminal
	cmd.Stdout = os.Stdout  // Display standard output (git clone details)
	cmd.Stderr = os.Stderr  // Display standard error output (git clone errors)

	utils.PrintInfo("Cloning project template...") // Print info message

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		utils.PrintError(fmt.Sprintf("Failed to clone repository: %v", err))
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	utils.PrintSuccess("Project cloned successfully!")
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
	utils.PrintInfo("Updating import paths in Go files...")

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

// Command to start the Go project
var StartProjectCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Go project",
	Long:  `The 'start' command runs 'go run cmd/api/main.go' to start the Go project.`,
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

// Command to list all Go projects
var ListProjectCmd = &cobra.Command{
	Use:   "list",
	Short: "List all created Go projects",
	Long:  `The 'list' command lists all created Go projects with basic information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current directory
		currentDir, err := os.Getwd()
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to get current directory: %v", err))
		}

		// Read the directory contents
		files, err := os.ReadDir(currentDir)
		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed to read the current directory: %v", err))
		}

		var projectCount int
		utils.PrintInfo("Listing Go Projects:")
		for _, file := range files {
			if file.IsDir() {
				goModPath := filepath.Join(currentDir, file.Name(), "go.mod")
				if _, err := os.Stat(goModPath); err == nil {
					// Found a Go project directory
					projectCount++
					utils.PrintSuccess(fmt.Sprintf(" - %s", file.Name())) // Success message for found project
				}
			}
		}

		// If no Go projects were found, print an info message
		if projectCount == 0 {
			utils.PrintInfo("No Go projects found.")
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
			utils.PrintError(fmt.Sprintf("Failed to get absolute path for directory '%s': %v", projectName, err))
			return fmt.Errorf("failed to get absolute path for directory '%s': %w", projectName, err)
		}

		// Check if the directory exists
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			utils.PrintError(fmt.Sprintf("Project directory '%s' does not exist", targetDir))
			return fmt.Errorf("project directory '%s' does not exist", targetDir)
		}

		// Remove the directory and its contents
		if err := os.RemoveAll(targetDir); err != nil {
			utils.PrintError(fmt.Sprintf("Failed to remove project directory '%s': %v", targetDir, err))
			return fmt.Errorf("failed to remove project directory '%s': %w", targetDir, err)
		}

		utils.PrintSuccess(fmt.Sprintf("Project '%s' removed successfully!", projectName))
		return nil
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
			utils.PrintError(fmt.Sprintf("Failed to get current directory: %v", err))
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Run `go mod tidy` to clean up dependencies
		goModTidyCmd := exec.Command("go", "mod", "tidy")
		goModTidyCmd.Dir = projectDir
		if err := goModTidyCmd.Run(); err != nil {
			utils.PrintError(fmt.Sprintf("Failed to run 'go mod tidy': %v", err))
			return fmt.Errorf("failed to run 'go mod tidy': %w", err)
		}

		utils.PrintSuccess("Dependencies installed successfully!")
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
			utils.PrintError(fmt.Sprintf("Failed to add package '%s': %v", packageName, err))
			return fmt.Errorf("failed to add package '%s': %w", packageName, err)
		}

		utils.PrintSuccess(fmt.Sprintf("Package '%s' installed successfully!", packageName))
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
			utils.PrintError(fmt.Sprintf("Failed to get current directory: %v", err))
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Run `go get -u` to update the dependencies
		updateCmd := exec.Command("go", "get", "-u")
		updateCmd.Dir = projectDir
		if err := updateCmd.Run(); err != nil {
			utils.PrintError(fmt.Sprintf("Failed to update dependencies: %v", err))
			return fmt.Errorf("failed to update dependencies: %w", err)
		}

		utils.PrintSuccess("Dependencies updated successfully!")
		return nil
	},
}

// Version command to show the current version of GoI
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of the goi CLI",
	Long:  `The 'version' command displays the current version of the goi command-line interface tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintSuccess(fmt.Sprintf("goi CLI version %s", CLI_VERSION)) // Success message in green
	},
}

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project for the specified platform(s)",
	RunE:  runBuildCommand,
}

func runBuildCommand(cmd *cobra.Command, args []string) error {
	// Set default build flags
	buildFlags := []string{
		"-ldflags", "-s -w", // Reduce binary size
		"-trimpath",         // Remove file paths from binary
	}

	// Default output name for the binary (for the current machine)
	// outputName := "goi" // Default name for the output binary
	platforms := []string{} // Initialize an empty slice

	// Parse flags for specific platforms
	if allFlag, _ := cmd.Flags().GetBool("all"); allFlag {
		platforms = append(platforms, "linux", "darwin", "windows")
	}

	if linuxFlag, _ := cmd.Flags().GetBool("linux"); linuxFlag {
		platforms = append(platforms, "linux")
	}

	if macFlag, _ := cmd.Flags().GetBool("mac"); macFlag {
		platforms = append(platforms, "darwin")
	}

	if windowsFlag, _ := cmd.Flags().GetBool("windows"); windowsFlag {
		platforms = append(platforms, "windows")
	}

	// If no flags are provided, default to the current platform
	if len(platforms) == 0 {
		platforms = append(platforms, runtime.GOOS)
	}

	// Iterate over the platforms and build for each one
	for _, platform := range platforms {
		var platformOutputName string

		// Detect the platform and architecture
		switch platform {
		case "linux":
			// Linux
			switch runtime.GOARCH {
			case "amd64":
				platformOutputName = "goi-linux-amd64"
			case "arm64":
				platformOutputName = "goi-linux-arm64"
			default:
				return fmt.Errorf("unsupported architecture %s for Linux", runtime.GOARCH)
			}
		case "darwin":
			// macOS (Apple Silicon or Intel)
			switch runtime.GOARCH {
			case "amd64":
				platformOutputName = "goi-macos-amd64"
			case "arm64":
				platformOutputName = "goi-macos-arm64"
			default:
				return fmt.Errorf("unsupported architecture %s for macOS", runtime.GOARCH)
			}
		case "windows":
			// Windows
			switch runtime.GOARCH {
			case "amd64":
				platformOutputName = "goi-win-amd64.exe"
			case "arm64":
				platformOutputName = "goi-win-arm64.exe"
			default:
				return fmt.Errorf("unsupported architecture %s for Windows", runtime.GOARCH)
			}
		default:
			return fmt.Errorf("unsupported operating system %s", platform)
		}

		// Construct the build command for the current platform
		cmdArgs := append([]string{"build"}, buildFlags...)
		cmdArgs = append(cmdArgs, "-o", "build/"+platformOutputName, ".")

		// Run the build command
		buildCommand := exec.Command("go", cmdArgs...)
		buildCommand.Stdout = os.Stdout
		buildCommand.Stderr = os.Stderr
		if err := buildCommand.Run(); err != nil {
			return fmt.Errorf("failed to run 'go build' for platform %s: %w", platform, err)
		}

		// Output success message using the utils
		utils.PrintSuccess(fmt.Sprintf("Successfully built Go project for %s and saved to build/%s", platform, platformOutputName))

		// If the --install flag is provided, move the binary to the correct directory
		installFlag, _ := cmd.Flags().GetBool("install")
		if installFlag {
			if err := moveBinaryToLocalBin(platformOutputName); err != nil {
				return err
			}
			utils.PrintSuccess("The binary has been moved to the appropriate directory")
		}
	}

	return nil
}

// Function to move the built binary to the correct directory based on the OS
func moveBinaryToLocalBin(binaryName string) error {
	// Check the current operating system
	switch runtime.GOOS {
	case "linux", "darwin":
		// For macOS and Linux, move the binary to /usr/local/bin
		destinationPath := "/usr/local/bin/goi"
		if err := os.Rename(filepath.Join("build", binaryName), destinationPath); err != nil {
			return fmt.Errorf("failed to move binary to /usr/local/bin: %w", err)
		}

		// Ensure the binary is executable
		if err := os.Chmod(destinationPath, 0755); err != nil {
			return fmt.Errorf("failed to make the binary executable: %w", err)
		}
	case "windows":
		// For Windows, move the binary to C:\Program Files\goi
		destinationPath := "C:\\Program Files\\goi\\goi.exe"
		if err := os.Rename(filepath.Join("build", binaryName), destinationPath); err != nil {
			return fmt.Errorf("failed to move binary to C:\\Program Files\\goi: %w", err)
		}
		// Optionally, update system PATH on Windows to include this directory (requires admin privileges)
		// Add your own logic here to update system PATH if necessary
	default:
		return fmt.Errorf("unsupported operating system %s", runtime.GOOS)
	}

	return nil
}

func init() {
	// Add the flags for specific platform builds
	BuildCmd.Flags().BoolP("all", "a", false, "Build for all platforms (linux, darwin, windows)")
	BuildCmd.Flags().BoolP("linux", "l", false, "Build for Linux")
	BuildCmd.Flags().BoolP("mac", "m", false, "Build for macOS")
	BuildCmd.Flags().BoolP("windows", "w", false, "Build for Windows")
	BuildCmd.Flags().BoolP("install", "i", false, "Move the binary to /usr/local/bin after build")
}

// Command to uninstall goi binary
var UninstallCmd = &cobra.Command{
    Use:   "uninstall",
    Short: "Uninstall goi CLI and remove the installed binary",
    Long:  `The 'uninstall' command removes the goi binary from your system and any associated files.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        var binaryPath string

        // Check the operating system and set the binary path accordingly
        switch runtime.GOOS {
        case "linux":
            binaryPath = "/usr/local/bin/goi" // Default install location on Linux
        case "darwin":
            binaryPath = "/usr/local/bin/goi" // Default install location on macOS
        case "windows":
            binaryPath = "C:\\Program Files\\goi\\goi.exe" // Default install location on Windows
        default:
            return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
        }

        // Check if the binary exists
        if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
            return fmt.Errorf("goi binary not found at %s", binaryPath)
        }

        // Remove the binary
        if err := os.Remove(binaryPath); err != nil {
            return fmt.Errorf("failed to remove goi binary: %w", err)
        }

        utils.PrintSuccess("goi binary removed successfully!")

        // Optional: Clean up other files (e.g., configuration files, logs, etc.)
        // Uncomment and update the following lines if needed for your project.
        // utils.CleanUpOtherFiles()

        return nil
    },
}


