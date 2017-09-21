package commander

import (
	"fmt"
	"github.com/shomali11/proper"
	"regexp"
	"strings"
)

const (
	empty            = ""
	space            = " "
	ignoreCase       = "(?i)"
	parameterPattern = "<\\S+>"
	spacePattern     = "\\s*"
	wordPattern      = "\\b(\\S+)?\\b"
	boundaryFormat   = "\\b%s\\b"
)

var parameterRegex *regexp.Regexp

func init() {
	parameterRegex = regexp.MustCompile(parameterPattern)
}

// NewCommand creates a new Command object from the format passed in
func NewCommand(format string) *Command {
	expression := compile(format)
	return &Command{format: format, expression: expression}
}

// Token represents the Token object
type Token struct {
	Word        string
	IsParameter bool
}

// Command represents the Command object
type Command struct {
	format     string
	expression *regexp.Regexp
}

// Match takes in the command and the text received, attempts to find the pattern and extract the parameters
func (c *Command) Match(text string) (*proper.Properties, bool) {
	if c.expression == nil {
		return nil, false
	}

	result := strings.TrimSpace(c.expression.FindString(text))
	if len(result) == 0 {
		return nil, false
	}

	parameters := make(map[string]string)
	commandTokens := strings.Split(c.format, space)
	resultTokens := strings.Split(result, space)

	for i, resultToken := range resultTokens {
		commandToken := commandTokens[i]
		if !isParameter(commandToken) {
			continue
		}

		parameters[commandToken[1:len(commandToken)-1]] = resultToken
	}
	return proper.NewProperties(parameters), true
}

// Tokenize returns Command info as tokens
func (c *Command) Tokenize() []*Token {
	words := strings.Split(c.format, space)
	tokens := make([]*Token, len(words))
	for i, word := range words {
		if isParameter(word) {
			tokens[i] = &Token{Word: word[1 : len(word)-1], IsParameter: true}
		} else {
			tokens[i] = &Token{Word: word, IsParameter: false}
		}
	}
	return tokens
}

func isParameter(text string) bool {
	return parameterRegex.MatchString(text)
}

func compile(commandFormat string) *regexp.Regexp {
	commandFormat = strings.TrimSpace(commandFormat)
	tokens := strings.Split(commandFormat, space)
	pattern := empty
	for _, token := range tokens {
		if len(token) == 0 {
			continue
		}

		if isParameter(token) {
			pattern += wordPattern
		} else {
			pattern += fmt.Sprintf(boundaryFormat, token)
		}
		pattern += spacePattern
	}

	if len(pattern) == 0 {
		return nil
	}
	return regexp.MustCompile(ignoreCase + pattern)
}
