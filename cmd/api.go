package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/api"
	"github.com/torbenconto/gopwd/internal/io"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

var pidFile = filepath.Join(GopwdPath, "gopwd.pid")

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Command palate for interacting with the API",

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of the API server",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if PID file exists
		pidData, err := io.ReadFile(pidFile)
		if err != nil {
			fmt.Println("Daemon is not running")
			return
		}

		// Parse PID from file
		pid, err := strconv.Atoi(string(pidData))
		if err != nil {
			fmt.Println("Invalid PID file content")
			return
		}

		// Check if process with PID is running
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println("Error finding process:", err)
			return
		}

		// Send signal 0 to check process existence
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			fmt.Println("Daemon is not running")
			return
		}

		fmt.Println("Daemon is running with PID:", pid)
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		portFlag, _ := cmd.Flags().GetString("port")

		// Remove log file if it exists
		logFile := filepath.Join(GopwdPath, "gopwd.log")
		if io.Exists(logFile) {
			err := os.Remove(logFile)
			if err != nil {
				fmt.Println("Error removing log file:", err)
				return
			}
		}
		// Start the API server as a daemon
		api.RunDaemon(GopwdPath, VaultPath, ":"+string(portFlag), os.Args)
		fmt.Println("API server started as daemon")
		fmt.Println("NEVER EVER EVER EXPOSE THIS TO THE INTERNET, IT IS NOT SECURE (LAN ONLY)")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the API server",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if PID file exists
		pidData, err := io.ReadFile(pidFile)
		if err != nil {
			fmt.Println("Daemon is not running")
			return
		}

		// Parse PID from file
		pid, err := strconv.Atoi(string(pidData))
		if err != nil {
			fmt.Println("Invalid PID file content")
			return
		}

		// Check if process with PID is running
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println("Error finding process:", err)
			return
		}

		// Send signal to stop the process
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			fmt.Println("Daemon is not running")
			return
		}

		// Remove PID file after stopping the process
		err = os.Remove(pidFile)
		if err != nil {
			fmt.Println("Error removing PID file:", err)
			return
		}

		fmt.Println("Daemon stopped")
	},
}

func init() {
	apiCmd.AddCommand(statusCmd)
	apiCmd.AddCommand(downCmd)
	apiCmd.AddCommand(upCmd)
	upCmd.Flags().StringP("port", "p", "8076", "Port to run the API server on")

	rootCmd.AddCommand(apiCmd)
}
