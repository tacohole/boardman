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
	ApiUrlEnvironmentName     = "API_URL"
)

type Configuration struct {
	LogLevel string
	DbUrl    string
	ApiUrl   string
}

var config *Configuration

func getConfig() *Configuration {
	if config == nil {
		config = &Configuration{}
	}

	config.LogLevel = viper.GetString("LOG_LEVEL")
	if config.LogLevel == "" {
		config.LogLevel = "debug"
	}

	config.DbUrl = viper.GetString(DbUrlEnvironmentName)
	if config.DbUrl == "" {
		log.Println("Unknown database URL")
	}

	config.ApiUrl = viper.GetString(ApiUrlEnvironmentName)
	if config.ApiUrl == "" {
		log.Printf("Unknown API URL")
	}

	return config
}
