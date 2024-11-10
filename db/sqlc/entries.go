package db

type HistoryEntry struct {
	CommandRequest
	SubCommands []SubCommand
}
