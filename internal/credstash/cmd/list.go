package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var listVersion bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all stored keys",
	Run: func(cmd *cobra.Command, _ []string) {
		app, err := CreateApp()
		if err != nil {
			log.Fatal(err)
		}
		items, err := app.DynamoDB.List()
		if err != nil {
			log.Fatal(err)
		}

		if listVersion {
			for _, i := range items {
				name := ""
				version := ""
				if n := i.GetName(); n != nil {
					name = *n
				}
				if v := i.GetVersion(); v != nil {
					version = *v
				}
				fmt.Printf("%s --version: %s\n", name, version)
			}
		} else {
			for _, i := range items {
				name := ""
				if n := i.GetName(); n != nil {
					name = *n
				}
				fmt.Printf("%s\n", name)
			}
		}
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listVersion, "version", "v", false, "Show all stored keys with their versions")
}
