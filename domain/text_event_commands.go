package domain

// TextEventCommands holds a map of input strings to actions that the
// system should take in response to the parsed text
type TextEventCommands struct {
	Commands []TextEventCommand
}

// TextEventCommand represents a known command the system should respond to
type TextEventCommand struct {
	Command      string
	ReturnOutput bool   `yaml:"return-output"`
	ActionFormat string `yaml:"action-format"`
	Action       string
}
