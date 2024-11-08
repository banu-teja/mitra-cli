package utils

import (
	"regexp"
	"strings"
)

var commandStartRegex = regexp.MustCompile(`^\w+\s`)

func ParseCommands(input string) []string {
	// Split the input into lines
	lines := strings.Split(strings.TrimSpace(input), "\n")

	var commands []string
	var currentCommand strings.Builder

	for _, line := range lines {
		// Trim leading and trailing whitespace
		line = strings.TrimSpace(line)

		// Check if the line is a new command (starts with a word followed by a space)
		if commandStartRegex.MatchString(line) {
			// If we have a current command, add it to the list
			if currentCommand.Len() > 0 {
				commands = append(commands, strings.TrimSpace(currentCommand.String()))
				currentCommand.Reset()
			}
			// Start a new command
			currentCommand.WriteString(line)
		} else {
			// Continue the current command
			if currentCommand.Len() > 0 {
				currentCommand.WriteString("\n")
			}
			currentCommand.WriteString(line)
		}
	}

	// Add the last command if there is one
	if currentCommand.Len() > 0 {
		commands = append(commands, strings.TrimSpace(currentCommand.String()))
	}

	return commands
}
