### **Documentation for `gobase` Installation and Uninstallation**

---

### **Introduction**

`gobase` is a Go project template generator that allows you to quickly set up a new Go project based on a pre-configured template. This documentation covers the process of installing and uninstalling `gobase` using the provided scripts.

---

### **Installation Script**

To install the `gobase` tool, you can execute a simple script that will download and set up the binary file on your system.

#### **Steps to Install `gobase`:**

1. **Run the Installation Command**

   Use `curl` to download and run the installation script directly from GitHub. This command will automatically download the `gobase` binary, make it executable, and place it in the appropriate directory.

   ```bash
   curl -o- https://raw.githubusercontent.com/toewailin/gobase/main/install_gobase.sh | bash
   ```

2. **What Happens During Installation:**

   * **Download**: The script downloads the `gobase` binary (e.g., `gobase-osx`) from the GitHub releases.
   * **Make Executable**: It sets the downloaded binary as executable using `chmod +x`.
   * **Move Binary**: It moves the binary to `/usr/local/bin/gobase`, making it accessible globally via the command line.
   * **Final Output**: Once the installation is complete, you should be able to use `gobase` as a command from anywhere on your system.

3. **Verify Installation**

   After the installation is complete, you can verify that `gobase` was installed correctly by running:

   ```bash
   gobase --version
   ```

   This should display the version of `gobase`, confirming that it was installed successfully.

---

### **Uninstallation Script**

To remove `gobase` from your system, you can execute the uninstallation script. This script removes the binary and any associated files.

#### **Steps to Uninstall `gobase`:**

1. **Run the Uninstallation Command**

   Use the following command to uninstall `gobase`:

   ```bash
   curl -o- https://raw.githubusercontent.com/toewailin/gobase/main/uninstall_gobase.sh | bash
   ```

2. **What Happens During Uninstallation:**

   * **Remove Binary**: The script checks if `gobase` is installed in `/usr/local/bin` and removes it.
   * **Remove Temporary Files**: If `gobase-osx` is found in the `~/Downloads` folder, it will be removed.
   * **Clean Up Project Template**: Optionally, the script can remove any project templates you may have cloned (this step is not mandatory).

3. **Verify Uninstallation**

   To confirm that `gobase` has been successfully uninstalled, you can run:

   ```bash
   which gobase
   ```

   This should return an empty result, indicating that `gobase` is no longer installed.

---

### **Scripts Details**

#### **Installation Script: `install_gobase.sh`**

```bash
#!/bin/bash

# Define the URL to download the binary
GOBASE_URL="https://github.com/toewailin/gobase/releases/download/v1.0.0-alpha/gobase-osx"

# Download the binary
echo "Downloading gobase binary..."
curl -L "$GOBASE_URL" -o ~/Downloads/gobase-osx

# Make the binary executable
echo "Making gobase binary executable..."
chmod +x ~/Downloads/gobase-osx

# Move the binary to /usr/local/bin
echo "Moving gobase binary to /usr/local/bin..."
sudo mv ~/Downloads/gobase-osx /usr/local/bin/gobase

echo "gobase installed successfully!"
```

* **Explanation**:

  * Downloads the binary from the GitHub release.
  * Makes it executable.
  * Moves the binary to `/usr/local/bin/gobase` for global access.

#### **Uninstallation Script: `uninstall_gobase.sh`**

```bash
#!/bin/bash

# Check if the binary exists in /usr/local/bin/gobase and remove it
if [ -f /usr/local/bin/gobase ]; then
  echo "Removing gobase binary from /usr/local/bin..."
  sudo rm /usr/local/bin/gobase
  echo "Successfully removed gobase binary."
else
  echo "gobase binary not found in /usr/local/bin."
fi

# Remove the downloaded gobase-osx file (if applicable)
if [ -f ~/Downloads/gobase-osx ]; then
  echo "Removing gobase-osx from Downloads..."
  rm ~/Downloads/gobase-osx
  echo "Successfully removed gobase-osx from Downloads."
else
  echo "gobase-osx file not found in Downloads."
fi

# Optionally, remove the cloned project or template repository if you want to clean up
# Uncomment the following lines if you'd like to remove the folder (not necessary in every case)
# if [ -d ~/go-project ]; then
#   echo "Removing project template folder..."
#   rm -rf ~/go-project
#   echo "Successfully removed project template folder."
# else
#   echo "Project template folder not found."
# fi

echo "Uninstallation completed."
```

* **Explanation**:

  * Removes the binary from `/usr/local/bin/`.
  * Optionally removes the `gobase-osx` binary from the `Downloads` folder.
  * Optionally removes the project template folder (uncomment to enable).

---

### **Conclusion**

With these installation and uninstallation scripts, users can easily manage the `gobase` tool on their systems. The process is simple, quick, and can be done entirely via the terminal. These scripts ensure that the installation and uninstallation are seamless and require minimal interaction from the user.

* **To Install**: Simply run the `install_gobase.sh` script via `curl` and `bash`.
* **To Uninstall**: Run the `uninstall_gobase.sh` script via `curl` and `bash`.

If you encounter any issues, please refer to the documentation or reach out to the project maintainers for further assistance.
