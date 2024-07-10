package cmd

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/crypt/gpg"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/qr"
	"github.com/torbenconto/gopwd/internal/util"
	"os"
	"path"
	"strings"
)

var showCmd = &cobra.Command{
	Use:   "show [service] [flags]",
	Short: "Show a password for a service",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]

		// Flags
		qrFlag, _ := cmd.Flags().GetBool("qr")
		copyFlag, _ := cmd.Flags().GetBool("copy")
		lineNumber, _ := cmd.Flags().GetInt("line")

		// Check if service exists
		if !io.Exists(path.Join(VaultPath, service+".gpg")) {
			return fmt.Errorf("service %s not found", service)
		}

		// Get password and decrypt
		file, err := io.ReadFile(path.Join(VaultPath, service+".gpg"))
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}

		GPGID, err := util.ReadGPGID(path.Join(VaultPath, ".gpg-id"))
		if err != nil {
			return fmt.Errorf("failed to read gpg-id: %v", err)
		}

		GPGModule := gpg.NewGPG(GPGID, gpg.Config{})

		password, err := GPGModule.Decrypt(file)
		if err != nil {
			return fmt.Errorf("failed to decrypt password: %v", err)
		}

		lines := strings.Split(string(password), "\n")
		var nonEmptyLines []string
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				nonEmptyLines = append(nonEmptyLines, line)
			}
		}

		if lineNumber > 0 && lineNumber <= len(nonEmptyLines) {
			line := nonEmptyLines[lineNumber-1]
			if line == "" {
				return fmt.Errorf("line %s is empty", line)
			}
			if qrFlag {
				qr.Generate(string(password), qr.M, os.Stdout)
			}
			if copyFlag {
				err = clipboard.WriteAll(line)
				if err != nil {
					return fmt.Errorf("failed to copy line to clipboard: %w", err)
				}
				fmt.Printf("Copied line %s for %s to clipboard\n", line, service)
			} else {
				fmt.Printf("%s\n", line)
			}
			return nil
		}

		password = []byte(strings.Join(nonEmptyLines, "\n"))

		if qrFlag {
			qr.Generate(string(password), qr.M, os.Stdout)
		}

		if copyFlag {
			err = clipboard.WriteAll(string(password))
			if err != nil {
				return fmt.Errorf("failed to copy password to clipboard: %w", err)
			}
			fmt.Printf("Copied password for %s to clipboard\n", service)
		} else {
			fmt.Println(string(password))
		}

		return nil
	},
}

func init() {
	fmt.Println(VaultPath)
	showCmd.RegisterFlagCompletionFunc("service", AutocompleteServices)
	showCmd.Flags().BoolP("qr", "q", false, "Show QR code of password")
	showCmd.Flags().IntP("line", "l", 0, "Show a specific line of the file")
	showCmd.Flags().BoolP("copy", "c", false, "Copy password to clipboard")
	rootCmd.AddCommand(showCmd)
}
