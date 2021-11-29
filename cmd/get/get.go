package get

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  "",
	Args:  cobra.MaximumNArgs(1),
}

var verbose bool
var source string
var writeTo string

func loadDefaultVariables() {
	verbose = viper.GetBool("verbose")
	writeTo = viper.GetString("writeTo")
}