package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/termio"
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

		// Flags
		copyFlag, _ := cmd.Flags().GetBool("copy")
		multilineFlag, _ := cmd.Flags().GetBool("multiline")

		defer func() {
			if !success {
				// Cleanup created directories and files if command fails
				err := io.Cleanup(createdDirs, createdFiles)
				if err != nil {
					fmt.Printf("failed to cleanup created directories and files: %v\n", err)
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

		var password string
		var err error

		if multilineFlag {
			fmt.Println("Enter the contents of the file and press Ctrl-D when finished.")
			password, err = termio.ReadMultiline()
		} else {
			password, err = termio.PromptPassword()
		}
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
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

		if copyFlag {
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
	insertCmd.Flags().BoolP("multiline", "m", false, "Input a multiline password")
	rootCmd.AddCommand(insertCmd)
}
