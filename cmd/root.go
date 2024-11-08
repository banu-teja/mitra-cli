package cmd

import (
	"fmt"
	"os"
	"strings"

	db "github.com/banu-teja/mitra-cli/db/sqlc"
	"github.com/banu-teja/mitra-cli/internal/completion"
	"github.com/banu-teja/mitra-cli/utils"
	"github.com/spf13/cobra"
)

func Execute(config utils.Config, store db.Store) {
	var rootCmd = &cobra.Command{
		Use:   "mitra [question]",
		Short: "Mitra CLI tool",
		Long:  `Mitra is a CLI tool for asking questions and managing interactions.`,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				questionText := strings.Join(args, " ")
				return completion.Handle(questionText, config, store)
			}
			return cmd.Help()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
