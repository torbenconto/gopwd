//go:build windows

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/api"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Command palette for interacting with the API",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of the API server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Status command is not implemented for Windows.")
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		portFlag, _ := cmd.Flags().GetString("port")

		// Start the API server
		fmt.Println("Starting API server on port:", portFlag)
		err := api.Run(GopwdPath, VaultPath, ":"+portFlag, filepath.Join(GopwdPath, "cert.pem"), filepath.Join(GopwdPath, "key.pem"))
		if err != nil {
			fmt.Println("Error starting API server:", err)
			os.Exit(1)
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the API server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stop command is not implemented for Windows.")
	},
}

func init() {
	apiCmd.AddCommand(statusCmd)
	apiCmd.AddCommand(downCmd)
	apiCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("port", "p", "8076", "Port to run the API server on")

	rootCmd.AddCommand(apiCmd)
}
