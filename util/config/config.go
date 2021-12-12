package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	// TODO config vars
	ConfigPath                = "etc/boardman"
	ConfigFileNameNoExtension = "boardman-config"
	ConfigFileName            = "boardman-config.json"
	DbUrlEnvironmentName      = "DB_URL"
)

type Configuration struct {
	DbUrl string
}

var config *Configuration

func getConfig() *Configuration {
	if config == nil {
		config = &Configuration{}
	}

	config.DbUrl = viper.GetString(DbUrlEnvironmentName)
	if config.DbUrl == "" {
		log.Println("Unknown database URL")
	}

	return config
}
