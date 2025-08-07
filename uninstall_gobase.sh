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
