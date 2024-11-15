package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var revertCmd = &cobra.Command{
	Use:   "revert <id>",
	Short: "Revert a previous action",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		fmt.Printf("Reverting action with ID: %s\n", id)
	},
}
