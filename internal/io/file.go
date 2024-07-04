package io

import (
	"io"
	"os"
)

func CreateFile(file string) (*os.File, error) {
	return os.Create(file)
}

func RemoveFile(file string) error {
	return os.Remove(file)
}

func WriteFile(file string, data []byte) error {
	return os.WriteFile(file, data, 0644)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func CopyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file for writing
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the content from the source file to the destination file
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return nil
}

func CreateTempFile(content []byte) (*os.File, error) {
	tmpfile, err := os.CreateTemp("", "tempfile")
	if err != nil {
		return nil, err
	}
	_, err = tmpfile.Write(content)
	if err != nil {
		return nil, err
	}
	tmpfile.Close()

	return tmpfile, nil
}
