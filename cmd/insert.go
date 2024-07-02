package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/pkg/crypt"
	"github.com/torbenconto/gopwd/pkg/promts"
	"github.com/torbenconto/gopwd/pkg/util"
	"os"
	"path"
	"strings"
)

var insertCmd = &cobra.Command{
	Use:   "insert [service] [flags]",
	Short: "Insert a password for a service",
	Long:  "Insert a password for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		servicePath := path.Join(VaultPath, service) + ".gpg"

		// Split the service argument to handle directory structure
		dirs := strings.Split(service, "/")

		// If the service includes a directory structure, create it
		if len(dirs) > 1 {
			dirPath := path.Join(VaultPath, strings.Join(dirs[:len(dirs)-1], "/"))
			if !util.Exists(dirPath) {
				err := os.MkdirAll(dirPath, 0755) // Create the directory structure with appropriate permissions
				if err != nil {
					return fmt.Errorf("failed to create directory structure for service: %s, error: %v", service, err)
				}
			}
		}

		// Check if the .gpg file for the service already exists
		if !util.Exists(servicePath) {
			_, err := util.CreateFile(servicePath)
			if err != nil {
				return fmt.Errorf("failed to create .gpg file for service: %s, error: %v", service, err)
			}
		} else {
			return fmt.Errorf("service already exists in vault: %s", service)
		}

		// Prompt the user for the password
		password, err := promts.PromptPassword()
		if err != nil {
			return fmt.Errorf("failed to prompt for password, error: %v", err)
		}

		// Get the GPG ID from the .gpg-id file
		gpgID, err := util.ReadGPGID(path.Join(VaultPath, ".gpg-id"))
		if err != nil {
			return fmt.Errorf("failed to read GPG ID from .gpg-id file, error: %v", err)
		}

		// Encrypt the password and write it to the .gpg file

		GPG := crypt.NewGPG(gpgID, crypt.Config{})
		encryptedPassword, err := GPG.Encrypt([]byte(password))
		if err != nil {
			return fmt.Errorf("failed to encrypt password for service: %s, error: %v", service, err)
		}

		err = util.WriteToFile(servicePath, encryptedPassword)
		if err != nil {
			return fmt.Errorf("failed to write encrypted password to .gpg file for service: %s, error: %v", service, err)
		}

		if copytoclipboard {
			err = clipboard.WriteAll(password)
			if err != nil {
				return fmt.Errorf("failed to copy password to clipboard, error: %v", err)
			}
			fmt.Printf("Copied password for %s to clipboard", service)
			return nil
		}

		fmt.Printf("Password for %s inserted successfully\n", service)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)
}
