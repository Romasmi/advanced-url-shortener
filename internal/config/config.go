package config

import (
	"fmt"
	"shorturl/internal/utils"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Port uint
}

const ConfigurationFilename = ".env"

func LoadConfig(config *viper.Viper) {
	utils.FatalIfFileDoesNotExists(ConfigurationFilename, fmt.Sprintf("config file %v does not exist", ConfigurationFilename))

	config = viper.New()
	config.SetConfigFile(ConfigurationFilename)

	err := viper.ReadInConfig()
	if err != nil {

	}
}
