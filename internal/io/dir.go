package io

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

// ListServices returns a slice of available service names by reading .gpg files in the vault.
func ListServices(vaultPath string) ([]string, error) {
	var services []string
	files, err := ioutil.ReadDir(vaultPath)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".gpg") {
			services = append(services, strings.TrimSuffix(f.Name(), ".gpg"))
		}
	}
	return services, nil
}
