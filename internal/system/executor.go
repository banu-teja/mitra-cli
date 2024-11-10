package system

import (
	"context"
	"log"
	"os/exec"
	"strings"

	db "github.com/banu-teja/mitra-cli/db/sqlc"
	"github.com/banu-teja/mitra-cli/utils"
)

func ExecuteAndStoreSubCommands(question, response string, store db.Store) {
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
		log.Printf("Executing command: %s", subCommand)
		output, status := executeCommand(subCommand)

		dbStatus := mapStatus(status)

		subCmdReq := db.CreateSubCommandParams{
			RequestID:      cmd.ID,
			Command:        subCommand,
			ExecutionOrder: int64(i + 1),
			CommandOutput:  output,
			CommandStatus:  dbStatus,
		}

		_, err := store.CreateSubCommand(ctx, subCmdReq)
		if err != nil {
			log.Fatal("cannot insert subcommand: ", err, " output: ", output)
			continue
		}
		log.Printf("Command executed: [%s] Status: %s Output: %s", subCommand, status, output)
	}
}
func mapStatus(status string) string {
	switch status {
	case "success":
		return "success"
	case "error":
		return "failure"
	default:
		return "in_progress"
	}
}

func executeCommand(command string) (string, string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "Empty command", "error"
	}

	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()

	outputStr := string(output)
	if len(outputStr) > 500 { // Adjust this number as needed
		outputStr = outputStr[:250] + "\n...\n" + outputStr[len(outputStr)-250:]
	}

	if err != nil {
		log.Printf("Command execution failed: %v", err)
		return outputStr, "error"
	}

	return outputStr, "success"
}
