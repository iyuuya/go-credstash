package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup credstash repository on DynamoDB",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := CreateApp()
		if err != nil {
			log.Fatal(err)
		}
		app.DynamoDB.Setup()
	},
}
