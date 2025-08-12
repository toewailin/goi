package commands

import (
	"fmt"
	"goi/utils" // Assuming goi/utils provides PrintSuccess, etc.
	"os"
	"os/exec"
	"path/filepath" // For joining paths safely
	"time"          // Add time package for timestamp generation

	"github.com/spf13/cobra"
)

// Flag variables for backup and restore
var backupPath string
var restorePath string
var mysqlUser string
var mysqlPassword string
var mysqlDatabase string

// MySQLBackupCmd is the 'backup' command to back up the MySQL database.
var MySQLBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Back up the MySQL database",
	Long: `The 'backup' command runs a MySQL dump command to create a backup of the database.

You can specify the path to save the backup file using the --path or -p flag.
You can also specify the MySQL credentials using the --user, --password, and --database flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current directory where the command is executed
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Determine the backup file path
		backupDir := backupPath
		if backupDir == "" {
			// Default to ./backups if no path is provided
			backupDir = "./backups"
		}

		// Create a timestamp for the backup file
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		backupFile := fmt.Sprintf("%s_backup_%s.sql", mysqlDatabase, timestamp)

		// Ensure the backup directory exists
		if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}

		// Full backup file path
		fullBackupPath := filepath.Join(backupDir, backupFile)

		// Validate MySQL credentials
		if mysqlUser == "" || mysqlPassword == "" || mysqlDatabase == "" {
			return fmt.Errorf("MySQL credentials are required: --user, --password, --database")
		}

		// Run the mysqldump command
		cmdArgs := []string{
			"-u", mysqlUser,
			"-p" + mysqlPassword, // Don't forget to use password directly after -p
			"-B", mysqlDatabase,   // Use the -B flag to specify the database
			"--result-file=" + fullBackupPath,
		}

		// Build the 'mysqldump' command
		backupCmd := exec.Command("mysqldump", cmdArgs...)
		backupCmd.Dir = projectDir // Ensure command runs from the project root

		// Pipe command output to current terminal
		backupCmd.Stdout = os.Stdout
		backupCmd.Stderr = os.Stderr
		backupCmd.Stdin = os.Stdin // Allow input (e.g., if the MySQL password is requested)

		// Run the command
		utils.PrintInfo(fmt.Sprintf("Starting MySQL backup to: %s", fullBackupPath))

		if err := backupCmd.Run(); err != nil {
			return fmt.Errorf("failed to back up MySQL database: %w", err)
		}

		utils.PrintSuccess("MySQL database backup was successful!")
		return nil
	},
}

// MySQLRestoreCmd is the 'restore' command to restore the MySQL database from a backup file.
var MySQLRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore the MySQL database from a backup file",
	Long: `The 'restore' command restores a MySQL database from a backup file.

You can specify the backup file to restore from using the --path or -p flag.
You can also specify the MySQL credentials using the --user, --password, and --database flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate MySQL credentials
		if mysqlUser == "" || mysqlPassword == "" || mysqlDatabase == "" {
			return fmt.Errorf("MySQL credentials are required: --user, --password, --database")
		}

		// Validate restore file path
		if restorePath == "" {
			return fmt.Errorf("backup file path is required: --path")
		}

		// Check if the backup file exists
		if _, err := os.Stat(restorePath); os.IsNotExist(err) {
			return fmt.Errorf("backup file does not exist: %s", restorePath)
		}

		// Run the mysql command to restore the database
		cmdArgs := []string{
			"-u", mysqlUser,
			"-p" + mysqlPassword, // Don't forget to use password directly after -p
			mysqlDatabase,         // Specify the database to restore to
			"<", restorePath,      // Use input redirection to load the backup
		}

		// Build the 'mysql' command
		restoreCmd := exec.Command("mysql", cmdArgs...)
		// Pipe command output to current terminal
		restoreCmd.Stdout = os.Stdout
		restoreCmd.Stderr = os.Stderr
		restoreCmd.Stdin = os.Stdin // Allow input if needed

		// Run the command
		utils.PrintInfo(fmt.Sprintf("Restoring MySQL database from: %s", restorePath))

		if err := restoreCmd.Run(); err != nil {
			return fmt.Errorf("failed to restore MySQL database: %w", err)
		}

		utils.PrintSuccess("MySQL database restore was successful!")
		return nil
	},
}

func init() {
	MySQLBackupCmd.Flags().StringVarP(&backupPath, "path", "b", "", "Path to save the backup file (e.g., /path/to/save/backup/)")
	MySQLBackupCmd.Flags().StringVarP(&mysqlUser, "user", "u", "", "MySQL username (e.g., root)")
	MySQLBackupCmd.Flags().StringVarP(&mysqlPassword, "password", "p", "", "MySQL password")
	MySQLBackupCmd.Flags().StringVarP(&mysqlDatabase, "database", "d", "", "MySQL database name to back up")

	MySQLRestoreCmd.Flags().StringVarP(&restorePath, "path", "r", "", "Path to the backup file to restore from (e.g., /path/to/backup.sql)")
	MySQLRestoreCmd.Flags().StringVarP(&mysqlUser, "user", "u", "", "MySQL username (e.g., root)")
	MySQLRestoreCmd.Flags().StringVarP(&mysqlPassword, "password", "p", "", "MySQL password")
	MySQLRestoreCmd.Flags().StringVarP(&mysqlDatabase, "database", "d", "", "MySQL database name to restore to")
}
