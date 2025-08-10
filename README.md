## **goi**

### **Introduction**

`goi` is a Go project template generator that allows you to quickly set up a new Go project based on a pre-configured template. This documentation covers the process of installing and uninstalling `goi` using the provided scripts.

---

### **Installation Script**

To install the `goi` tool, you can execute a simple script that will download and set up the binary file on your system.

#### **Steps to Install `goi`:**

1. **Run the Installation Command**

   Use `curl` to download and run the installation script directly from GitHub. This command will automatically download the `goi` binary, make it executable, and place it in the appropriate directory.

   ```bash
   curl -o- https://raw.githubusercontent.com/toewailin/goi/main/install_goi.sh | bash
   ```

2. **What Happens During Installation:**

   * **Download**: The script downloads the `goi` binary (e.g., `goi-osx`) from the GitHub releases.
   * **Make Executable**: It sets the downloaded binary as executable using `chmod +x`.
   * **Move Binary**: It moves the binary to `/usr/local/bin/goi`, making it accessible globally via the command line.
   * **Final Output**: Once the installation is complete, you should be able to use `goi` as a command from anywhere on your system.

3. **Verify Installation**

   After the installation is complete, you can verify that `goi` was installed correctly by running:

   ```bash
   goi --version
   ```

   This should display the version of `goi`, confirming that it was installed successfully.

---

### **Uninstallation Script**

To remove `goi` from your system, you can execute the uninstallation script. This script removes the binary and any associated files.

#### **Steps to Uninstall `goi`:**

1. **Run the Uninstallation Command**

   Use the following command to uninstall `goi`:

   ```bash
   curl -o- https://raw.githubusercontent.com/toewailin/goi/main/uninstall_goi.sh | bash
   ```

2. **What Happens During Uninstallation:**

   * **Remove Binary**: The script checks if `goi` is installed in `/usr/local/bin` and removes it.
   * **Remove Temporary Files**: If `goi-osx` is found in the `~/Downloads` folder, it will be removed.
   * **Clean Up Project Template**: Optionally, the script can remove any project templates you may have cloned (this step is not mandatory).

3. **Verify Uninstallation**

   To confirm that `goi` has been successfully uninstalled, you can run:

   ```bash
   which goi
   ```

   This should return an empty result, indicating that `goi` is no longer installed.

---

### **Go Project Scaffolded Application Structure**

This directory structure is designed to provide a solid foundation for a Golang application. It includes the main entry points for the application, configuration files, and templates for generating controllers, models, and migrations. The `config` directory contains configuration files for the database, Redis, and JWT, while the `scaffolding` directory contains templates for generating controllers, models, and migrations. The `cli` directory contains the CLI tools entry point, which can be used for generating scaffolding files.

---

### **Usage**

After installation, you can run the `goi` CLI with the following command format:

```bash
goi [command]
```

#### **Available Commands:**

* **`completion`**: Generate the autocompletion script for the specified shell.

  Example usage:

  ```bash
  goi completion
  ```

* **`help`**: Display help about any command.

  Example usage:

  ```bash
  goi help
  ```

* **`new`**: Create a new Go project using a project template (based on the `goi` templates).

  Example usage:

  ```bash
  goi new project-name
  ```

* **`version`**: Display the current version of `goi`.

  Example usage:

  ```bash
  goi version
  ```

---

### **Build Commands for Cross-Compilation**

To build the Go project for different platforms, you can use the following commands with cross-compilation options for **Linux**, **MacOS**, and **Windows**.

#### **Build the Project for the Current Platform:**

```bash
  chmod +x build.sh
  ./build.sh
```

### **Project Structure**

```plaintext
myapp/
├── cmd/                                  # Main entry points for the app
│   ├── api/                              # API server entry point
│   │   └── main.go                       # Main entry point for the API (Gin setup)
│   ├── worker/                           # Cron job or background worker entry point
│   │   └── main.go                       # Starts background processing tasks (e.g., cron jobs)
│   ├── cli/                              # CLI tools entry point (optional)
│   │   └── main.go                       # CLI commands for different tasks (including scaffolding)
├── config/                               # Centralized configuration files
│   ├── config.go                         # Global configuration (DB, Redis, JWT)
│   ├── db_config.go                      # Database configuration (PostgreSQL)
│   ├── redis_config.go                   # Redis connection and config (pub/sub)
│   ├── environment.go                    # Environment-specific settings (dev, prod)
│   └── .env                              # Environment variables
├── scaffolding/                          # Scaffolding files and templates
│   ├── controller_template.go            # Template for generating controllers
│   ├── model_template.go                 # Template for generating models
│   ├── migration_template.go             # Template for generating migrations
│   ├── factory_template.go               # Template for generating factories
│   └── route_template.go                 # Template for generating routes
├── handler/                              # HTTP handlers (controllers)
│   ├── adminHandler.go                   # Admin logic (system health, user management)
├── models/                               # GORM models (database entities)
│   ├── user.go                           # User model (all roles)
├── repository/                           # Database interaction layer
│   ├── userRepo.go                       # CRUD for users
├── services/                             # Core business logic
│   ├── authService.go                    # Authentication logic (JWT, sessions)
│   └── cronService.go                    # Cron job or background task logic
├── routes/                               # Route definitions
│   ├── adminRoutes.go                    # Admin-related routes
│   └── gameRoutes.go                     # Game-related routes
├── middlewares/                          # Middleware
│   ├── authMiddleware.go                 # JWT authentication
│   ├── loggingMiddleware.go              # Logging requests/responses
│   ├── rateLimitMiddleware.go            # Rate-limiting
│   ├── roleMiddleware.go                 # Role-based access control
│   ├── errorHandlingMiddleware.go        # Centralized error handler
│   ├── corsMiddleware.go                 # CORS handling
├── utils/                                # Utility functions
│   ├── jwt.go                            # JWT helpers
│   ├── logger.go                         # Logger
│   ├── validation.go                     # Input validation
│   ├── pagination.go                     # Pagination helper
│   ├── i18n.go                           # Translations (i18n helpers)
├── locales/                              # Translation files
│   ├── en.json                           # English translations
│   ├── my.json                           # Myanmar translations
├── tests/                                # Unit and integration tests
│   ├── controllers/
│   │   └── adminController_test.go
│   ├── services/
│   │   └── authService_test.go
│   └── repositories/
│       └── userRepository_test.go
├── scripts/                              # Database migrations and scripts
│   ├── genkeys.sh                        # Generate PEM keys
│   └── upload.sh                         # Upload script
├── docker/                               # Docker configuration
│   ├── Dockerfile                        # Dockerfile for production
│   ├── Dockerfile.dev                    # Dockerfile for development
│   └── docker-compose.yml                # Docker-compose setup
├── go.mod                                # Go module dependencies
└── go.sum                                # Go module checksum
```

---

### **Conclusion**

With these installation and uninstallation scripts, users can easily manage the `goi` tool on their systems. The process is simple, quick, and can be done entirely via the terminal. These scripts ensure that the installation and uninstallation are seamless and require minimal interaction from the user.
