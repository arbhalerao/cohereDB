package commands

import (
	"fmt"
	"log"

	"github.com/arbha1erao/cohereDB/cli/client"
	"github.com/spf13/cobra"
)

// GetCmd represents the "get" command
var GetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a value from the database by key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		key := args[0]

		response, err := client.Get(addr, key)
		if err != nil {
			log.Fatalf("Error getting value: %v", err)
		}
		fmt.Printf("Value for '%s': %s\n", key, response)
	},
}
