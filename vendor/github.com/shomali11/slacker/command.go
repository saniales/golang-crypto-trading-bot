package slacker

import (
	"github.com/shomali11/commander"
	"github.com/shomali11/proper"
)

// NewBotCommand creates a new bot command object
func NewBotCommand(usage string, description string, handler func(request *Request, response ResponseWriter)) *BotCommand {
	command := commander.NewCommand(usage)
	return &BotCommand{usage: usage, description: description, handler: handler, command: command}
}

// BotCommand structure contains the bot's command, description and handler
type BotCommand struct {
	usage       string
	description string
	handler     func(request *Request, response ResponseWriter)
	command     *commander.Command
}

// Match determines whether the bot should respond based on the text received
func (c *BotCommand) Match(text string) (*proper.Properties, bool) {
	return c.command.Match(text)
}

// Tokenize returns the command format's tokens
func (c *BotCommand) Tokenize() []*commander.Token {
	return c.command.Tokenize()
}

// Execute executes the handler logic
func (c *BotCommand) Execute(request *Request, response ResponseWriter) {
	c.handler(request, response)
}
