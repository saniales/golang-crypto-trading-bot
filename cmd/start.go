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

package bot

import (
	"fmt"
	"io"
	"os"
	"strings"

	helpers "github.com/saniales/golang-crypto-trading-bot/bot_helpers"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
	"github.com/saniales/golang-crypto-trading-bot/strategies"
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
	contentToMarshal, err := io.ReadAll(configFile)
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
	fmt.Print("Getting configurations ... ")
	if err := initConfigs(); err != nil {
		fmt.Println("Cannot read from configuration file, please create or replace the current one using gobot init")
		return
	}
	fmt.Println("DONE")

	fmt.Print("Getting exchange info ... ")
	wrappers := make([]exchanges.ExchangeWrapper, len(botConfig.ExchangeConfigs))
	for i, config := range botConfig.ExchangeConfigs {
		wrappers[i] = helpers.InitExchange(config, botConfig.SimulationModeOn, config.FakeBalances, config.DepositAddresses)
	}
	fmt.Println("DONE")

	fmt.Print("Getting markets cold info ... ")
	for _, strategyConf := range botConfig.Strategies {
		mkts := make([]*environment.Market, len(strategyConf.Markets))
		for i, mkt := range strategyConf.Markets {
			currencies := strings.SplitN(mkt.Name, "-", 2)
			mkts[i] = &environment.Market{
				Name:           mkt.Name,
				BaseCurrency:   currencies[0],
				MarketCurrency: currencies[1],
			}

			mkts[i].ExchangeNames = make(map[string]string, len(wrappers))

			for _, exName := range mkt.Exchanges {
				mkts[i].ExchangeNames[exName.Name] = exName.MarketName
			}
		}
		err := strategies.MatchWithMarkets(strategyConf.Strategy, mkts)
		if err != nil {
			fmt.Println("Cannot add tactic : ", err)
		}
	}
	fmt.Println("DONE")

	fmt.Println("Starting bot ... ")
	executeBotLoop(wrappers)
	fmt.Println("EXIT, good bye :)")
}

func executeBotLoop(wrappers []exchanges.ExchangeWrapper) {
	strategies.ApplyAllStrategies(wrappers)
}
