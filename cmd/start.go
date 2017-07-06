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

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/botHelpers"
	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/environment"
	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/exchangeWrappers"
	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/strategies"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts trading using saved configs",
	Long:  `Starts trading using saved configs`,
	Run:   executeStartCommand,
}

var startFlags struct {
	Simulate bool
}

var botConfig environment.BotConfig

func init() {
	RootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&startFlags.Simulate, "simulate", "s", false, "Simulates the trades instead of actually doing them")
}

func initConfigs() error {
	configFile, err := os.Open(GlobalFlags.ConfigFile)
	if err != nil {
		return err
	}
	contentToMarshal, err := ioutil.ReadAll(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(contentToMarshal, &botConfig)
	if err != nil {
		return err
	}
	return nil
}

func executeStartCommand(cmd *cobra.Command, args []string) {
	if err := initConfigs(); err != nil {
		fmt.Println("Cannot read from configuration file, please create or replace the current one using gobot init")
		return
	}

	exchangeWrapper := botHelpers.InitExchange(botConfig.Exchange)

	markets, err := botHelpers.InitMarkets(exchangeWrapper)
	if err != nil {
		fmt.Printf("Cannot initialize Markets data : %s\n", err)
		return
	}

	tactics := make(map[string]strategies.Strategy, len(botConfig.Strategies))

	for _, strategyConf := range botConfig.Strategies {
		tactics[strategyConf.Market] = strategies.Get(strategyConf.Strategy)
	}

	err = executeBotLoop(exchangeWrapper, markets, tactics)
	if err != nil {
		fmt.Printf("Bot exited with an error : %s\n", err)
	} else {
		fmt.Println("Bot exited gracefully")
	}
}

func executeBotLoop(wrapper exchangeWrappers.ExchangeWrapper, markets map[string]*environment.Market, tactics map[string]strategies.Strategy) error {
	for marketName, strategy := range tactics {
		market := markets[marketName]
		if strategy.SetUpStrategy(wrapper, market) != nil {
			fmt.Printf("Cannot initialize tactic %s in market %s, ignoring...\n", strategy.Name(), marketName)
			delete(tactics, marketName)
			continue
		}
		//timer := time.NewTicker(time.Minute)

		defer strategy.TearDownStrategy(wrapper, market)
	}

	for {
		if len(tactics) == 0 {
			return errors.New("No available strategy")
		}
		for marketName, strategy := range tactics {
			market := markets[marketName]
			action, limit, amount, err := strategy.OnCandleUpdate(wrapper, market)
			if err != nil {
				fmt.Printf("Error while performing tactic %s in market %s : %s \nstopping that strategy...\n", strategy.Name(), marketName, err)
				delete(tactics, marketName)
			} else {
				err = applyAction(wrapper, *market, action, amount, limit)
				if err != nil {
					fmt.Printf("Error while applying action : strategy %s on market %s, action was %d", strategy.Name(), marketName, action)
				}
			}
		}
	}
}

func applyAction(wrapper exchangeWrappers.ExchangeWrapper, market environment.Market, action strategies.Action, amount float64, limit float64) error {
	switch action {
	case strategies.Buy:
		if startFlags.Simulate == false {
			_, err := wrapper.BuyMarket(market, amount)
			if err != nil {
				return err
			}
			fmt.Printf("Buy Market on : %s market\n", market.Name)
		} else {
			fmt.Printf("Buy Market on : %s market --simulated\n", market.Name)
		}
		break
	case strategies.BuyLimit:
		if startFlags.Simulate == false {
			_, err := wrapper.BuyLimit(market, amount, limit)
			if err != nil {
				return err
			}
			fmt.Printf("Buy Limit on : %s market\n", market.Name)
		} else {
			fmt.Printf("Buy Limit on : %s market --simulated\n", market.Name)
		}
		break
	case strategies.Sell:
		if startFlags.Simulate == false {
			_, err := wrapper.SellMarket(market, amount)
			if err != nil {
				return err
			}
			fmt.Printf("Sell Market on : %s market\n", market.Name)
		} else {
			fmt.Printf("Sell Market on : %s market --simulated\n", market.Name)
		}
		break
	case strategies.SellLimit:
		if startFlags.Simulate == false {
			_, err := wrapper.SellLimit(market, amount, limit)
			if err != nil {
				return err
			}
			fmt.Printf("Sell Limit on : %s market\n", market.Name)
		} else {
			fmt.Printf("Sell Limit on : %s market --simulated\n", market.Name)
		}
		break
	case strategies.Invalid:
		fmt.Println("Invalid action")
		break
	case strategies.DoNothing:
		fmt.Println("Chilling")
	default:
		break
	}
	return nil
}
