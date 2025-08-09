Here’s the updated guide for **`goi`** replacing **`gobase`** in the installation, usage, and project creation process. The updated guide will walk you through how to create, upload, and install `goi` on your system, using the CLI commands to generate new Go projects.

---

## **How to Create, Upload, and Install `goi` Project**

This guide will walk you through the steps to:

1. **Create the `goi` project.**
2. **Build it as a binary file.**
3. **Upload it to GitHub as a release.**
4. **Install the binary on macOS and make it executable.**
5. **Use the `goi` command to download a Go project template.**

---

### **Step 1: Create the `goi` Project**

The `goi` CLI project allows you to easily set up a new Go project based on a predefined template from GitHub.

1. **Initialize your Go project**:

   * Create a new directory for the project and initialize it.

   ```bash
   mkdir goi && cd goi
   go mod init github.com/toewailin/goi
   ```

2. **Write the `goi` project code**:

   Below is a sample code for `goi` that clones the Go template repository from GitHub and sets it up.

   ```go
   package main

   import (
       "fmt"
       "os"
       "os/exec"
       "path/filepath"

       "github.com/spf13/cobra"
   )

   // cloneRepo clones the Go project template repository
   func cloneRepo(projectName string) error {
       // Specify the repository URL
       repoURL := "https://github.com/toewailin/go-project" // Your repo URL
       cmd := exec.Command("git", "clone", repoURL, projectName)
       cmd.Dir = "./" // Set the directory to clone into
       cmd.Stdout = os.Stdout
       cmd.Stderr = os.Stderr
       return cmd.Run()
   }

   // createProjectCmd is the 'new' command to create a new Go project
   var createProjectCmd = &cobra.Command{
       Use:   "new",
       Short: "Create a new Go project",
       RunE: func(cmd *cobra.Command, args []string) error {
           if len(args) < 1 {
               return fmt.Errorf("please specify a project name")
           }

           projectName := args[0]
           // Clone the project template into the specified directory
           if err := cloneRepo(projectName); err != nil {
               return fmt.Errorf("failed to clone repository: %w", err)
           }

           // Run any initialization commands, such as `go mod tidy` (optional)
           projectDir := filepath.Join("./", projectName)
           cmdInit := exec.Command("go", "mod", "tidy")
           cmdInit.Dir = projectDir
           if err := cmdInit.Run(); err != nil {
               return fmt.Errorf("failed to run 'go mod tidy': %w", err)
           }

           fmt.Println("Project created successfully!")
           fmt.Printf("Go to the project folder and run: cd %s\n", projectName)
           return nil
       },
   }

   func main() {
       var rootCmd = &cobra.Command{Use: "goi"}

       // Add the 'new' command to create a project
       rootCmd.AddCommand(createProjectCmd)

       // Execute the root command
       if err := rootCmd.Execute(); err != nil {
           fmt.Println(err)
           os.Exit(1)
       }
   }
   ```

3. **Build the `goi` binary**:

   Compile the Go project into a binary.

   ```bash
   go build -o goi-osx
   GOOS=linux GOARCH=amd64 go build -o goi-linux
   GOOS=windows GOARCH=amd64 go build -o goi.exe
   ```

   This will generate a binary file called `goi` in the current directory.

---

### **Step 2: Upload the Binary to GitHub as a Release**

1. **Create a GitHub repository**:

   * Create a new GitHub repository (e.g., `goi`).

2. **Upload the binary as a release**:

   * Go to the **Releases** section of your repository on GitHub.
   * Click **Draft a new release**.
   * Tag the release (e.g., `v1.0.0`).
   * Upload the `goi` binary for your platform (e.g., `goi-osx` for macOS).
   * Click **Publish release**.

---

### **Step 3: Install the Binary on macOS and Make it Executable**

1. **Download the binary**:

   * Go to the **Releases** section of your GitHub repository.
   * Download the appropriate binary for your operating system (e.g., `goi-osx`).

2. **Make the binary executable**:

   * Open the terminal and navigate to the directory where the binary was downloaded.
   * Run the following command to make it executable:

   ```bash
   chmod +x goi-osx
   ```

3. **Move the binary to `/usr/local/bin`**:

   * To make the `goi` binary accessible globally, move it to `/usr/local/bin`.

   ```bash
   sudo mv goi-osx /usr/local/bin/goi
   ```

4. **Verify the installation**:

   * Run the following command to verify that `goi` is installed correctly:

   ```bash
   goi --version
   ```

   This should return the version of `goi`.

---

### **Step 4: Use the `goi` Command to Download a Go Project Template**

1. **Create a new Go project**:

   * Use the `goi` command to create a new Go project from the GitHub template:

   ```bash
   goi new myproject
   ```

   This will clone the Go project template from `https://github.com/toewailin/go-project` into the `myproject` directory.

2. **Go to your new project folder**:

   ```bash
   cd myproject
   ```

3. **Run the new Go project**:

   * Run the Go project setup commands (optional, like `go mod tidy`):

   ```bash
   go mod tidy
   ```

---

### **Conclusion**

You’ve now created the `goi` CLI tool that allows users to quickly set up a Go project based on a predefined template from GitHub. By following this guide, you can efficiently distribute and install the `goi` binary on macOS, as well as create new projects using a simple command.

This setup can be extended to support more features such as customizing the project template, adding different project configurations, and more.

Let me know if you need any further assistance or clarifications!
