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
			if entry.Name() == ".gpg-id" {
				continue
			}

			isLastEntry := i == len(dirEntries)-1
			entryName := entry.Name()

			if strings.HasSuffix(entryName, ".gpg") {
				entryName = strings.TrimSuffix(entryName, ".gpg")
			}

			var linePrefix string
			if isLastEntry {
				linePrefix = prefix + "    "
			} else {
				linePrefix = prefix + "│   "
			}

			if entry.IsDir() {
				if isLastEntry {
					fmt.Println(prefix + "└── " + entryName)
				} else {
					fmt.Println(prefix + "├── " + entryName)
				}
				printStructure(filepath.Join(path, entryName), linePrefix, isLastEntry)
			} else {
				if isLastEntry {
					fmt.Println(prefix + "└── " + entryName)
				} else {
					fmt.Println(prefix + "├── " + entryName)
				}
			}
		}
	}

	printStructure(vaultPath, "", false)
	return nil
}
