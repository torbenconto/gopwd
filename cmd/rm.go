package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/pkg/prompts"
	"github.com/torbenconto/gopwd/pkg/util"
	"os"
	"path"
)

var rmCmd = &cobra.Command{
	Use:   "rm [service] [flags]",
	Short: "Remove a password for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		servicePath := path.Join(VaultPath, service) + ".gpg"
		dirPath := path.Dir(servicePath)

		if util.Exists(servicePath) {
			action, err := prompts.ConfirmAction()
			if err != nil {
				return fmt.Errorf("failed to confirm action: %v", err)
			}
			if action {
				err := util.RemoveFile(servicePath)
				if err != nil {
					return fmt.Errorf("failed to remove service: %s, error: %v", service, err)
				}

				// Check if the directory is empty and not the root vault directory
				isEmpty, err := util.IsDirEmpty(dirPath)
				if err != nil {
					return fmt.Errorf("failed to check if directory is empty: %v", err)
				}
				if isEmpty && dirPath != VaultPath {
					err := os.Remove(dirPath)
					if err != nil {
						return fmt.Errorf("failed to remove directory: %s, error: %v", dirPath, err)
					}
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
