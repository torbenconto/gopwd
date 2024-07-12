package util

import (
	"fmt"
	"github.com/torbenconto/gopwd/internal/io"
	"os"
	"path"
	"strings"
)

// CreateStructureAndClean creates the directory structure for a service and cleans up if an error occurs, helper function for insert and generate commands
func CreateStructureAndClean(service, vaultPath, servicePath string, encrypted []byte) error {
	var createdDirs []string
	var createdFiles []string
	var success = false

	defer func() {
		if !success {
			// Cleanup created directories and files if command fails
			err := io.Cleanup(createdDirs, createdFiles)
			if err != nil {
				fmt.Printf("failed to cleanup created directories and files: %v\n", err)
			}
		}
	}()

	dirs := strings.Split(service, "/")
	var dirPath string
	if len(dirs) > 1 {
		dirPath = path.Join(vaultPath, strings.Join(dirs[:len(dirs)-1], "/"))
		if !io.Exists(dirPath) {
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory structure for service: %s, error: %v", service, err)
			}
			createdDirs = append(createdDirs, dirPath)
		}
	}

	_, err := io.CreateFile(servicePath)
	if err != nil {
		return fmt.Errorf("failed to create .gpg file for service: %s, error: %v", service, err)
	}
	createdFiles = append(createdFiles, servicePath)

	// Write the encrypted password to the file
	err = io.WriteFile(servicePath, encrypted)
	if err != nil {
		return fmt.Errorf("failed to write encrypted password to file: %v", err)
	}
	success = true

	return nil
}
