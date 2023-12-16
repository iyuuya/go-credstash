package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var getVersion string

var getCmd = &cobra.Command{
	Use:   "get [KEY_NAME]",
	Short: "Show a value for key name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]

		app, err := CreateApp()
		if err != nil {
			log.Fatal(err)
		}

		item, err := app.Get(keyName, nil, getVersion)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(item)
	},
}

func init() {
	getCmd.Flags().StringVarP(&getVersion, "version", "v", "", "Show a value for key name with their versions")
}
