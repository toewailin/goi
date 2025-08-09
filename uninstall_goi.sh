#!/bin/bash

# Check if the binary exists in /usr/local/bin/goi and remove it
if [ -f /usr/local/bin/goi ]; then
  echo "Removing goi binary from /usr/local/bin..."
  sudo rm /usr/local/bin/goi
  echo "Successfully removed goi binary."
else
  echo "goi binary not found in /usr/local/bin."
fi

# Remove the downloaded goi-osx file (if applicable)
if [ -f ~/Downloads/goi-osx ]; then
  echo "Removing goi-osx from Downloads..."
  rm ~/Downloads/goi-osx
  echo "Successfully removed goi-osx from Downloads."
else
  echo "goi-osx file not found in Downloads."
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
