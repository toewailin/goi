package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		utils.PrintInfo("     goi serve    # This will run 'go run cmd/api/main.go'")

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

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Go module (if not already initialized)",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the current directory
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}

		// Check if go.mod exists
		if !fileExists(filepath.Join(dir, "go.mod")) {
			// Run `go mod init`
			fmt.Println("Initializing Go module...")
			cmd := exec.Command("go", "mod", "init", filepath.Base(dir))
			cmd.Dir = dir
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error initializing Go module:", err)
				return
			}
			fmt.Println("Go module initialized successfully!")
		} else {
			fmt.Println("Go module already initialized!")
		}
	},
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

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

// syncCmd represents the sync command
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync project dependencies by pulling the latest versions",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the sync function to update dependencies
		err := syncDependencies()
		if err != nil {
			fmt.Println("Error syncing dependencies:", err)
			os.Exit(1)
		}
	},
}

// syncDependencies ensures that dependencies are up-to-date
func syncDependencies() error {
	// First, run `go mod tidy` to clean up dependencies
	fmt.Println("Running 'go mod tidy' to clean up modules...")
	cmd := exec.Command("go", "mod", "tidy")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %w", err)
	}

	// Then, run `go get` to ensure all modules are fetched and updated to their latest versions
	fmt.Println("Running 'go get' to fetch missing dependencies...")
	cmd = exec.Command("go", "get", "-u")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run 'go get': %w", err)
	}

	fmt.Println("Project dependencies synced successfully!")
	return nil
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

// cleanCmd represents the clean command
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up build artifacts and temporary files",
	Run: func(cmd *cobra.Command, args []string) {
		// Check for dry-run flag
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		verbose, _ := cmd.Flags().GetBool("verbose")
		cleanCache, _ := cmd.Flags().GetBool("clean-cache")

		// Clean all files under the "build" directory using pattern
		filesToClean := []string{
			"build/*",  // Use wildcard to clean all files under the build folder
		}

		// Handle dry-run mode: show what would be cleaned without deleting anything
		if dryRun {
			fmt.Println("Dry-run mode: The following files would be removed:")
			for _, file := range filesToClean {
				fmt.Printf("- %s\n", file)
			}
			return
		}

		// Perform the cleaning
		for _, filePattern := range filesToClean {
			if verbose {
				fmt.Printf("Cleaning files matching pattern: %s\n", filePattern)
			}

			// Use filepath.Glob to match all files under the build folder
			matchedFiles, err := filepath.Glob(filePattern)
			if err != nil {
				fmt.Printf("Error matching files: %v\n", err)
				return
			}

			// Iterate over the matched files and remove them
			for _, file := range matchedFiles {
				if info, err := os.Stat(file); err == nil {
					if info.IsDir() {
						if err := os.RemoveAll(file); err != nil {
							fmt.Printf("Error removing directory %s: %v\n", file, err)
						} else if verbose {
							fmt.Printf("Directory %s removed successfully\n", file)
						}
					} else {
						if err := os.Remove(file); err != nil {
							fmt.Printf("Error removing file %s: %v\n", file, err)
						} else if verbose {
							fmt.Printf("File %s removed successfully\n", file)
						}
					}
				} else {
					// If file doesn't exist, print a message
					if verbose {
						fmt.Printf("No such file or directory: %s\n", file)
					}
				}
			}
		}

		// Optionally clean Go module cache if the user requested it
		if cleanCache {
			removeGoCache(verbose)
		}

		// Success message
		fmt.Println("Cleanup completed!")
	},
}

// Helper function to remove cached Go modules
func removeGoCache(verbose bool) {
	cacheDirs := []string{
		filepath.Join(os.Getenv("GOPATH"), "pkg", "mod"), // Go module cache
	}

	for _, dir := range cacheDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			if verbose {
				fmt.Printf("Cleaning Go module cache at %s...\n", dir)
			}
			cmd := exec.Command("go", "clean", "-modcache")
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error cleaning Go module cache: %v\n", err)
			} else if verbose {
				fmt.Printf("Go module cache cleaned\n")
			}
		}
	}
}

