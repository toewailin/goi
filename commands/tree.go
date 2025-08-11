package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

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

// Initialize flags for the TreeCmd
func init() {
	TreeCmd.Flags().BoolP("dirs", "d", false, "Show directories only, exclude files")
}