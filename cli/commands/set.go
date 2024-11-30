package commands

import (
	"fmt"
	"log"

	"github.com/arbha1erao/cohereDB/cli/client"
	"github.com/spf13/cobra"
)

// SetCmd represents the "set" command
var SetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a value in the database",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		key := args[0]
		value := args[1]

		response, err := client.Set(addr, key, value)
		if err != nil {
			log.Fatalf("Error setting value: %v", err)
		}
		fmt.Printf("Response: %s\n", response)
	},
}
