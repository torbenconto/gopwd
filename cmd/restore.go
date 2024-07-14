package cmd

import (
	"github.com/spf13/cobra"

	"github.com/torbenconto/gopwd/internal/backup"
)

var restoreCmd = &cobra.Command{
	Use:   "restore [archivePath] [restoreTo] [flags]",
	Short: "Restore the gopwd vault from an archive",
	Args:  cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		// Restore to provided path or one in config
		archivePath := args[0]
		var restoreTo string
		if len(args) >= 2 {
			restoreTo = args[1]
		} else {
			restoreTo = VaultPath
		}

		err := backup.RestoreGopwdVault(archivePath, restoreTo)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
