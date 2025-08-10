package utils

import (
	"os/exec"
)

// Check if Git is installed
func CheckGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// Check if Go is installed
func CheckGoInstalled() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

