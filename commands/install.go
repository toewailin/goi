package commands

import (
	"fmt"
	"goi/utils" // Assuming you have a utils package for printing success/error messages
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or update the goi CLI tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Install or update the goi CLI tool
		if err := installGoBinary(); err != nil {
			utils.PrintError(fmt.Sprintf("Failed to install goi: %v", err))
			return err
		}

		utils.PrintSuccess("goi CLI installed successfully!")
		return nil
	},
}

// installGoBinary installs or updates the goi CLI tool
func installGoBinary() error {
	// Get the appropriate download URL based on the current OS and architecture
	downloadURL := getDownloadURL()

	// Determine the appropriate binary name for the current OS and architecture
	binaryName := getBinaryName()

	// Check if the binary exists in the build folder first
	if _, err := os.Stat(filepath.Join("build", binaryName)); err == nil {
		// If the binary exists in the build folder, move it to the correct directory
		utils.PrintSuccess("Using locally compiled binary...")
		if err := moveBinaryToLocalBin(filepath.Join("build", binaryName)); err != nil {
			return fmt.Errorf("failed to move locally compiled binary: %w", err)
		}
	} else {
		// Otherwise, download the binary
		utils.PrintSuccess("Downloading binary from GitHub...")
		if err := downloadAndInstallBinary(downloadURL, binaryName); err != nil {
			return err
		}
	}

	// Ensure the binary is executable
	if err := os.Chmod("/usr/local/bin/goi", 0755); err != nil {
		return fmt.Errorf("failed to make the binary executable: %w", err)
	}

	return nil
}

// getDownloadURL determines the download URL based on the current OS and architecture
func getDownloadURL() string {
	// Determine the download URL based on the current OS and architecture
	var url string
	switch runtime.GOOS {
	case "linux":
		url = "https://github.com/toewailin/go-project/releases/latest/download/goi-linux-amd64"
	case "darwin":
		url = "https://github.com/toewailin/go-project/releases/latest/download/goi-macos-amd64"
	case "windows":
		url = "https://github.com/toewailin/go-project/releases/latest/download/goi-win-amd64.exe"
	default:
		url = ""
	}
	return url
}

// getBinaryName determines the binary name based on the current OS and architecture
func getBinaryName() string {
	// Define the binary name based on OS and architecture
	var binaryName string
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			binaryName = "goi-linux-amd64"
		case "arm64":
			binaryName = "goi-linux-arm64"
		default:
			binaryName = "goi-linux-amd64" // Fallback to amd64 if architecture is not supported
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			binaryName = "goi-macos-amd64"
		case "arm64":
			binaryName = "goi-macos-arm64"
		default:
			binaryName = "goi-macos-amd64" // Fallback to amd64 if architecture is not supported
		}
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			binaryName = "goi-win-amd64.exe"
		case "arm64":
			binaryName = "goi-win-arm64.exe"
		default:
			binaryName = "goi-win-amd64.exe" // Fallback to amd64 if architecture is not supported
		}
	default:
		binaryName = "goi-linux-amd64" // Default for unsupported OS
	}
	return binaryName
}

// moveBinaryToLocalBin moves the locally compiled binary to the system's bin directory
func moveBinaryToLocalBin(binaryPath string) error {
	// Check the current operating system
	switch runtime.GOOS {
	case "linux", "darwin":
		// For macOS and Linux, move the binary to /usr/local/bin
		destinationPath := "/usr/local/bin/goi"
		if err := os.Rename(binaryPath, destinationPath); err != nil {
			return fmt.Errorf("failed to move binary to /usr/local/bin: %w", err)
		}

		// Ensure the binary is executable
		if err := os.Chmod(destinationPath, 0755); err != nil {
			return fmt.Errorf("failed to make the binary executable: %w", err)
		}
	case "windows":
		// For Windows, move the binary to C:\Program Files\goi
		destinationPath := "C:\\Program Files\\goi\\goi.exe"
		// Ensure that the folder exists, otherwise create it
		if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for binary: %w", err)
		}

		if err := os.Rename(binaryPath, destinationPath); err != nil {
			return fmt.Errorf("failed to move binary to C:\\Program Files\\goi: %w", err)
		}
		// Optionally, update system PATH on Windows to include this directory (requires admin privileges)
		// Add your own logic here to update system PATH if necessary
	default:
		return fmt.Errorf("unsupported operating system %s", runtime.GOOS)
	}

	return nil
}

// downloadAndInstallBinary downloads the binary from the specified URL and installs it
func downloadAndInstallBinary(url, binaryName string) error {
	// Use curl to download the binary based on the OS
	var destinationPath string
	switch runtime.GOOS {
	case "linux", "darwin":
		destinationPath = "/usr/local/bin/goi"
	case "windows":
		destinationPath = "C:\\Program Files\\goi\\goi.exe"
	default:
		return fmt.Errorf("unsupported operating system %s", runtime.GOOS)
	}

	// Download the binary from GitHub
	cmd := exec.Command("curl", "-L", url, "-o", destinationPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	return nil
}

