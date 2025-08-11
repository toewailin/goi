package commands

import (
	"fmt"           // Assuming you have a utils package for printing success/error messages
	"goi/templates" // Import the templates package
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// MakeCmd represents the 'make' command
var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate custom scaffolds for Go projects (e.g., models, controllers)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("subcommand is required. Example: goi make controller <name>")
	},
}

// MakeHandlerCmd generates a new handler file
var MakeHandlerCmd = &cobra.Command{
	Use:   "handler <name>",
	Short: "Generate a new handler",
	Args:  cobra.ExactArgs(1), // We expect exactly one argument: the name of the handler
	RunE:  generateHandler,
}

// MakeModelCmd generates a new model file
var MakeModelCmd = &cobra.Command{
	Use:   "model <name>",
	Short: "Generate a new model",
	Args:  cobra.ExactArgs(1), // We expect exactly one argument: the name of the model
	RunE:  generateModel,
}

// ServiceCmd represents the 'make service' command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Generate a new service for handling business logic",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("subcommand is required. Example: goi make service <name>")
	},
}

// MakeServiceCmd generates a new service file
var MakeServiceCmd = &cobra.Command{
	Use:   "service <name>",
	Short: "Generate a new service",
	Args:  cobra.ExactArgs(1), // We expect exactly one argument: the service name (e.g., User)
	RunE:  generateService,
}

// RepositoryCmd represents the 'make repository' command
var RepositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Generate a new repository for database interactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("subcommand is required. Example: goi make repository <name>")
	},
}

// MakeRepositoryCmd generates a new repository file
var MakeRepositoryCmd = &cobra.Command{
	Use:   "repo <name>",
	Short: "Generate a new repository",
	Args:  cobra.ExactArgs(1), // We expect exactly one argument: the repository name
	RunE:  generateRepository,
}

// Initialize flags for the MakeCmd
func init() {
	// Add the make command and its subcommands
	MakeCmd.AddCommand(MakeHandlerCmd)
	MakeCmd.AddCommand(MakeModelCmd)
	MakeCmd.AddCommand(MakeServiceCmd)
	MakeCmd.AddCommand(MakeRepositoryCmd)
}

// generateHandler creates a new handler file based on the name provided
func generateHandler(cmd *cobra.Command, args []string) error {
	handlerName := args[0]

	// Ensure the handlers directory exists
	if err := ensureDirectoryExists("handlers"); err != nil {
		return err
	}

	// Parse the handler template
	tmpl, err := template.New("handler").Parse(templates.HandlerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse handler template: %w", err)
	}

	// Create the handler file
	fileName := fmt.Sprintf("handlers/%s_handler.go", strings.ToLower(handlerName))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create handler file: %w", err)
	}
	defer file.Close()

	// Generate the handler content from the template
	err = tmpl.Execute(file, map[string]string{"HandlerName": handlerName})
	if err != nil {
		return fmt.Errorf("failed to write to handler file: %w", err)
	}

	fmt.Printf("Handler %s created successfully\n", handlerName)
	return nil
}

// generateModel creates a new model file based on the name provided
func generateModel(cmd *cobra.Command, args []string) error {
	modelName := args[0]

	// Ensure the models directory exists
	if err := ensureDirectoryExists("models"); err != nil {
		return err
	}

	// Parse the model template
	tmpl, err := template.New("model").Parse(templates.ModelTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse model template: %w", err)
	}

	// Create the model file
	fileName := fmt.Sprintf("models/%s.go", strings.ToLower(modelName))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create model file: %w", err)
	}
	defer file.Close()

	// Generate the model content from the template
	err = tmpl.Execute(file, map[string]string{"ModelName": modelName})
	if err != nil {
		return fmt.Errorf("failed to write to model file: %w", err)
	}

	fmt.Printf("Model %s created successfully\n", modelName)
	return nil
}

// generateService creates a new service file based on the name provided
func generateService(cmd *cobra.Command, args []string) error {
	serviceName := args[0]

	// Get the current working directory (the current project folder)
	projectPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get the module name from the current project's go.mod
	moduleName, err := getModuleNameFromGoMod(projectPath)
	if err != nil {
		return err
	}

	// Ensure the service directory exists
	if err := ensureDirectoryExists("service"); err != nil {
		return err
	}

	// Parse the service template
	tmpl, err := template.New("service").Parse(templates.ServiceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	// Create the service file
	fileName := fmt.Sprintf("service/%s_service.go", strings.ToLower(serviceName))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}
	defer file.Close()

	// Generate the service content from the template
	err = tmpl.Execute(file, map[string]string{
		"ServiceName": serviceName,
		"ModuleName":  moduleName,  // Pass the dynamically fetched module name
	})
	if err != nil {
		return fmt.Errorf("failed to write to service file: %w", err)
	}

	fmt.Printf("Service %s created successfully\n", serviceName)
	return nil
}

// getModuleNameFromGoMod reads the go.mod file in the specified project directory and extracts the module name
func getModuleNameFromGoMod(projectPath string) (string, error) {
	// Construct the path to the go.mod file in the current project directory
	goModPath := fmt.Sprintf("%s/go.mod", projectPath)

	// Read the go.mod file
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Extract the module name from the go.mod file
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			// Strip "module " and return the module name
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

// generateRepository creates a new repository file based on the name provided
func generateRepository(cmd *cobra.Command, args []string) error {
	repoName := args[0]

	// Ensure the repository directory exists
	if err := ensureDirectoryExists("repository"); err != nil {
		return err
	}

	// Parse the repository template
	tmpl, err := template.New("repository").Parse(templates.RepositoryTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse repository template: %w", err)
	}

	// Create the repository file
	fileName := fmt.Sprintf("repository/%s_repo.go", strings.ToLower(repoName))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create repository file: %w", err)
	}
	defer file.Close()

	// Generate the repository content from the template
	err = tmpl.Execute(file, map[string]string{"RepoName": repoName})
	if err != nil {
		return fmt.Errorf("failed to write to repository file: %w", err)
	}

	fmt.Printf("Repository %s created successfully\n", repoName)
	return nil
}

// ensureDirectoryExists checks if the directory exists, and creates it if it doesn't
func ensureDirectoryExists(dir string) error {
	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// If it doesn't exist, create the directory
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}