package cmd

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  "",
	Run: ,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func get() {

}
