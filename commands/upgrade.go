package commands

import (
	"encoding/json"
	"fmt"
	"goi/config"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

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
	currentVersion := config.CLI_VERSION

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