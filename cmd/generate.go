package cmd

import (
	"fmt"
	"path"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	"github.com/torbenconto/gopwd/internal/pwgen"
	"github.com/torbenconto/gopwd/internal/util"
)

var generateCmd = &cobra.Command{
	Use:               "generate [service] [flags]",
	Short:             "Generate a password for a service",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: AutocompleteServices,

	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		servicePath := path.Join(VaultPath, service) + ".gpg"

		length, _ := cmd.Flags().GetInt("length")
		symbols, _ := cmd.Flags().GetBool("symbols")
		numbers, _ := cmd.Flags().GetBool("numbers")
		uppercase, _ := cmd.Flags().GetBool("uppercase")
		lowercase, _ := cmd.Flags().GetBool("lowercase")
		human, _ := cmd.Flags().GetBool("memorable")
		copyFlag, _ := cmd.Flags().GetBool("copy")

		// Generate password
		generator := pwgen.NewPasswordGenerator(pwgen.PasswordGeneratorConfig{
			Length:    length,
			Humanized: human,
			Symbols:   symbols,
			Numbers:   numbers,
			Lowercase: lowercase,
			Uppercase: uppercase,
		})

		password, err := generator.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate password: %v", err)
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
		if err != nil {
			return fmt.Errorf("failed to create structure and clean up, error: %v", err)
		}

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
	generateCmd.Flags().IntP("length", "l", 16, "Length of the generated password")
	generateCmd.Flags().BoolP("symbols", "s", true, "Include symbols in the generated password")
	generateCmd.Flags().BoolP("numbers", "n", true, "Include numbers in the generated password")
	generateCmd.Flags().BoolP("uppercase", "u", true, "Include uppercase letters in the generated password")
	generateCmd.Flags().BoolP("lowercase", "L", true, "Include lowercase letters in the generated password")
	generateCmd.Flags().BoolP("memorable", "m", false, "Include words in the generated password")
	generateCmd.Flags().BoolP("copy", "c", false, "Copy the generated password to the clipboard")
	rootCmd.AddCommand(generateCmd)
}
