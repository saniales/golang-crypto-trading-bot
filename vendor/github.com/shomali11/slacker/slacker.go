package slacker

import (
	"context"
	"errors"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
	"strings"
)

const (
	space               = " "
	dash                = "-"
	newLine             = "\n"
	invalidToken        = "Invalid token"
	helpCommand         = "help"
	directChannelMarker = "D"
	userMentionFormat   = "<@%s>"
	codeMessageFormat   = "`%s`"
	boldMessageFormat   = "*%s*"
	italicMessageFormat = "_%s_"
)

// NewClient creates a new client using the Slack API
func NewClient(token string) *Slacker {
	client := slack.New(token)
	rtm := client.NewRTM()

	slacker := &Slacker{
		Client: client,
		rtm:    rtm,
	}

	slacker.Command(helpCommand, helpCommand, func(request *Request, response ResponseWriter) {
		helpMessage := empty
		for _, command := range slacker.botCommands {
			tokens := command.Tokenize()
			for _, token := range tokens {
				if token.IsParameter {
					helpMessage += fmt.Sprintf(codeMessageFormat, token.Word) + space
				} else {
					helpMessage += fmt.Sprintf(boldMessageFormat, token.Word) + space
				}
			}
			helpMessage += dash + space + fmt.Sprintf(italicMessageFormat, command.description) + newLine
		}
		response.Reply(helpMessage)
	})

	return slacker
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	Client         *slack.Client
	rtm            *slack.RTM
	botCommands    []*BotCommand
	initHandler    func()
	errorHandler   func(err string)
	defaultHandler func(request *Request, response ResponseWriter)
}

// Init handle the event when the bot is first connected
func (s *Slacker) Init(initHandler func()) {
	s.initHandler = initHandler
}

// Err handle when errors are encountered
func (s *Slacker) Err(errorHandler func(err string)) {
	s.errorHandler = errorHandler
}

// Default handle when none of the commands are matched
func (s *Slacker) Default(defaultHandler func(request *Request, response ResponseWriter)) {
	s.defaultHandler = defaultHandler
}

// Command define a new command and append it to the list of existing commands
func (s *Slacker) Command(usage string, description string, handler func(request *Request, response ResponseWriter)) {
	s.botCommands = append(s.botCommands, NewBotCommand(usage, description, handler))
}

// Listen receives events from Slack and each is handled as needed
func (s *Slacker) Listen() error {
	go s.rtm.ManageConnection()

	for msg := range s.rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.ConnectedEvent:
			if s.initHandler == nil {
				continue
			}
			go s.initHandler()

		case *slack.MessageEvent:
			if s.isFromBot(event) {
				continue
			}

			if !s.isBotMentioned(event) && !s.isDirectMessage(event) {
				continue
			}
			go s.handleMessage(event)

		case *slack.RTMError:
			if s.errorHandler == nil {
				continue
			}
			go s.errorHandler(event.Error())

		case *slack.InvalidAuthEvent:
			return errors.New(invalidToken)
		}
	}
	return nil
}

func (s *Slacker) sendMessage(text string, channel string) {
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(text, channel))
}

func (s *Slacker) isFromBot(event *slack.MessageEvent) bool {
	info := s.rtm.GetInfo()
	return event.User == info.User.ID || event.BotID != ""
}

func (s *Slacker) isBotMentioned(event *slack.MessageEvent) bool {
	info := s.rtm.GetInfo()
	return strings.Contains(event.Text, fmt.Sprintf(userMentionFormat, info.User.ID))
}

func (s *Slacker) isDirectMessage(event *slack.MessageEvent) bool {
	return strings.HasPrefix(event.Channel, directChannelMarker)
}

func (s *Slacker) handleMessage(event *slack.MessageEvent) {
	response := NewResponse(event.Channel, s.rtm)
	ctx := context.Background()

	for _, cmd := range s.botCommands {
		parameters, isMatch := cmd.Match(event.Text)
		if !isMatch {
			continue
		}

		cmd.Execute(NewRequest(ctx, event, parameters), response)
		return

	}

	if s.defaultHandler != nil {
		s.defaultHandler(NewRequest(ctx, event, &proper.Properties{}), response)
	}
}
