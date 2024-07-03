package io

import (
	"errors"
	"io"
	"os"
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