func init() {
	// Add the flags for specific platform builds
	BuildCmd.Flags().BoolP("all", "a", false, "Build for all platforms (linux, darwin, windows)")
	BuildCmd.Flags().BoolP("linux", "l", false, "Build for Linux")
	BuildCmd.Flags().BoolP("mac", "m", false, "Build for macOS")
	BuildCmd.Flags().BoolP("windows", "w", false, "Build for Windows")
	BuildCmd.Flags().BoolP("install", "i", false, "Move the binary to /usr/local/bin after build")
	// --- Clean Command Flags ---
	// Flag to perform a dry run: show what would be cleaned without actually deleting anything
	CleanCmd.Flags().Bool("dry-run", false, "Show what would be cleaned without actually deleting anything")
	CleanCmd.Flags().Bool("verbose", false, "Show detailed logs of what is being cleaned")
	CleanCmd.Flags().Bool("clean-cache", false, "Clean the Go module cache (optional)")
	TreeCmd.Flags().BoolP("dirs", "d", false, "Show directories only, exclude files")
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

// upgradeCmd represents the upgrade command
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade goi to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the upgrade function
		err := upgrade()
		if err != nil {
			fmt.Println("Error upgrading goi:", err)
			os.Exit(1)
		}
	},
}

// upgrade checks for the latest version and updates goi if needed
func upgrade() error {
	// Use the CLI_VERSION constant to check the current version
	currentVersion := CLI_VERSION

	// Fetch the latest release version from GitHub
	latestVersion, err := getLatestVersionFromGitHub()
	if err != nil {
		return err
	}

	// Compare versions
	if currentVersion == latestVersion {
		fmt.Println("You are already on the latest version:", currentVersion)
		return nil
	}

	// Perform upgrade
	fmt.Printf("Upgrading goi to version %s...\n", latestVersion)
	err = performUpgrade(latestVersion)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully upgraded goi to version %s\n", latestVersion)
	return nil
}

// getLatestVersionFromGitHub fetches the latest release version from GitHub
func getLatestVersionFromGitHub() (string, error) {
	// Make an HTTP request to GitHub API to get the latest release
	resp, err := http.Get("https://api.github.com/repos/toewailin/goi/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response to get the latest version (tag_name)
	var releaseInfo struct {
		TagName string `json:"tag_name"`
	}

	// Parse the JSON data
	err = json.Unmarshal(body, &releaseInfo)
	if err != nil {
		return "", fmt.Errorf("failed to parse the JSON response: %w", err)
	}

	// Return the latest version found
	return releaseInfo.TagName, nil
}

// performUpgrade handles the actual upgrade process (downloading and installing the new version)
func performUpgrade(latestVersion string) error {
	// Example: Download the latest binary for the appropriate platform (Linux, macOS, Windows)
	// You can use a platform detection method to download the right binary for the user

	// Construct the download URL for the latest version
	downloadURL := fmt.Sprintf("https://github.com/toewailin/goi/releases/download/%s/goi-linux-amd64", latestVersion)

	// Download the binary
	fmt.Println("Downloading the latest version of goi...")
	cmd := exec.Command("curl", "-L", downloadURL, "-o", "/usr/local/bin/goi")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to download the latest version of goi: %w", err)
	}

	// Make the binary executable
	cmd = exec.Command("chmod", "+x", "/usr/local/bin/goi")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to make the binary executable: %w", err)
	}

	return nil
}

// historyCmd represents the history command
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Displays a list of recently executed commands",
	Run: func(cmd *cobra.Command, args []string) {
		err := showHistory()
		if err != nil {
			fmt.Println("Error displaying history:", err)
		}
	},
}

