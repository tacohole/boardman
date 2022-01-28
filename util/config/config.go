package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	// TODO config vars
	ConfigPath                = "~/"
	ConfigFileNameNoExtension = "boardman-config"
	ConfigFileName            = "boardman-config.env"
	DbUrlEnvironmentName      = "DATABASE_URL"
	DbTimeout                 = "DB_TIMEOUT"
	Verbose                   = "VERBOSE"
)

type Configuration struct {
	DbUrl     string
	DbTimeout string
	Verbose   string
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

	config.DbTimeout = viper.GetString(DbTimeout)
	if config.DbTimeout == "" {
		log.Println("No database timeout set")
	}

	config.Verbose = viper.GetString(Verbose)
	if config.Verbose == "false" {
		log.Println("additional logging disabled")
	}

	return config
}
