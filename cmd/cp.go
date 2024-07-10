package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/termio"
	"path"
)

var cpCmd = &cobra.Command{
	Use:   "cp [service] [new-service] [flags]",
	Short: "Copy a password for a service",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		newService := args[1]

		servicePath := path.Join(VaultPath, service) + ".gpg"
		newServicePath := path.Join(VaultPath, newService) + ".gpg"

		// Check if service exists
		if !io.Exists(servicePath) {
			return fmt.Errorf("service %s not found", service)
		}

		// Check if new service already exists
		if io.Exists(newServicePath) {
			fmt.Println("There is already a service at this name, would you like to overwrite it? This action cannot be undone.")
			if confirm, _ := termio.ConfirmAction(); confirm {
				err := io.CopyFile(servicePath, newServicePath)
				if err != nil {
					return fmt.Errorf("failed to copy file: %v", err)
				}
			} else {
				fmt.Println("Aborted")
			}
		} else {
			err := io.CopyFile(servicePath, newServicePath)
			if err != nil {
				return fmt.Errorf("failed to copy file: %v", err)
			}
			fmt.Printf("Copied %s to %s\n", service, newService)
		}

		return nil
	},
}

func init() {
	cpCmd.RegisterFlagCompletionFunc("service", AutocompleteServices)
	cpCmd.RegisterFlagCompletionFunc("new-service", AutocompleteServices)

	rootCmd.AddCommand(cpCmd)
}
