package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/torbenconto/gopwd/internal/io"
	"github.com/torbenconto/gopwd/internal/util"
)

var (
	GopwdPath      = path.Join(io.GetHomeDir(), ".gopwd")
	configFilePath string
	VaultPath      string
)

// AutocompleteServices provides autocompletion for service names.
func AutocompleteServices(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	services, err := io.ListServices(VaultPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return services, cobra.ShellCompDirectiveNoFileComp
}

var rootCmd = &cobra.Command{
	Use:   "gopwd",
	Short: "gopwd is a password manager",
	Long:  "gopwd is a password manager that stores your passwords in an encrypted file",

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(strings.Split(VaultPath, "/")[len(strings.Split(VaultPath, "/"))-1])
		util.PrintVaultStructure(VaultPath)
		return nil
	},
}

func Execute() {
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Name() == "init" || cmd.Name() == "completion" || cmd.Name() == "__complete" {
			return
		}

		// check if .goPwd directory exists
		// if not, tell the user to run gopwd init
		if !io.Exists(GopwdPath) {
			panic("gopwd not initialized. Run 'gopwd init' to initialize")
		}

		// Load config
		util.LoadConfig(configFilePath)

		// check if vaultPath exists
		// if not, tell the user to run gopwd init
		if !io.Exists(viper.Get("vaultPath").(string)) {
			panic("gopwd vault not initialized. Run 'gopwd init' to initialize")
		}

		VaultPath = viper.Get("vaultPath").(string)
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.CompletionOptions.DisableDescriptions = true

	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", path.Join(GopwdPath, ".gopwd.yaml"), "config file (default is $HOME/.gopwd/.gopwd.yaml)")
	rootCmd.Execute()
}
