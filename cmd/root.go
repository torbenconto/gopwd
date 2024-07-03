package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/util"
	"path"
)

var (
	GopwdPath      = path.Join(io.GetHomeDir(), ".gopwd")
	configFilePath string
	VaultPath      string
)

var rootCmd = &cobra.Command{
	Use:   "gopwd",
	Short: "gopwd is a password manager",
	Long:  "gopwd is a password manager that stores your passwords in an encrypted file",

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("vault")
		util.PrintVaultStructure(VaultPath)
		return nil
	},
}

func Execute() {
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "init" {
			return
		}

		// check if .goPwd directory exists
		// if not, tell the user to run gopwd init
		if !io.Exists(GopwdPath) {
			panic("gopwd not initialized. Run 'gopwd init' to initialize")
		}

		// Load config
		err := util.LoadConfig(configFilePath)
		if err != nil {
			panic(err)
		}

		// check if vaultPath exists
		// if not, tell the user to run gopwd init
		if !io.Exists(viper.Get("vaultPath").(string)) {
			panic("gopwd vault not initialized. Run 'gopwd init' to initialize")
		}

		VaultPath = viper.Get("vaultPath").(string)
	}

	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", path.Join(GopwdPath, ".gopwd.yaml"), "config file (default is $HOME/.gopwd/.gopwd.yaml)")
	rootCmd.Execute()
}
