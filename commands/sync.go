package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

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