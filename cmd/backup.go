package cmd

import (
	"github.com/spf13/cobra"

	"github.com/torbenconto/gopwd/internal/backup"
)

var backupCmd = &cobra.Command{
	Use:   "backup [toDir] [flags]",
	Short: "Backup the gopwd vault",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		to := args[0]

		err := backup.ArchiveGopwdVault(VaultPath, to)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
