Sure! Here’s a documentation guide for creating a `gobase` project, uploading it as a release on GitHub, and setting it up for installation on macOS. After downloading the binary, the user will be able to make it executable, move it to `/usr/local/bin`, and use it to download the Go project template from the specified URL.

---

# **How to Create, Upload, and Install `gobase` Project**

This guide will walk you through the steps to:

1. **Create the `gobase` project.**
2. **Build it as a binary file.**
3. **Upload it to GitHub as a release.**
4. **Install the binary on macOS and make it executable.**
5. **Use the `gobase` command to download a Go project template.**

---

### **Step 1: Create the `gobase` Project**

The `gobase` CLI project allows you to easily set up a new Go project based on a predefined template from GitHub.

1. **Initialize your Go project**:

   * Create a new directory for the project and initialize it.

   ```bash
   mkdir gobase && cd gobase
   go mod init github.com/toewailin/gobase
   ```

2. **Write the `gobase` project code**:

   * Below is a sample code for `gobase` that clones the Go template repository from GitHub and sets it up.

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
       var rootCmd = &cobra.Command{Use: "gobase"}

       // Add the 'new' command to create a project
       rootCmd.AddCommand(createProjectCmd)

       // Execute the root command
       if err := rootCmd.Execute(); err != nil {
           fmt.Println(err)
           os.Exit(1)
       }
   }
   ```

3. **Build the `gobase` binary**:

   * Compile the Go project into a binary.

   ```bash
   go build -o gobase
    GOOS=linux GOARCH=amd64 go build -o gobase-linux
    GOOS=windows GOARCH=amd64 go build -o gobase.exe
   
   
   ```

   This will generate a binary file called `gobase` in the current directory.

---

### **Step 2: Upload the Binary to GitHub as a Release**

1. **Create a GitHub repository**:

   * Create a new GitHub repository (e.g., `gobase`).

2. **Upload the binary as a release**:

   * Go to the **Releases** section of your repository on GitHub.
   * Click **Draft a new release**.
   * Tag the release (e.g., `v1.0.0`).
   * Upload the `gobase` binary for your platform (e.g., `gobase-osx` for macOS).
   * Click **Publish release**.

---

### **Step 3: Install the Binary on macOS and Make it Executable**

1. **Download the binary**:

   * Go to the **Releases** section of your GitHub repository.
   * Download the appropriate binary for your operating system (e.g., `gobase-osx`).

2. **Make the binary executable**:

   * Open the terminal and navigate to the directory where the binary was downloaded.
   * Run the following command to make it executable:

   ```bash
   chmod +x gobase-osx
   ```

3. **Move the binary to `/usr/local/bin`**:

   * To make the `gobase` binary accessible globally, move it to `/usr/local/bin`.

   ```bash
   sudo mv gobase-osx /usr/local/bin/gobase
   ```

4. **Verify the installation**:

   * Run the following command to verify that `gobase` is installed correctly:

   ```bash
   gobase --version
   ```

   This should return the version of `gobase`.

---

### **Step 4: Use the `gobase` Command to Download a Go Project Template**

1. **Create a new Go project**:

   * Use the `gobase` command to create a new Go project from the GitHub template:

   ```bash
   gobase new myproject
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

You’ve now created the `gobase` CLI tool that allows users to quickly set up a Go project based on a predefined template from GitHub. By following this guide, you can efficiently distribute and install the `gobase` binary on macOS, as well as create new projects using a simple command.

This setup can be extended to support more features such as customizing the project template, adding different project configurations, and more.
