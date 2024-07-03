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

		// Check if service exists
		if !io.Exists(path.Join(VaultPath, service+".gpg")) {
			return fmt.Errorf("service %s not found", service)
		}

		// Get password and decrypt
		file, err := io.ReadFromFile(path.Join(VaultPath, service+".gpg"))
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

		line, _ := cmd.Flags().GetInt("line")
		lines := strings.Split(string(password), "\n")
		if line > 0 {
			line := lines[line-1]
			if line == "" {
				return fmt.Errorf("line %s is empty", line)
			}
			if qrFlag, _ := cmd.Flags().GetBool("qr"); qrFlag {
				qr.Generate(string(password), qr.M, os.Stdout)
			}
			if copyFlag, _ := cmd.Flags().GetBool("copy"); copyFlag {
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

		if qrFlag, _ := cmd.Flags().GetBool("qr"); qrFlag {
			qr.Generate(string(password), qr.M, os.Stdout)
		}

		if copyFlag, _ := cmd.Flags().GetBool("copy"); copyFlag {
			err = clipboard.WriteAll(string(password))
			if err != nil {
				return fmt.Errorf("failed to copy password to clipboard: %w", err)
			}
			fmt.Printf("Copied password for %s to clipboard\n", service)
		} else {
			fmt.Printf("%s\n", password)
		}

		return nil
	},
}

func init() {
	showCmd.Flags().BoolP("qr", "q", false, "Show QR code of password")
	showCmd.Flags().IntP("line", "l", 1, "Show a specific line of the file")
	showCmd.Flags().BoolP("copy", "c", false, "Copy password to clipboard")
	rootCmd.AddCommand(showCmd)
}
