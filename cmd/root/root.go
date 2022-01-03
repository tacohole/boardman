/*
Copyright © 2021 Troy Coll troy.coll@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package root

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	genConfig "github.com/tacohole/boardman/cmd/config"
	"github.com/tacohole/boardman/cmd/get"
)

// var (
// 	// cfgFile   string
// 	// dbTimeout time.Duration
// 	// verbose   bool
// )

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "boardman",
	Short: "A data loader for NBA player/team stats",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", config.ConfigPath+config.ConfigFileName))
	// RootCmd.PersistentFlags().DurationVar(&dbTimeout, "dbTimeout", 60*time.Second, "database timeout default is 60s, use 60m for minutes or 60h for hours")
	// RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, fmt.Sprintln("Include additional logging information"))

	RootCmd.AddCommand(get.GetCmd)
	RootCmd.AddCommand(genConfig.ConfigCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	godotenv.Load("boardman-config.env")

	// viper.SetDefault(config.DbUrlEnvironmentName, "postgresql://user:secret@localhost:5432")

	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Search config in config path with name ".boardman" (without extension).
	// 	viper.AddConfigPath(config.ConfigPath)
	// 	viper.SetConfigName(config.ConfigFileNameNoExtension)
	// }

	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// } else {
	// 	fmt.Printf("There was a problem with config file %s and it has been ignored: %s", cfgFile, err)
	// }

	// if err := viper.BindPFlag("verbose", RootCmd.Flags().Lookup("verbose")); err != nil {
	// 	log.Println(err.Error())
	// }

}
