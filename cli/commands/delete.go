package commands

import (
	"fmt"
	"log"

	"github.com/arbha1erao/cohereDB/cli/client"
	"github.com/spf13/cobra"
)

// DeleteCmd represents the "delete" command
var DeleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a value from the database by key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		key := args[0]

		response, err := client.Delete(addr, key)
		if err != nil {
			log.Fatalf("Error deleting key: %v", err)
		}
		fmt.Printf("Response: %s\n", response)
	},
}
