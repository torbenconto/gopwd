package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/torbenconto/gopwd/internal/io"
)

var renameCmd = &cobra.Command{
	Use:   "rename [service] [new-service] [flags]",
	Short: "Rename a service",
	Args:  cobra.ExactArgs(2),

	ValidArgsFunction: AutocompleteServices,

	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		newService := args[1]

		// Check if service exists
		if !io.Exists(path.Join(VaultPath, service+".gpg")) {
			return fmt.Errorf("service %s not found", service)
		}

		// Check if new service already exists
		if io.Exists(path.Join(VaultPath, newService+".gpg")) {
			return fmt.Errorf("service %s already exists", newService)
		}

		// Rename service
		err := os.Rename(path.Join(VaultPath, service+".gpg"), path.Join(VaultPath, newService+".gpg"))
		if err != nil {
			return fmt.Errorf("failed to rename service: %v", err)
		}

		fmt.Printf("Service %s renamed to %s\n", service, newService)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
