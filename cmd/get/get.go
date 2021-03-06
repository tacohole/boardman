package get

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "",
	Long:  "",
	Args:  cobra.MaximumNArgs(1),
}

var verbose bool

func loadDefaultVariables() {
	verbose = viper.GetBool("VERBOSE")
}
