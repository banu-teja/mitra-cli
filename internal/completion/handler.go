package completion

import (
	"context"
	"fmt"
	"log"

	db "github.com/banu-teja/mitra-cli/db/sqlc"
	"github.com/banu-teja/mitra-cli/internal/system"
	"github.com/banu-teja/mitra-cli/utils"
)

func Handle(question string, config utils.Config, store db.Store) error {
	client := NewClient(config)

	sysInfo := GetSystemInfo()
	systemInfo := FormatSystemInfo(sysInfo)

	shellCommand, err := GetShellCommand(sysInfo.CurrentShell, sysInfo.OperatingSystem)
	if err != nil {
		fmt.Printf("Error getting shell command: %v\n", err)

	}

	prompt := fmt.Sprintf(`%s
%s

Human: %s`, shellCommand, systemInfo, question)

	// println(prompt)

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
