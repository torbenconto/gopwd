package util

import (
	"github.com/spf13/viper"
	"path"
)

func LoadConfig(configFilePath string) error {
	if configFilePath == "" {
		// Load config from default location
		viper.SetConfigName(".gopwd")
		viper.AddConfigPath(path.Join(GetHomeDir(), ".gopwd"))
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	} else {
		// Load config from specified location
		viper.SetConfigFile(configFilePath)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	}
	return nil
}
