package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/pwgen"
	"github.com/torbenconto/gopwd/internal/util"
	"os"
	"path"
	"strings"
)

var generateCmd = &cobra.Command{
	Use:               "generate [service] [flags]",
	Short:             "Generate a password for a service",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: AutocompleteServices,

	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		servicePath := path.Join(VaultPath, service) + ".gpg"
		var createdDirs []string
		var createdFiles []string
		var success = false

		length, _ := cmd.Flags().GetInt("length")
		symbols, _ := cmd.Flags().GetBool("symbols")
		numbers, _ := cmd.Flags().GetBool("numbers")
		uppercase, _ := cmd.Flags().GetBool("uppercase")
		lowercase, _ := cmd.Flags().GetBool("lowercase")
		human, _ := cmd.Flags().GetBool("memorable")
		copyFlag, _ := cmd.Flags().GetBool("copy")

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

		err = io.WriteFile(servicePath, encryptedPassword)
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
	generateCmd.Flags().IntP("length", "l", 16, "Length of the generated password")
	generateCmd.Flags().BoolP("symbols", "s", true, "Include symbols in the generated password")
	generateCmd.Flags().BoolP("numbers", "n", true, "Include numbers in the generated password")
	generateCmd.Flags().BoolP("uppercase", "u", true, "Include uppercase letters in the generated password")
	generateCmd.Flags().BoolP("lowercase", "L", true, "Include lowercase letters in the generated password")
	generateCmd.Flags().BoolP("memorable", "m", false, "Include words in the generated password")
	generateCmd.Flags().BoolP("copy", "c", false, "Copy the generated password to the clipboard")
	rootCmd.AddCommand(generateCmd)
}
