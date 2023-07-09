package util

import "os"

func GetHomeDir() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	return homeDir
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDirectory(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func RemoveDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(path string) error {
	_, err := os.Create(path)
	if err != nil {
		return err
	}
	return nil
}
