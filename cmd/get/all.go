package get

import "github.com/spf13/cobra"

var getAllCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getAll,
}

func init() {
	getAllCmd.Flags().StringVar(&writeTo, "output", "", "output type")

}

func getAll(cmd *cobra.Command, args []string) {
	loadDefaultVariables()
}
