package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/iyuuya/go-credstash/internal/credstash/app"
)

var rootCmd = &cobra.Command{
	Use:   "credstash",
	Short: "credstash cli",
}

var endpoint string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "Endpoint for dynamodb-local")
	rootCmd.AddCommand(
		getCmd,
		putCmd,
		listCmd,
		deleteCmd,
		setupCmd,
	)
}

func CreateApp() (*app.App, error) {
	return app.NewApp(endpoint)
}
