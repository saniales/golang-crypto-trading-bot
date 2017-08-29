package strategies

import (
	"log"
	"time"

	"github.com/nlopes/slack"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchangeWrappers"
	"github.com/shomali11/slacker"
)

var bot *slacker.Slacker

// The following slack integration allows to send messages as a strategy.
// RTM not supported (and usually not requested when trading, this is an automated bot).
var slackIntegrationExample = IntervalStrategy{
	model: StrategyModel{
		Setup: func(exchangeWrappers.ExchangeWrapper, *environment.Market) error {
			// connect slack token
			bot = slacker.NewClient("YOUR-TOKEN-HERE")
			bot.Init(func() {
				log.Println("Slackbot Connected")
			})
			bot.Err(func(err string) {
				log.Println("Error during slack bot connection: ", err)
			})
			return bot.Listen()
		},
		OnUpdate: func(exchangeWrappers.ExchangeWrapper, *environment.Market) error {
			//if updates has requirements
			bot.Client.PostMessage("DESIRED-CHANNEL", "OMG something happening!!!!!", slack.PostMessageParameters{})
			return nil
		},
		OnError: func(err error) {
			//log
		},
	},
	Interval: time.Second * 10,
}
