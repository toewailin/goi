package commands

import (
	"fmt"
	"goi/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

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

// Initialize flags for the BuildCmd
func init() {
	// Add flags for specific platform builds
	BuildCmd.Flags().BoolP("all", "a", false, "Build for all platforms (linux, darwin, windows)")
	BuildCmd.Flags().BoolP("linux", "l", false, "Build for Linux")
	BuildCmd.Flags().BoolP("mac", "m", false, "Build for macOS")
	BuildCmd.Flags().BoolP("windows", "w", false, "Build for Windows")
	BuildCmd.Flags().BoolP("install", "i", false, "Move the binary to /usr/local/bin after build")
}