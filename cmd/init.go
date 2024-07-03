package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/torbenconto/gopwd/internal/io"
	util2 "github.com/torbenconto/gopwd/internal/util"
	"path"
)

var vaultPath string

var initCmd = &cobra.Command{
	Use:   "init [gpg-id] [flags]",
	Short: "Initialize gopwd vault at a path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		//check if vaultPath exists
		//if not, create it
		if !io.Exists(vaultPath) {
			err := io.CreateDir(vaultPath)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("vault already exists at %s", vaultPath)
		}

		err := io.WriteToFile(path.Join(GopwdPath, ".gopwd.yaml"), []byte(fmt.Sprintf("vaultPath: %s", vaultPath)))
		if err != nil {
			return err
		}

		//create .gpg-id file in vaultPath
		_, err = io.CreateFile(path.Join(vaultPath, ".gpg-id"))
		if err != nil {
			return err
		}

		//write gpg-id to .gpg-id file
		err = io.WriteToFile(path.Join(vaultPath, ".gpg-id"), []byte(args[0]))
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
		if !io.Exists(GopwdPath) {
			err := io.CreateDir(GopwdPath)
			if err != nil {
				panic(err)
			}
		}

		if !io.Exists(path.Join(GopwdPath, ".gopwd.yaml")) {
			_, err := io.CreateFile(path.Join(GopwdPath, ".gopwd.yaml"))
			if err != nil {
				panic(err)
			}
		}

		err := util2.LoadConfig(configFilePath)

		if err != nil {
			panic(err)
		}
	}

	initCmd.Flags().StringVarP(&vaultPath, "path", "p", path.Join(GopwdPath, "vault"), "path to gopwd vault")
	rootCmd.AddCommand(initCmd)
}
