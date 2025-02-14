package completion

import (
	"context"
	"fmt"
	"log"
	"strings"

	db "github.com/banu-teja/mitra-cli/db/sqlc"
	"github.com/banu-teja/mitra-cli/internal/system"
	"github.com/banu-teja/mitra-cli/utils"
)

func formatHistory(history []db.GetLastNEntriesRow) string {
	var result strings.Builder
	currentRequestID := int64(-1)

	result.WriteString("\n=== Previous History ===\n")
	for _, entry := range history {
		if entry.RequestID != currentRequestID {
			// New request, print input content
			if currentRequestID != -1 {
				result.WriteString("\n") // Add extra newline between requests
			}
			result.WriteString(fmt.Sprintf("Input: %s\n", entry.InputContent))
			result.WriteString("Commands:\n")
			currentRequestID = entry.RequestID
		}

		// Print command details if present
		if entry.SubcommandID.Valid {
			result.WriteString(fmt.Sprintf("  - Command: %s\n", entry.Command.String))
			result.WriteString(fmt.Sprintf("    Output: %s", entry.CommandOutput.String))
			result.WriteString(fmt.Sprintf("    Status: %s\n", entry.CommandStatus.String))
		}
	}
	result.WriteString("=== End of Previous History ===\n")
	return result.String()
}

func Handle(question string, config utils.Config, store db.Store) error {
	client := NewGoogleAIClient(config)

	sysInfo := GetSystemInfo()
	systemInfo := FormatSystemInfo(sysInfo)

	shellCommand, err := GetShellCommand(sysInfo.CurrentShell, sysInfo.OperatingSystem)
	if err != nil {
		fmt.Printf("Error getting shell command: %v\n", err)

	}

	history, err := store.GetLastNEntries(context.Background(), 3)
	if err != nil {
		fmt.Printf("Error fetching history: %v\n", err)
	}

	formattedHistory := formatHistory(history)

	prompt := fmt.Sprintf(`%s
%s
%s

Human: %s`, shellCommand, systemInfo, formattedHistory, question)

	println(prompt)

	commands, err := client.Send(context.Background(), prompt)

	if err != nil {
		return err
	}

	system.ExecuteAndStoreSubCommands(question, commands, store)

	return nil
}

func splitIntoSubCommands(question, response string, store db.Store) {
	cmdReq := db.CreateCommandRequestParams{
		InputContent: question,
		CommandType:  "question",
	}

	ctx := context.Background()

	cmd, err := store.CreateCommandRequest(ctx, cmdReq)
	if err != nil {
		log.Fatal("cannot insert command: ", err)
	}
	subCommands := utils.ParseCommands(response)

	for i, subCommand := range subCommands {
		subCmdReq := db.CreateSubCommandParams{
			RequestID:      cmd.ID,
			Command:        subCommand,
			ExecutionOrder: int64(i + 1),
			CommandOutput:  "hard",
			CommandStatus:  "success",
		}
		subCmd, err := store.CreateSubCommand(ctx, subCmdReq)
		if err != nil {
			log.Fatal("cannot insert command: ", err)
		}
		log.Printf("SubCommand created: %+v", subCmd)
	}
}
