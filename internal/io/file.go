package io

import "os"

func CreateFile(file string) (*os.File, error) {
	return os.Create(file)
}

func RemoveFile(file string) error {
	return os.Remove(file)
}

func WriteToFile(file string, data []byte) error {
	return os.WriteFile(file, data, 0644)
}

func ReadFromFile(path string) ([]byte, error) {
	return os.ReadFile(path)
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
