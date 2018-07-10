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

package strategies

import (
	"log"
	"time"

	"github.com/nlopes/slack"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
	"github.com/shomali11/slacker"
)

var bot *slacker.Slacker

// The following slack integration allows to send messages as a strategy.
// RTM not supported (and usually not requested when trading, this is an automated bot).
var slackIntegrationExample = IntervalStrategy{
	Model: StrategyModel{
		Setup: func(exchanges.ExchangeWrapper, *environment.Market) error {
			// connect slack token
			bot = slacker.NewClient("YOUR-TOKEN-HERE")
			bot.Init(func() {
				log.Println("Slackbot Connected")
			})
			bot.Err(func(err string) {
				log.Println("Error during slack bot connection: ", err)
			})
			go func() {
				err := bot.Listen()
				if err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnUpdate: func(exchanges.ExchangeWrapper, *environment.Market) error {
			//if updates has requirements
			_, _, err := bot.Client.PostMessage("DESIRED-CHANNEL", "OMG something happening!!!!!", slack.PostMessageParameters{})
			return err
		},
		OnError: func(err error) {
			bot.Client.PostMessage("DESIRED-CHANNEL", "Got An Error: "+err.Error(), slack.PostMessageParameters{})
		},
	},
	Interval: time.Second * 10,
}
