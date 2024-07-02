package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/pkg/util"
	"path"
)

var vaultPath string

var initCmd = &cobra.Command{
	Use:   "init [gpg-id] [flags]",
	Short: "Initialize gopwd vault at a path",
	Long:  "Initialize gopwd vault at a path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		//check if vaultPath exists
		//if not, create it
		if !util.Exists(vaultPath) {
			err := util.CreateDir(vaultPath)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("vault already exists at %s", vaultPath)
		}

		err := util.WriteToFile(path.Join(GopwdPath, ".gopwd.yaml"), []byte(fmt.Sprintf("vaultPath: %s", vaultPath)))
		if err != nil {
			return err
		}

		//create .gpg-id file in vaultPath
		_, err = util.CreateFile(path.Join(vaultPath, ".gpg-id"))
		if err != nil {
			return err
		}

		//write gpg-id to .gpg-id file
		err = util.WriteToFile(path.Join(vaultPath, ".gpg-id"), []byte(args[0]))
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		// check if .goPwd directory exists
		// check if .goPwd/.goPwd.yaml exists
		// if not, create it
		if !util.Exists(GopwdPath) {
			err := util.CreateDir(GopwdPath)
			if err != nil {
				panic(err)
			}
		}

		if !util.Exists(path.Join(GopwdPath, ".gopwd.yaml")) {
			_, err := util.CreateFile(path.Join(GopwdPath, ".gopwd.yaml"))
			if err != nil {
				panic(err)
			}
		}

		err := util.LoadConfig(configFilePath)

		if err != nil {
			panic(err)
		}
	}

	initCmd.Flags().StringVarP(&vaultPath, "path", "p", path.Join(GopwdPath, "vault"), "path to gopwd vault")
	rootCmd.AddCommand(initCmd)
}
