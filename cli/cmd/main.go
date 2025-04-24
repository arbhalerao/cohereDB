package main

import (
	"fmt"
	"log"

	"github.com/arbha1erao/cohereDB/cli/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cohere-cli",
		Short: "CLI client to interact with cohereDB",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CLI client to interact with cohereDB.")
		},
	}

	var addr string
	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "a", "", "Address of the database server")
	rootCmd.MarkPersistentFlagRequired("addr")

	rootCmd.AddCommand(commands.GetCmd)
	rootCmd.AddCommand(commands.SetCmd)
	rootCmd.AddCommand(commands.DeleteCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
