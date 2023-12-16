package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var putKmsKeyId string

var putCmd = &cobra.Command{
	Use:   "put [KEY_NAME]",
	Short: "Put a value for key name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]
		fmt.Print("secret value> ")
		v, err := readline()
		if err != nil {
			log.Fatal(err)
		}

		app, err := CreateApp()
		if err != nil {
			log.Fatal(err)
		}

		err = app.Put(keyName, v, putKmsKeyId, nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	putCmd.Flags().StringVarP(&putKmsKeyId, "kms_key_id", "k", "", "the KMS key-id of the master key to use. Defaults to alias/credstash")
}

func readline() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return s[:len(s)-1], nil
}