// showHistory reads the history file and displays the last 10 commands
func showHistory() error {
	// Get the path to the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %v", err)
	}

	// Construct the path to the history file (.goi_history in the user's home directory)
	historyFilePath := filepath.Join(homeDir, ".goi_history")

	// Open the history file
	file, err := os.Open(historyFilePath)
	if err != nil {
		return fmt.Errorf("could not open history file: %v", err)
	}
	defer file.Close()

	// Read the lines from the history file
	var commands []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}

	// Check for errors reading the file
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading history file: %v", err)
	}

	// Display the last 10 commands (or fewer if there aren't enough)
	numCommands := len(commands)
	if numCommands > 10 {
		commands = commands[numCommands-10:]
	}

	// Print the commands to the console
	fmt.Println("Recent commands:")
	for i, cmd := range commands {
		fmt.Printf("%d: %s\n", i+1, cmd)
	}

	return nil
}

// saveToHistory appends a command to the history file
func SaveToHistory(command string) {
	// Get the path to the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}

	// Construct the path to the history file (.goi_history in the user's home directory)
	historyFilePath := filepath.Join(homeDir, ".goi_history")

	// Open the history file in append mode
	file, err := os.OpenFile(historyFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening history file:", err)
		return
	}
	defer file.Close()

	// Write the command to the file
	_, err = file.WriteString(command + "\n")
	if err != nil {
		fmt.Println("Error saving command to history:", err)
	}
}

// treeCmd represents the tree command
var TreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Display the project folder structure (or current directory), excluding hidden files/folders",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if -d flag is set (show directories only)
		showDirsOnly, _ := cmd.Flags().GetBool("dirs")
		err := displayTree(".", "", 0, 0, showDirsOnly) // Start from the current directory with proper formatting
		if err != nil {
			fmt.Println("Error displaying folder structure:", err)
		}
	},
}

// displayTree recursively lists the files and directories in a tree format
// It also counts directories and files recursively, ignoring hidden files/folders.
func displayTree(dir string, indent string, dirsCount, filesCount int, showDirsOnly bool) error {
	// Open the specified directory
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not open directory %s: %v", dir, err)
	}

	// Skip hidden files and directories
	dirEntries = filterHiddenFiles(dirEntries)

	// Print the root directory (.)
	if dir == "." {
		fmt.Println(".")
	}

	// Iterate through the directory entries
	for i, entry := range dirEntries {
		// If we are only showing directories and the entry is not a directory, skip it
		if showDirsOnly && !entry.IsDir() {
			continue
		}

		// If the entry is a directory, print it
		if entry.IsDir() {
			// Print the directory with appropriate indentation
			if i == len(dirEntries)-1 {
				// Last directory in this level, use "└──"
				printIndented(entry.Name(), indent, true)
			} else {
				// Not the last, use "├──"
				printIndented(entry.Name(), indent, false)
			}

			// Increment the directory counter
			dirsCount++

			// Recursively display subdirectories
			err = displayTree(filepath.Join(dir, entry.Name()), indent+"│   ", dirsCount, filesCount, showDirsOnly)
			if err != nil {
				return err
			}
		} else {
			// For files, print them in the correct format (only if not in dirs-only mode)
			if !showDirsOnly {
				if i == len(dirEntries)-1 {
					// Last file in this level, use "└──"
					printIndented(entry.Name(), indent, true)
				} else {
					// Not the last, use "├──"
					printIndented(entry.Name(), indent, false)
				}
				// Increment the file counter
				filesCount++
			}
		}
	}

	// Print the final count of directories and files
	if dir == "." {
		fmt.Printf("\n%d directories, %d files\n", dirsCount, filesCount)
	}

	return nil
}

// printIndented prints the file/folder with indentation based on the depth level
func printIndented(name string, indent string, isLast bool) {
	// Only remove the extra '│' if indent is not empty
	if len(indent) > 0 {
		indent = indent[:len(indent)-1] // Remove the extra '│' at the end
	}

	// Add the appropriate branch character (├── or └──)
	if isLast {
		indent += "└── "
	} else {
		indent += "├── "
	}

	// Print the name with indentation
	fmt.Println(indent + name)
}

// filterHiddenFiles filters out hidden files and directories that start with "."
func filterHiddenFiles(entries []os.DirEntry) []os.DirEntry {
	var filtered []os.DirEntry
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
