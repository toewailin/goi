package utils

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty" // Package to check if stdout is a terminal
)

// colorize applies ANSI color codes to a string for colored output
func colorize(colorCode, text string) string {
	if isTerminal() {
		return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, text)
	}
	return text
}

// isTerminal checks if stdout is a terminal that supports colors
func isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// PrintInfo prints an informational message in blue
func PrintInfo(message string) {
	fmt.Println(colorize("0;34", "INFO: "+message))
}

// PrintSuccess prints a success message in green
func PrintSuccess(message string) {
	fmt.Println(colorize("0;32", "SUCCESS: "+message))
}

// PrintError prints an error message in red and exits the program
func PrintError(message string) {
	fmt.Fprintln(os.Stderr, colorize("0;31", "ERROR: "+message))
	os.Exit(1)
}
