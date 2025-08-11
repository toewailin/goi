package commands

import (
	"fmt"
	"goi/utils"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

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