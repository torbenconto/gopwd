package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func GetHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func CreateDir(dir string) error {
	return os.Mkdir(dir, 0755)
}

func CreateFile(file string) (*os.File, error) {
	return os.Create(file)
}

func RemoveFile(file string) error {
	return os.Remove(file)
}

func WriteToFile(file string, data []byte) error {
	return os.WriteFile(file, data, 0644)
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

func PrintVaultStructure(vaultPath string) error {
	var printStructure func(path string, prefix string, isLast bool)
	printStructure = func(path string, prefix string, isLast bool) {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		for i, entry := range dirEntries {
			// Skip .gpg-id file
			if entry.Name() == ".gpg-id" {
				continue
			}

			isLastEntry := i == len(dirEntries)-1
			entryName := entry.Name()

			// Remove .gpg extension from file names
			if strings.HasSuffix(entryName, ".gpg") {
				entryName = strings.TrimSuffix(entryName, ".gpg")
			}

			// Determine the correct prefix
			var linePrefix string
			if isLast {
				linePrefix = prefix + "    " // No vertical line needed for the last item
			} else {
				linePrefix = prefix + "│   "
			}

			if entry.IsDir() {
				// Print directory with the appropriate tree branch
				if isLastEntry {
					fmt.Println(prefix + "└── " + entryName)
				} else {
					fmt.Println(prefix + "├── " + entryName)
				}
				// Recursively print the structure of the directory
				printStructure(filepath.Join(path, entryName), linePrefix, isLastEntry)
			} else {
				// Print file with the appropriate tree branch
				if isLastEntry {
					fmt.Println(prefix + "└── " + entryName)
				} else {
					fmt.Println(prefix + "├── " + entryName)
				}
			}
		}
	}

	// Start the recursive print from the root of the vault with an empty prefix
	printStructure(vaultPath, "", true)
	return nil
}
