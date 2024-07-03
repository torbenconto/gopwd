package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
