package config

import (
	"github.com/tacohole/boardman/util/config"

	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage config variables for an environment",
	Long:  "Generate or source a config file and add env variables:" + config.ConfigPath + config.ConfigFileName,
	Args:  cobra.MinimumNArgs(1),
}

// TODO config vars
var logLevelConfig string
var dbUrlConfig string
var apiUrlConfig string
