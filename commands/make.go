package commands

import (
	"fmt"
	"goi/templates"
	"goi/utils"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// MakeCmd represents the 'make' command
var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate custom scaffolds for Go projects, goi make <type> <name>",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("subcommand is required. Example: goi make controller <name>")
	},
}

// Initialize the "make" subcommands
func init() {
	MakeCmd.AddCommand(MakeHandlerCmd)
	MakeCmd.AddCommand(MakeModelCmd)
	MakeCmd.AddCommand(MakeServiceCmd)
	MakeCmd.AddCommand(MakeRepositoryCmd)
	MakeCmd.AddCommand(MakeResponseCmd)
}

// MakeHandlerCmd generates a new handler file
var MakeHandlerCmd = &cobra.Command{
	Use:   "handler <name>",
	Short: "Generate a new handler",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(args[0], "handler")
	},
}

// MakeModelCmd generates a new model file
var MakeModelCmd = &cobra.Command{
	Use:   "model <name>",
	Short: "Generate a new model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(args[0], "model")
	},
}

// MakeServiceCmd generates a new service file
var MakeServiceCmd = &cobra.Command{
	Use:   "service <name>",
	Short: "Generate a new service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(args[0], "service")
	},
}

// MakeRepositoryCmd generates a new repository file
var MakeRepositoryCmd = &cobra.Command{
	Use:   "repo <name>",
	Short: "Generate a new repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(args[0], "repository")
	},
}

// MakeResponseCmd generates response files (success, error, pagination)
var MakeResponseCmd = &cobra.Command{
	Use:   "response",
	Short: "Generate response files (SuccessResponse, ErrorResponse, PaginationResponse)",
	RunE:  generateResponse,
}

// generateFile creates files based on the provided resource type and name
func generateFile(resourceName, resourceType string) error {
	// Ensure the directory exists for the resource type
	dir := getDirectoryForResource(resourceType)
	if err := ensureDirectoryExists(dir); err != nil {
		return err
	}

	// Get the module name dynamically from the go.mod file
	moduleName, err := getModuleNameFromGoMod(".") // Assuming the go.mod is in the current directory
	if err != nil {
		return fmt.Errorf("failed to get module name from go.mod: %w", err)
	}

	// Parse the appropriate template
	tmpl, err := parseTemplateForResource(resourceType)
	if err != nil {
		return err
	}

	// Create the file
	fileName := fmt.Sprintf("%s/%s_%s.go", dir, strings.ToLower(resourceName), strings.ToLower(resourceType))
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create %s file: %w", resourceType, err)
	}
	defer file.Close()

	// Capitalize the resourceType using cases.Title (proper Unicode handling)
	titleCase := cases.Title(language.Und, cases.Compact).String(resourceType)

	// Generate content from the template, passing the ModuleName along with the resource name
	err = tmpl.Execute(file, map[string]string{
		titleCase:      resourceName,
		"ModuleName":   moduleName,
	})
	if err != nil {
		return fmt.Errorf("failed to write to %s file: %w", resourceType, err)
	}

	utils.PrintSuccess(fmt.Sprintf("%s '%s' created successfully", titleCase, resourceName))
	return nil
}

// parseTemplateForResource selects the appropriate template for the resource type
func parseTemplateForResource(resourceType string) (*template.Template, error) {
	var tmplContent string

	// Choose the correct template based on resourceType
	switch resourceType {
	case "handler":
		tmplContent = templates.HandlerTemplate
	case "model":
		tmplContent = templates.ModelTemplate
	case "service":
		tmplContent = templates.ServiceTemplate
	case "repository":
		tmplContent = templates.RepositoryTemplate
	default:
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	// Parse and return the template
	return template.New(resourceType).Parse(tmplContent)
}

// getDirectoryForResource returns the correct directory based on the resource type
func getDirectoryForResource(resourceType string) string {
	switch resourceType {
	case "handler":
		return "handlers"
	case "model":
		return "models"
	case "service":
		return "services"
	case "repository":
		return "repository"
	default:
		return ""
	}
}

// generateResponse creates response files based on templates
func generateResponse(cmd *cobra.Command, args []string) error {
	responseFiles := map[string]string{
		"success_response.go":       templates.SuccessResponseTemplate,       // Already exists (SuccessResponse)
		"error_response.go":         templates.ErrorResponseTemplate,         // Already exists (ErrorResponse)
		"pagination_response.go":    templates.PaginationResponseTemplate,    // Already exists (Pagination)
	}

	if err := ensureDirectoryExists("response"); err != nil {
		return err
	}

	for fileName, tmplContent := range responseFiles {
		tmpl, err := template.New(fileName).Parse(tmplContent)
		if err != nil {
			return fmt.Errorf("failed to parse response template %s: %w", fileName, err)
		}

		filePath := fmt.Sprintf("response/%s", fileName)
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create response file %s: %w", filePath, err)
		}
		defer file.Close()

		if err := tmpl.Execute(file, nil); err != nil {
			return fmt.Errorf("failed to write to response file %s: %w", filePath, err)
		}

		fmt.Printf("Response file '%s' created successfully\n", fileName)
	}

	return nil
}

// ensureDirectoryExists checks if the directory exists, and creates it if it doesn't
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
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
