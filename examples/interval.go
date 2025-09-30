// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package examples

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
	"github.com/saniales/golang-crypto-trading-bot/strategies"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

// Watch5Sec prints out the info of the market every 5 seconds.
var Watch5Sec = strategies.IntervalStrategy{
	Model: strategies.StrategyModel{
		Name: "Watch5Sec",
		Setup: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) error {
			fmt.Println("Watch5Sec starting")
			return nil
		},
		OnUpdate: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) error {
			_, err := wrappers[0].GetMarketSummary(markets[0])
			if err != nil {
				return err
			}
			logrus.Info(markets)
			logrus.Info(wrappers)
			return nil
		},
		OnError: func(err error) {
			fmt.Println(err)
		},
		TearDown: func(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) error {
			fmt.Println("Watch5Sec exited")
			return nil
		},
	},
	Interval: time.Second * 5,
}

var telegramBot *tb.Bot

// TelegramIntegrationExample send messages to Telegram as a strategy.
var TelegramIntegrationExample = strategies.IntervalStrategy{
	Model: strategies.StrategyModel{
		Name: "TelegramIntegrationExample",
		Setup: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
			telegramBot, err := tb.NewBot(tb.Settings{
				Token:  "YOUR-TELEGRAM-TOKEN",
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

var discordBot *discordgo.Session

// DiscordIntegrationExample sends messages to a specified discord channel.
var DiscordIntegrationExample = strategies.IntervalStrategy{
	Model: strategies.StrategyModel{
		Name: "DiscordIntegrationExample",
		Setup: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
			// Create a new Discord session using the provided bot token.
			discordBot, err := discordgo.New("Bot " + "YOUR-DISCORD-TOKEN")
			if err != nil {
				return err
			}

			go func() {
				err = discordBot.Open()
				if err != nil {
					return
				}
			}()

			//sleep some time
			time.Sleep(time.Second * 5)
			if err != nil {
				return err
			}

			return nil
		},
		OnUpdate: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
			_, err := discordBot.ChannelMessageSend("CHANNEL-ID", "OMG SOMETHING HAPPENING!!!!!")
			if err != nil {
				return err
			}
			return nil
		},
		OnError: func(err error) {
			logrus.Errorf("I Got an error %s", err)
			telegramBot.Stop()
		},
		TearDown: func([]exchanges.ExchangeWrapper, []*environment.Market) error {
			err := discordBot.Close()
			if err != nil {
				return err
			}
			return nil
		},
	},
}
