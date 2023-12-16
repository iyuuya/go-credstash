package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var deleteVersion string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a key",
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]

		app, err := CreateApp()

		if err != nil {
			log.Fatal(err)
		}

		if err := app.Delete(keyName, deleteVersion); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&deleteVersion, "version", "v", "", "Specify version")
}
