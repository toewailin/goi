package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

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