## To Fix
goi build should be for project, and goi serve (start)

### 2. **`config`** – Configure project settings or environment variables

* **Description**: A command to allow users to configure environment variables or project settings like database configurations, API keys, etc.
* **Usage**: `goi config set <key> <value>`
* **Example**: `goi config set DB_HOST localhost`
* **Alternative**: `goi config get <key>` to retrieve settings.

### 3. **`test`** – Run tests for the Go project

* **Description**: A simple command to run tests (`go test`) for the current Go project. It could have options for running tests in specific packages or files.
* **Usage**: `goi test [package]`
* **Example**: `goi test ./...` – Runs tests in the entire project.

2. **Generate Middleware**

   * **Description**: Generate middleware files for your Go project with predefined code for things like authentication, logging, CORS, or rate limiting.
   * **Usage**: `goi make middleware <name>`
   * **Example**: `goi make middleware auth` – Generates a new authentication middleware file.
   * **Benefit**: Reduces boilerplate code when adding common middleware to the project.

3. **Create Database Migrations**

   * **Description**: Generate a new database migration file for GORM or other ORM libraries.
   * **Usage**: `goi make migration <name>`
   * **Example**: `goi make migration add_users_table` – Generates a migration file to create a `users` table.
   * **Benefit**: Helps developers quickly generate migration files in a consistent format.

4. **Generate Custom Templates**

   * **Description**: Allow users to generate customized templates for various parts of their Go project. This could be based on parameters (e.g., CRUD operations, database structures, API definitions).
   * **Usage**: `goi make template <template_name>`
   * **Example**: `goi make template crud-api` – Generates a basic CRUD API template with routes, handlers, and models.
   * **Benefit**: Allows users to create their own project structures or templates based on common patterns in their team or organization.

5. **Create Configuration Files**

   * **Description**: Generate common configuration files, such as `.env` files, `docker-compose.yml`, or project-specific configurations (e.g., Redis, JWT, or database).
   * **Usage**: `goi make config <config_name>`
   * **Example**: `goi make config .env` – Generates a default `.env` configuration file for the project.
   * **Benefit**: Helps create the initial configuration structure, so users don’t have to manually create and remember every key.

