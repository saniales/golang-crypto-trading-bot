package strategies

import (
	"time"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
	"github.com/saniales/golang-crypto-trading-bot/strategies"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var telegramBot *tb.Bot

var telegramIntegrationExample = strategies.IntervalStrategy{
	Model: strategies.StrategyModel{
		Name: "telegramIntegrationExample",
		Setup: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
			telegramBot, err := tb.NewBot(tb.Settings{
				Token:  "TOKEN_HERE",
				Poller: &tb.LongPoller{Timeout: 10 * time.Second},
			})

			if err != nil {
				return err
			}

			telegramBot.Start()
			return nil
		},
		OnUpdate: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
			telegramBot.Send(&tb.User{
				Username: "YOUR-USERNAME-GROUP-OR-USER",
			}, "OMG SOMETHING HAPPENING!!!!!", tb.SendOptions{})

			/*
				// Optionally it can have options
				telegramBot.Send(tb.User{
					Username: "YOUR-JOINED-GROUP-USERNAME",
				}, "OMG SOMETHING HAPPENING!!!!!", tb.SendOptions{})
			*/
			return nil
		},
		OnError: func(err error) {
			logrus.Errorf("I Got an error %s", err)
			telegramBot.Stop()
		},
		TearDown: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
			telegramBot.Stop()
			return nil
		},
	},
}
