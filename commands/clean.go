package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

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

// Initialize flags for the CleanCmd
func init() {
	// --- Clean Command Flags ---
	CleanCmd.Flags().Bool("dry-run", false, "Show what would be cleaned without actually deleting anything")
	CleanCmd.Flags().Bool("verbose", false, "Show detailed logs of what is being cleaned")
	CleanCmd.Flags().Bool("clean-cache", false, "Clean the Go module cache (optional)")
}