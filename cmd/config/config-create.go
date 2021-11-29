package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tacohole/boardman/util/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var generateConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Generates a config file at" + config.ConfigPath,
	Long:  "",
	Run:   generateConfig,
}

func init() {
	generateConfigCmd.Flags().StringVar(&logLevelConfig, "logLevel", "", "log verbosity")
	generateConfigCmd.Flags().StringVar(&dbUrlConfig, "dbUrl", "", "connection string to your database")
	generateConfigCmd.Flags().StringVar(&apiUrlConfig, "apiUrl", "", "API URL")

	ConfigCmd.AddCommand(generateConfigCmd)
}

func generateConfig(cmd *cobra.Command, args []string) {
	logLevel, err := cmd.Flags().GetString(logLevelConfig)
	if err != nil {
		log.Fatalf("Could not read log level: %s", err)
	}

	dbUrl, err := cmd.Flags().GetString(dbUrlConfig)
	if err != nil {
		log.Fatalf("Could not read DB URL: %s", err)
	}

	apiUrl, err := cmd.Flags().GetString(apiUrlConfig)
	if err != nil {
		log.Fatalf("Could not read API URL: %s", err)
	}

	setConfigVars(logLevel, dbUrl, apiUrl)
	writeConfig()
}

func setConfigVars(logLevel string, dbUrl string, apiUrl string) {
	var config *config.Configuration

	if strings.TrimSpace(logLevel) != "" {
		viper.Set(config.LogLevel, logLevel)
	}

	if strings.TrimSpace(dbUrl) != "" {
		viper.Set(config.DbUrl, dbUrl)
	}

	if strings.TrimSpace(apiUrl) != "" {
		viper.Set(config.ApiUrl, apiUrl)
	}
}

func writeConfig() {
	if _, err := os.Stat(config.ConfigPath); os.IsNotExist(err) {
		err := os.Mkdir(config.ConfigPath, 0555)
		if err != nil {
			fmt.Printf("Failed to create config: %s", err)
		}
	}

	_, err := os.Create(config.ConfigPath + config.ConfigFileName)
	if err != nil {
		fmt.Printf("Failed to create config file: %s", err)
	}

	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("Failed to write config file: %s", err)
	}
}
