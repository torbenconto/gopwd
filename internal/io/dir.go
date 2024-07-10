package io

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func CreateDir(dir string) error {
	return os.Mkdir(dir, 0755)
}

func IsDirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Try to read at least one entry
	if errors.Is(err, io.EOF) {
		return true, nil // EOF means the directory is empty
	}
	return false, err // Return false if the directory is not empty or an error occurred
}

func Cleanup(files []string, dirs []string) error {
	for _, file := range files {
		if Exists(file) {
			err := os.Remove(file)
			if err != nil {
				return fmt.Errorf("error cleaning up file %s: %v\n", file, err)
			}
		}
	}
	for _, dir := range dirs {
		if Exists(dir) {
			err := os.RemoveAll(dir)
			if err != nil {
				return fmt.Errorf("error cleaning up directory %s: %v\n", dir, err)
			}
		}
	}

	return nil
}

// listServicesRecursive is a helper function that recursively searches for .gpg files.
func listServicesRecursive(vaultPath string, services *[]string, relativePath string) error {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(vaultPath, entry.Name())
		if entry.IsDir() {
			err := listServicesRecursive(fullPath, services, filepath.Join(relativePath, entry.Name()))
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(entry.Name(), ".gpg") {
			serviceRelativePath := filepath.Join(relativePath, strings.TrimSuffix(entry.Name(), ".gpg"))
			*services = append(*services, serviceRelativePath)
		}
	}

	return nil
}

// ListServices returns a slice of available service names by recursively reading .gpg files in the vault.
func ListServices(vaultPath string) ([]string, error) {
	var services []string
	err := listServicesRecursive(vaultPath, &services, "")
	if err != nil {
		return nil, err
	}
	return services, nil
}
