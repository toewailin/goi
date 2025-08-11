package commands

import (
	"fmt"
	"goi/utils" // Assuming you have a utils package for printing success/error messages
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// DeployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the Go project to a cloud service or server",
	RunE:  runDeployCommand,
}

// runDeployCommand handles the deployment logic
func runDeployCommand(cmd *cobra.Command, args []string) error {
	// Get deployment target (Docker, Heroku, etc.)
	target, _ := cmd.Flags().GetString("target")

	switch target {
	case "docker":
		// Deploy using Docker
		if err := deployWithDocker(); err != nil {
			return fmt.Errorf("docker deployment failed: %w", err)
		}
	case "heroku":
		// Deploy to Heroku
		if err := deployWithHeroku(); err != nil {
			return fmt.Errorf("heroku deployment failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown target '%s'. supported targets are 'docker' or 'heroku'", target)
	}

	utils.PrintSuccess(fmt.Sprintf("deployment to %s successful!", target))
	return nil
}

// deployWithDocker builds a Docker image and pushes it to the Docker registry
func deployWithDocker() error {
	// Ensure there's a Dockerfile in the project
	if _, err := os.Stat("Dockerfile"); os.IsNotExist(err) {
		return fmt.Errorf("dockerfile not found in the current directory")
	}

	// Build Docker image
	cmd := exec.Command("docker", "build", "-t", "goi-project", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build docker image: %w", err)
	}

	// Optionally, push the image to a Docker registry (e.g., Docker Hub)
	// You can add this section if you want to push the image to a registry
	// cmdPush := exec.Command("docker", "push", "goi-project")
	// cmdPush.Stdout = os.Stdout
	// cmdPush.Stderr = os.Stderr
	// if err := cmdPush.Run(); err != nil {
	// 	return fmt.Errorf("failed to push docker image: %w", err)
	// }

	return nil
}

// deployWithHeroku deploys the project to Heroku
func deployWithHeroku() error {
	// Ensure the Heroku CLI is installed and the user is logged in
	if err := checkHerokuCLI(); err != nil {
		return err
	}

	// Create a Heroku app if it doesn't exist
	cmdCreateApp := exec.Command("heroku", "create")
	cmdCreateApp.Stdout = os.Stdout
	cmdCreateApp.Stderr = os.Stderr
	if err := cmdCreateApp.Run(); err != nil {
		return fmt.Errorf("failed to create heroku app: %w", err)
	}

	// Deploy to Heroku
	cmdDeploy := exec.Command("git", "push", "heroku", "master")
	cmdDeploy.Stdout = os.Stdout
	cmdDeploy.Stderr = os.Stderr
	if err := cmdDeploy.Run(); err != nil {
		return fmt.Errorf("failed to deploy to heroku: %w", err)
	}

	return nil
}

// checkHerokuCLI checks if the Heroku CLI is installed and the user is logged in
func checkHerokuCLI() error {
	cmd := exec.Command("heroku", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("heroku cli is not installed or the user is not logged in: %w", err)
	}
	return nil
}

// Initialize flags for the DeployCmd
func init() {
	// Add flags to specify the target for deployment (docker or heroku)
	DeployCmd.Flags().StringP("target", "t", "docker", "specify deployment target: 'docker' or 'heroku'")
}
