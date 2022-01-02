package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	// TODO config vars
	ConfigPath                = "~/"
	ConfigFileNameNoExtension = "boardman-config"
	ConfigFileName            = "boardman-config.env"
	DbUrlEnvironmentName      = "DATABASE_URL"
	DbTimeout                 = "DB_TIMEOUT"
)

type Configuration struct {
	DbUrl     string
	DbTimeout time.Duration
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

	config.DbTimeout = viper.GetDuration(fmt.Sprint(DbTimeout))
	if config.DbTimeout == 0 {
		log.Println("No database timeout set")
	}

	return config
}
