package get

import "github.com/spf13/cobra"

var getOneCmd = &cobra.Command{
	Short: "",
	Long:  "",
	Run:   getOne,
}

func init() {
	getOneCmd.Flags().StringVar(&source, "source", "", "name of source")
	getOneCmd.Flags().StringVar(&writeTo, "output", "", "output type")

	getOneCmd.MarkFlagRequired("source")
	getOneCmd.MarkFlagRequired("writeTo")

}

func getOne(cmd *cobra.Command, args []string) {
	loadDefaultVariables()
}
