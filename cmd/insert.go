package cmd

import (
	"fmt"
	"path"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	"github.com/torbenconto/gopwd/internal/termio"
	"github.com/torbenconto/gopwd/internal/util"
)

var insertCmd = &cobra.Command{
	Use:               "insert [service] [flags]",
	Short:             "Insert a password for a service",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: AutocompleteServices,

	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		servicePath := path.Join(VaultPath, service) + ".gpg"

		// Flags
		copyFlag, _ := cmd.Flags().GetBool("copy")
		multilineFlag, _ := cmd.Flags().GetBool("multiline")

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

		err = util.CreateStructureAndClean(service, VaultPath, servicePath, encryptedPassword)

		if copyFlag {
			err = clipboard.WriteAll(password)
			if err != nil {
				return fmt.Errorf("failed to copy password to clipboard, error: %v", err)
			}
			fmt.Printf("Copied password for %s to clipboard\n", service)
			return nil
		}

		fmt.Printf("Password for %s inserted successfully\n", service)

		return nil
	},
}

func init() {
	insertCmd.Flags().BoolP("copy", "c", false, "Copy the password to the clipboard")
	insertCmd.Flags().BoolP("multiline", "m", false, "Input a multiline password")
	rootCmd.AddCommand(insertCmd)
}
