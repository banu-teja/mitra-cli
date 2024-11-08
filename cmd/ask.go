package cmd

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/spf13/cobra"
// )

// var rootCmd1 = &cobra.Command{
// 	Use:   "mitra [question]",
// 	Short: "Mitra CLI tool",
// 	Long:  `Mitra is a CLI tool for asking questions and managing interactions.`,
// 	// This ensures that any unknown command is treated as part of the question
// 	Args: cobra.ArbitraryArgs,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		if len(args) > 0 {
// 			question := strings.Join(args, " ")
// 			// Process the question
// 			fmt.Printf("Processing question: %s\n", question)
// 			return nil
// 		}
// 		// If no args, show help
// 		return cmd.Help()
// 	},
// }

// // Subcommands
// var revertCmd = &cobra.Command{
// 	Use:   "revert <id>",
// 	Short: "Revert a previous action",
// 	Args:  cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		id := args[0]
// 		fmt.Printf("Reverting action with ID: %s\n", id)
// 	},
// }

// var historyCmd = &cobra.Command{
// 	Use:   "history",
// 	Short: "Display interaction history",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Displaying history...")
// 	},
// }

// var checkCmd = &cobra.Command{
// 	Use:   "check <id>",
// 	Short: "Check information about a specific interaction",
// 	Args:  cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		id := args[0]
// 		fmt.Printf("Checking interaction with ID: %s\n", id)
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(revertCmd)
// 	rootCmd.AddCommand(historyCmd)
// 	rootCmd.AddCommand(checkCmd)
// }

// func Execute1() error {
// 	return rootCmd1.Execute()
// }
