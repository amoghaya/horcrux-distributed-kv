package server

//Parse raw TCP text → structured command

import "strings"

// Command represents a parsed request
type Command struct {
	Name string
	Args []string
}

// ParseCommand converts raw input into Command struct
func ParseCommand(input string) Command {
	input = strings.TrimSpace(input)
	parts := strings.Split(input, " ")

	return Command{
		Name: strings.ToUpper(parts[0]),
		Args: parts[1:],
	}
}
