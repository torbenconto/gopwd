package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	editor "github.com/torbenconto/gopwd/internal/editor_darwin"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/util"
	"os"
	"os/exec"
	"path"
)

var editCmd = &cobra.Command{
	Use:   "edit [service] [flags]",
	Short: "Edit a password for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceFile := path.Join(VaultPath, args[0]+".gpg")
		// Check if service exists
		if !io.Exists(serviceFile) {
			return fmt.Errorf("service %s not found", args[0])
		}

		// Get gpgid
		gpgid, err := util.ReadGPGID(path.Join(VaultPath, ".gpg-id"))
		if err != nil {
			return fmt.Errorf("failed to read gpg-id: %v", err)
		}

		// Read password from file
		password, err := io.ReadFile(serviceFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}

		// Decrypt password
		gpgmodule := gpg.NewGPG(gpgid, gpg.Config{})
		password, err = gpgmodule.Decrypt(password)
		if err != nil {
			return fmt.Errorf("failed to decrypt password: %v", err)
		}

		// Write password to temporary file
		tmpFile, err := io.CreateTempFile(password)
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %v", err)
		}

		editorExecuteable := editor.Editor()
		editorCmd := exec.Command(editorExecuteable, tmpFile.Name())
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		err = editorCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to open editor: %v", err)
		}

		// Check if password was changed
		newPassword, err := io.ReadFile(tmpFile.Name())
		if err != nil {
			return fmt.Errorf("failed to read temporary file: %v", err)
		}

		if string(password) == string(newPassword) {
			fmt.Println("Password was not changed")
			return nil
		}

		encryptedPassword, err := gpgmodule.Encrypt(newPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt password: %v", err)
		}

		err = io.WriteFile(serviceFile, encryptedPassword)
		if err != nil {
			return fmt.Errorf("failed to write encrypted password to file: %v", err)
		}

		// Remove temporary file
		err = os.Remove(tmpFile.Name())
		if err != nil {
			return fmt.Errorf("failed to remove temporary file: %v", err)
		}

		fmt.Println("Password updated successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
