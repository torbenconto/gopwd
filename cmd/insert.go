package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/util"
	"os"
	"path"
	"strings"
)

var insertCmd = &cobra.Command{
	Use:   "insert [service] [flags]",
	Short: "Insert a password for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		servicePath := path.Join(VaultPath, service) + ".gpg"
		var createdDirs []string
		var createdFiles []string
		var success = false

		defer func() {
			if !success {
				// Cleanup created directories and files if command fails
				for _, file := range createdFiles {
					if io.Exists(file) {
						err := os.Remove(file)
						if err != nil {
							fmt.Printf("Error cleaning up file %s: %v\n", file, err)
						}
					}
				}
				for _, dir := range createdDirs {
					if io.Exists(dir) {
						err := os.RemoveAll(dir)
						if err != nil {
							fmt.Printf("Error cleaning up directory %s: %v\n", dir, err)
						}
					}
				}
			}
		}()

		dirs := strings.Split(service, "/")
		var dirPath string
		if len(dirs) > 1 {
			dirPath = path.Join(VaultPath, strings.Join(dirs[:len(dirs)-1], "/"))
			if !io.Exists(dirPath) {
				err := os.MkdirAll(dirPath, 0755)
				if err != nil {
					return fmt.Errorf("failed to create directory structure for service: %s, error: %v", service, err)
				}
				createdDirs = append(createdDirs, dirPath)
			}
		}

		if !io.Exists(servicePath) {
			_, err := io.CreateFile(servicePath)
			if err != nil {
				return fmt.Errorf("failed to create .gpg file for service: %s, error: %v", service, err)
			}
			createdFiles = append(createdFiles, servicePath)
		} else {
			return fmt.Errorf("service already exists in vault: %s", service)
		}

		// Prompt the user for the password
		password, err := io.PromptPassword()
		if err != nil {
			return fmt.Errorf("failed to prompt for password, error: %v", err)
		}

		// Get the GPG ID from the .gpg-id file
		gpgID, err := util.ReadGPGID(path.Join(VaultPath, ".gpg-id"))
		if err != nil {
			return fmt.Errorf("failed to read GPG ID from .gpg-id file, error: %v", err)
		}

		// Encrypt the password and write it to the .gpg file

		GPG := gpg.NewGPG(gpgID, gpg.Config{})
		encryptedPassword, err := GPG.Encrypt([]byte(password))
		if err != nil {
			return fmt.Errorf("failed to encrypt password for service: %s, error: %v", service, err)
		}

		err = io.WriteToFile(servicePath, encryptedPassword)
		if err != nil {
			return fmt.Errorf("failed to write encrypted password to .gpg file for service: %s, error: %v", service, err)
		}

		if copyFlag, _ := cmd.Flags().GetBool("copy"); copyFlag {
			err = clipboard.WriteAll(password)
			if err != nil {
				return fmt.Errorf("failed to copy password to clipboard, error: %v", err)
			}
			fmt.Printf("Copied password for %s to clipboard\n", service)
			success = true
			return nil
		}

		success = true
		fmt.Printf("Password for %s inserted successfully\n", service)

		return nil
	},
}

func init() {
	insertCmd.Flags().BoolP("copy", "c", false, "Copy the password to the clipboard")
	rootCmd.AddCommand(insertCmd)
}
