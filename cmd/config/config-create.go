package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	generateConfigCmd.Flags().StringVar(&dbUrlConfig, "dbUrl", "", "connection string to your database")
	generateConfigCmd.Flags().StringVar(&dbTimeOutConfig, "dbTimeout", "", "database timeout setting")

	ConfigCmd.AddCommand(generateConfigCmd)
}

func generateConfig(cmd *cobra.Command, args []string) {

	dbUrl, err := cmd.Flags().GetString(dbUrlConfig)
	if err != nil {
		log.Fatalf("Could not read DB URL: %s", err)
	}

	dbTimeOut, err := cmd.Flags().GetString(dbTimeOutConfig)
	if err != nil {
		log.Fatalf("Could not read database timeout: %s", err)
	}

	setConfigVars(dbUrl, dbTimeOut)
	writeConfig()
}

func setConfigVars(dbUrl string, dbTimeOut string) {
	var config *config.Configuration

	if strings.TrimSpace(dbUrl) != "" {
		viper.Set(config.DbUrl, dbUrl)
	}

	if strings.TrimSpace(dbTimeOut) != "" {
		timeDur, err := strconv.Atoi(dbTimeOut)
		if err != nil {
			log.Printf("%s", err)
		}
		viper.Set(config.DbTimeout, timeDur)
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
